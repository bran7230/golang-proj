package service

import (
	"fmt"
	"golang-proj/internal/models"
)

func ValidateTycoonRequest(r *models.TycoonRequest) error {

	if r.ServerId == "" {
		return fmt.Errorf("Missing server id.")
	}

	for _, player := range r.Players {
		if player.PlayerId <= 0 {
			return fmt.Errorf("Player id is null or <= 0.")
		}

	}

	return nil
}
