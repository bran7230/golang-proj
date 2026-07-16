package service

import (
	"errors"
	"fmt"
	"golang-proj/internal/models"
	"golang-proj/internal/repository"
	"log"
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
	for _, player := range r.Players {
		query := `
            INSERT INTO players (player_id, rebirths) 
            VALUES ($1, $2)
            ON CONFLICT(player_id)
            DO UPDATE SET
                rebirths = EXCLUDED.rebirths
            `
		err := s.repo.InsertUser(query, player.PlayerId, player.Stats.Rebirths)
		if err != nil {
			return fmt.Errorf("failed to insert player %d: %w", player.PlayerId, err)
		}
	}
	fmt.Printf("Affected players: %d\n", len(r.Players))
	return nil
}
