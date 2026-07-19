package server

import (
	"errors"
	"log"
	"net/http"

	"golang-proj/internal/models"
	"golang-proj/internal/service"
)

// HandleSaves TODO: Implement hmac-sha256 encoding in the response / requests.
func HandleSaves(tycoonSvc service.TycoonProcessor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestData, decodingErr := decode[models.TycoonRequest](w, r)

		if decodingErr != nil {
			var maxBytesError *http.MaxBytesError
			if errors.As(decodingErr, &maxBytesError) {
				err := sendError(w, http.StatusRequestEntityTooLarge, decodingErr.Error())
				if err != nil {
					return
				}
				return
			}
			err := sendError(w, http.StatusBadRequest, decodingErr.Error())
			if err != nil {
				return
			}
			return
		}

		if err := service.ValidateTycoonRequest(&requestData); err != nil {
			err := sendError(w, http.StatusBadRequest, err.Error())
			if err != nil {
				return
			}
			return
		}

		if err := tycoonSvc.ProcessTycoonData(&requestData); err != nil {
			if errors.Is(err, service.ErrQueueFull) {
				err := sendError(w, http.StatusTooManyRequests, "Server is busy, please try again later.")
				if err != nil {
					return
				}
				return
			}
			err := sendError(w, http.StatusInternalServerError, err.Error())
			if err != nil {
				return
			}
			return
		}

		response := map[string]any{
			"message": "Accepted payload, processing now!",
		}
		if err := encode(w, http.StatusAccepted, response); err != nil {
			log.Print("Failed to encode response: ", err)

			http.Error(w, "Error parsing data.", http.StatusInternalServerError)
			return
		}

		if tycoonSvc == nil {
			log.Print("Tycoon service dependency is not configured")
			return
		}
	}
}
