package server

import (
	"errors"
	"log"
	"net/http"

	"golang-proj/internal/models"
	"golang-proj/internal/service"
)

// TODO: Implement hmac-sha256 encoding in the response / requests.
func HandleSaves(w http.ResponseWriter, r *http.Request) {

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

	response := models.TestResponse{
		ErrorCode: 200,
		Data: []any{
			"123",
			123,
		},
	}
	err := encode(w, http.StatusOK, response)

	if err != nil {
		log.Print("Failed to encode response: ", err)

		http.Error(w, "Error parsing data.", http.StatusInternalServerError)
		return
	}

}
