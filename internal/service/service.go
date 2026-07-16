package service

import (
	"fmt"
	"golang-proj/internal/models"
	"golang-proj/internal/repository"
)

type TycoonProcessor interface {
	ProcessTycoonData(r *models.TycoonRequest) error
}

type TycoonService struct {
	repo repository.Repository
}

func NewTycoonService(repo repository.Repository) *TycoonService {
	return &TycoonService{repo: repo}
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

	fmt.Printf("Inserted players: %d", len(r.Players))

	return nil
}
