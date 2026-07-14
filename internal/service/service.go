package service

import (
	"fmt"
	"golang-proj/internal/models"
)

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

func ProcessTycoonData(r *models.TycoonRequest) error {
	return fmt.Errorf("Failed to decode data.")
}
