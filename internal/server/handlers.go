package server

import (
	"log"
	"net/http"

	"golang-proj/internal/models"
)

func HandleSaves(w http.ResponseWriter, r *http.Request) {

	requestData, decodingErr := decode[models.TestRequest](r)

	if decodingErr != nil {
		sendError(w, http.StatusBadRequest, "error parsing data.")
		return
	}

	if requestData.Email == "" {
		sendError(w, http.StatusBadRequest, "Missing email field.")
		return
	} else if requestData.Age <= 0 {
		sendError(w, http.StatusBadRequest, "Missing age field.")
		return
	} else if requestData.Name == "" {
		sendError(w, http.StatusBadRequest, "Missing name field")
		return
	}

	response := models.TestResponse{ErrorCode: 200, Data: "test"}
	err := encode(w, 200, response)

	if err != nil {
		log.Print("Failed to encode response: ", err)

		http.Error(w, "Error parsing data.", http.StatusInternalServerError)
		return
	}

}
