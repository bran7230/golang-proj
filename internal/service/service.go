package service

import (
	"fmt"
	"golang-proj/internal/models"
)

func ValidateTycoonRequest(r *models.TycoonRequest) error {
	if r.Players == nil {
		return fmt.Errorf("Missing player object.")
	}

	for _, player := range r.Players {
		if player.PlayerId <= 0 {
			return fmt.Errorf("Player id is null or <= 0.")
		}
	}

	return nil
}
