package server

import (
	"errors"
	"log"
	"net/http"

	"golang-proj/internal/models"
	"golang-proj/internal/service"
)

// TODO: Implement hmac-sha256 encoding in the response / requests.
func HandleSaves(tycoonSvc service.TycoonProcessor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestData, decodingErr := decode[models.TycoonRequest](w, r)

		if decodingErr != nil {
			var maxBytesError *http.MaxBytesError
			if errors.As(decodingErr, &maxBytesError) {
				sendError(w, http.StatusRequestEntityTooLarge, decodingErr.Error())
				return
			}
			sendError(w, http.StatusBadRequest, decodingErr.Error())
			return
		}

		if err := service.ValidateTycoonRequest(&requestData); err != nil {
			sendError(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := tycoonSvc.ProcessTycoonData(&requestData); err != nil {
			if errors.Is(err, service.ErrQueueFull) {
				sendError(w, http.StatusTooManyRequests, "Server is busy, please try again later.")
				return
			}
			sendError(w, http.StatusInternalServerError, err.Error())
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
