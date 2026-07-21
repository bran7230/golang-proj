package server

import (
	"errors"
	"golang-proj/internal/models"
	"golang-proj/internal/service"
	"log/slog"
	"net/http"
)

// HandleSaves TODO: Implement hmac-sha256 encoding in the response / requests.
func HandleSaves(tycoonSvc service.TycoonProcessor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestData, decodingErr := decode[models.TycoonRequest](w, r)

		if decodingErr != nil {
			var maxBytesError *http.MaxBytesError
			if errors.As(decodingErr, &maxBytesError) {
				err := SendError(w, http.StatusRequestEntityTooLarge, decodingErr.Error())
				if err != nil {
					return
				}
				return
			}
			err := SendError(w, http.StatusBadRequest, decodingErr.Error())
			if err != nil {
				return
			}
			return
		}

		if err := service.ValidateTycoonRequest(&requestData); err != nil {
			err := SendError(w, http.StatusBadRequest, err.Error())
			if err != nil {
				return
			}
			return
		}

		if tycoonSvc == nil {
			err := SendError(w, http.StatusInternalServerError, "fatal error, our validation is currently down. Please try again later")
			if err != nil {
				return
			}
		}

		if err := tycoonSvc.ProcessTycoonData(&requestData); err != nil {
			if errors.Is(err, service.ErrQueueFull) {
				err := SendError(w, http.StatusTooManyRequests, "Server is busy, please try again later.")
				if err != nil {
					return
				}
				return
			}
			err := SendError(w, http.StatusInternalServerError, err.Error())
			if err != nil {
				return
			}
			return
		}

		response := map[string]any{
			"message": "Accepted payload, processing now!",
		}
		if err := encode(w, http.StatusAccepted, response); err != nil {
			slog.Error("Failed to encode response", "error", err.Error(), slog.Group("response",
				slog.Any("response", response),
			))

			http.Error(w, "Error parsing data.", http.StatusInternalServerError)
			return
		}

		if tycoonSvc == nil {
			slog.Error("Tycoon service dependency is not configured")
			return
		}
	}
}
