package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang-proj/internal/models"
	"golang-proj/internal/repository"
	"log/slog"
	"strings"
	"sync"
	"time"
)

var ErrQueueFull = errors.New("save queue is full, system is overloaded")

type TycoonProcessor interface {
	ProcessTycoonData(r *models.TycoonRequest) error
}

type TycoonService struct {
	repo   repository.Repository
	queue  chan *models.TycoonRequest
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

// NewTycoonService creates a TycoonService with worker goroutines. Returns an error if repo is nil.
func NewTycoonService(repo repository.Repository, queueSize, workerCount int) (*TycoonService, error) {
	if repo == nil {
		return nil, fmt.Errorf("repository cannot be nil")
	}
	if queueSize <= 0 {
		queueSize = 1000
	}
	if workerCount <= 0 {
		workerCount = 1
	}

	ctx, cancel := context.WithCancel(context.Background())

	s := &TycoonService{
		repo:   repo,
		queue:  make(chan *models.TycoonRequest, queueSize),
		ctx:    ctx,
		cancel: cancel,
	}

	// Spin up the background worker pool
	for i := 0; i < workerCount; i++ {
		s.wg.Add(1)
		go s.worker()
	}

	// background goroutine to log queue length periodically (optional)
	go func() {
		t := time.NewTicker(30 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-s.ctx.Done():
				return
			}
		}
	}()

	return s, nil
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
	for {
		select {
		case <-s.ctx.Done():
			return
		case r, ok := <-s.queue:
			if !ok {
				return
			}
			if err := s.insertData(r); err != nil {
				slog.Warn("Worker failed to process data", "error", err.Error())
			}
		}
	}
}

func (s *TycoonService) Close() {
	// cancel background tasks and stop accepting new work
	if s.cancel != nil {
		s.cancel()
	}
	// close the queue to allow workers to drain remaining items
	select {
	case <-s.ctx.Done():
		// already cancelled
	default:
	}
	close(s.queue)
	// wait for workers to finish
	s.wg.Wait()
}

func (s *TycoonService) insertData(r *models.TycoonRequest) error {
	if len(r.Players) == 0 {
		return nil
	}

	queryHeader := `
				INSERT INTO player (
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
			return fmt.Errorf("failed to marshal placed_objects: %w", err)
		}

		// Create the placeholder string for this specific player
		placeholder := fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)",
			placeholderIndex, placeholderIndex+1, placeholderIndex+2,
			placeholderIndex+3, placeholderIndex+4, placeholderIndex+5)

		queryPlaceholders = append(queryPlaceholders, placeholder)

		// Add the actual variables to the arguments slice (parameterized to prevent SQL injection)
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

	// Pass the arguments alongside the query string (parameterized, safe from injection)
	err := s.repo.InsertUser(query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert players: %w", err)
	}

	slog.Debug("Successfully inserted/updated players", "playerCount", len(r.Players), "serverId", r.ServerId)
	return nil
}
