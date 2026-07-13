package server

import (
	"errors"
	"log"
	"net/http"

	"golang-proj/internal/models"
)

// TODO: Implement hmac-sha256 encoding in the response / requests.
func HandleSaves(w http.ResponseWriter, r *http.Request) {

	requestData, decodingErr := decode[models.TestRequest](w, r)

	if decodingErr != nil {
		var maxBytesError *http.MaxBytesError
		if errors.As(decodingErr, &maxBytesError) {
			sendError(w, http.StatusRequestEntityTooLarge, "Payload exceeded max limit.")
			return
		}
		sendError(w, http.StatusBadRequest, "error parsing data.")
		return
	}

	if requestData.Name == "" {
		sendError(w, http.StatusBadRequest, "Name field missing.")
		return
	} else if requestData.Age <= 0 {
		sendError(w, http.StatusBadRequest, "Age field missing / is 0.")
		return
	} else if requestData.Email == "" {
		sendError(w, http.StatusBadRequest, "Email field is missing.")
		return
	}

	response := models.TestResponse{
		ErrorCode: 200,
		Data: []any{
			"123",
			123,
		},
	}
	err := encode(w, 200, response)

	if err != nil {
		log.Print("Failed to encode response: ", err)

		http.Error(w, "Error parsing data.", http.StatusInternalServerError)
		return
	}

}
