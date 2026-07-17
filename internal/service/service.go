package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang-proj/internal/models"
	"golang-proj/internal/repository"
	"log"
	"strings"
)

var ErrQueueFull = errors.New("save queue is full, system is overloaded.")

type TycoonProcessor interface {
	ProcessTycoonData(r *models.TycoonRequest) error
}

type TycoonService struct {
	repo  repository.Repository
	queue chan *models.TycoonRequest
}

func NewTycoonService(repo repository.Repository, queueSize, workerCount int) *TycoonService {
	s := &TycoonService{
		repo:  repo,
		queue: make(chan *models.TycoonRequest, queueSize),
	}

	// Spin up the background worker pool
	for i := 0; i < workerCount; i++ {
		go s.worker()
	}

	return s
}

func ValidateTycoonRequest(r *models.TycoonRequest) error {

	if r.ServerId == "" {
		return fmt.Errorf("Missing server id / invalid.")
	}

	if r.Timestamp.IsZero() {
		return fmt.Errorf("Timestamp is missing / invalid.")
	}

	for _, player := range r.Players {
		if player.PlayerId <= 0 {
			return fmt.Errorf("Player id is null or <= 0.")
		}
		if player.Stats == nil {
			return fmt.Errorf("Player has no stats(it's null).")
		}
	}

	return nil
}

func (s *TycoonService) ProcessTycoonData(r *models.TycoonRequest) error {
	if r == nil {
		return fmt.Errorf("Request cannot be processed / is null.")
	}
	// queue the requests
	select {
	case s.queue <- r:
		return nil
	default:
		return ErrQueueFull
	}
}

func (s *TycoonService) worker() {
	for r := range s.queue {
		if err := s.insertData(r); err != nil {
			log.Printf("Worker failed to process data: %v", err)
		}
	}
}

func (s *TycoonService) insertData(r *models.TycoonRequest) error {

	queryHeader := `
	            INSERT INTO players (
				player_id,
				total_currency,
				rebirths,
				placed_objects,
				current_server_id,
				date_last_updated)
	            VALUES`

	queryValues := make([]string, len(r.Players))

	queryTail := `
				ON CONFLICT(player_id)
	            DO UPDATE SET
					total_currency = EXCLUDED.total_currency,
	                rebirths = EXCLUDED.rebirths,
					placed_objects = EXCLUDED.placed_objects,
					current_server_id = EXCLUDED.current_server_id,
					date_last_updated = EXCLUDED.date_last_updated
	            `
	for i, player := range r.Players {
		placedObjects, err := json.Marshal(player.PlacedObjects)
		if err != nil {
			return err
		}
		// Escape any single quotes in the JSON string to avoid SQL syntax errors
		placedEsc := strings.ReplaceAll(string(placedObjects), "'", "''")
		// Format a single row of VALUES. Use explicit formatting to avoid fmt.Sprintf extra-args behavior.
		queryValues[i] = fmt.Sprintf("(%d,%d,%d,'%s','%s','%s')", player.PlayerId,
			player.Stats.TotalCurrency,
			player.Stats.Rebirths,
			placedEsc,
			r.ServerId,
			r.Timestamp.Format("2006-01-02 15:04:05"))
	}

	valuesSection := strings.Join(queryValues, ",")
	query := queryHeader + "\n" + valuesSection + "\n" + queryTail

	err := s.repo.InsertUser(query)
	if err != nil {
		return fmt.Errorf("failed to insert player, error: %s", err)
	}

	fmt.Printf("Affected players: %d\n", len(r.Players))
	return nil
}
