package server

import (
	"errors"
	"golang-proj/internal/models"
	"golang-proj/internal/service"
	"io"
	"log"
	"net/http"
)

// HandleSaves processes tycoon save data with HMAC-SHA256 verification.
func HandleSaves(tycoonSvc service.TycoonProcessor, secretKey []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Enforce size limit to prevent DOS attacks (1MB limit)
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

		// Read the actual request body
		bodyBytes, readErr := io.ReadAll(r.Body)
		if readErr != nil {
			var maxBytesError *http.MaxBytesError
			if errors.As(readErr, &maxBytesError) {
				err := SendError(w, http.StatusRequestEntityTooLarge, readErr.Error())
				if err != nil {
					return
				}
				return
			}
			err := SendError(w, http.StatusBadRequest, readErr.Error())
			if err != nil {
				return
			}
			return
		}

		requestData, decodingErr := decode[models.TycoonRequest](bodyBytes, secretKey, r.Header.Get("HMAC-Signature"))

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
