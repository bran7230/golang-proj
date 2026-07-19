package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang-proj/internal/models"
	"golang-proj/internal/repository"
	"log"
	"strings"
	"sync"
)

var ErrQueueFull = errors.New("save queue is full, system is overloaded")

type TycoonProcessor interface {
	ProcessTycoonData(r *models.TycoonRequest) error
}

type TycoonService struct {
	repo  repository.Repository
	queue chan *models.TycoonRequest
	wg    sync.WaitGroup
}

func NewTycoonService(repo repository.Repository, queueSize, workerCount int) *TycoonService {
	s := &TycoonService{
		repo:  repo,
		queue: make(chan *models.TycoonRequest, queueSize),
	}

	// Spin up the background worker pool
	for range workerCount {
		go s.worker()
	}

	return s
}

// ValidateTycoonRequest uses the request to verify that they have the required fields.
// This runs after the client receives a 202(accepted) response code.
func ValidateTycoonRequest(r *models.TycoonRequest) error {

	if r.ServerId == "" {
		return fmt.Errorf("missing server id / invalid")
	}

	if r.Timestamp.IsZero() {
		return fmt.Errorf("timestamp is missing / invalid")
	}

	for _, player := range r.Players {
		if player.PlayerId <= 0 {
			return fmt.Errorf("player id is null or <= 0")
		}
		if player.Stats == nil {
			return fmt.Errorf("player has no stats(it's null)")
		}
	}

	return nil
}

func (s *TycoonService) ProcessTycoonData(r *models.TycoonRequest) error {
	if r == nil {
		return fmt.Errorf("request cannot be processed / is null")
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
	defer s.wg.Done()
	for r := range s.queue {
		if err := s.insertData(r); err != nil {
			log.Printf("Worker failed to process data: %v", err)
		}
	}
}

func (s *TycoonService) Close() {
	close(s.queue)
	s.wg.Wait()
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
	            VALUES `

	var queryPlaceholders []string
	var args []any

	placeholderIndex := 1
	for _, player := range r.Players {
		placedObjects, err := json.Marshal(player.PlacedObjects)
		if err != nil {
			return err
		}

		// Create the placeholder string for this specific player
		placeholder := fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)",
			placeholderIndex, placeholderIndex+1, placeholderIndex+2,
			placeholderIndex+3, placeholderIndex+4, placeholderIndex+5)

		queryPlaceholders = append(queryPlaceholders, placeholder)

		// Add the actual variables to the arguments slice
		args = append(args,
			player.PlayerId,
			player.Stats.TotalCurrency,
			player.Stats.Rebirths,
			placedObjects,
			r.ServerId,
			r.Timestamp,
		)

		placeholderIndex += 6
	}

	queryTail := `
				ON CONFLICT(player_id)
	            DO UPDATE SET
					total_currency = EXCLUDED.total_currency,
	                rebirths = EXCLUDED.rebirths,
					placed_objects = EXCLUDED.placed_objects,
					current_server_id = EXCLUDED.current_server_id,
					date_last_updated = EXCLUDED.date_last_updated
	            `

	// Join the placeholders with commas
	valuesSection := strings.Join(queryPlaceholders, ",")
	query := queryHeader + valuesSection + queryTail

	// Pass the arguments alongside the query string
	err := s.repo.InsertUser(query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert player, error: %s", err)
	}

	fmt.Printf("Affected players: %d\n", len(r.Players))
	return nil
}
