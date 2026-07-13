package server

import (
	"bytes"
	"encoding/json"
	"golang-proj/internal/models"
	"golang-proj/internal/server"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// test commit
func TestHandleSaves(t *testing.T) {
	// test cases
	tests := []struct {
		name               string
		reqBody            models.TestRequest
		expectedHTTPStatus int
	}{
		{
			name: "Valid request",
			reqBody: models.TestRequest{
				Name:  "TestUser",
				Email: "Test",
				Age:   123,
			},
			expectedHTTPStatus: http.StatusOK,
		},
		{
			name: "Missing name field",
			reqBody: models.TestRequest{
				Email: "Test",
				Age:   123,
			},
			expectedHTTPStatus: http.StatusBadRequest,
		},
		{
			name: "Missing email(invalid req)",
			reqBody: models.TestRequest{
				Name: "TestUser",
				Age:  123,
			},
			expectedHTTPStatus: http.StatusBadRequest,
		},
		{
			name: "Missing age field",
			reqBody: models.TestRequest{
				Name:  "TestUser",
				Email: "Test",
			},
			expectedHTTPStatus: http.StatusBadRequest,
		},
		{
			name: "Over 1mb limit.",
			reqBody: models.TestRequest{
				Name:  strings.Repeat("A", (1<<20)+100),
				Age:   123,
				Email: "Test",
			},
			expectedHTTPStatus: http.StatusRequestEntityTooLarge,
		},
	}

	// iterate over each test case
	for _, tc := range tests {
		// run eac test case as a sub case(isolated test)
		t.Run(tc.name, func(t *testing.T) {
			// parse the json
			bodyBytes, err := json.Marshal(tc.reqBody)
			if err != nil {
				t.Fatal(err)
			}

			// post to test endpoint
			req, err := http.NewRequest("POST", "/", bytes.NewBuffer(bodyBytes))
			if err != nil {
				t.Fatal(err)
			}

			// json type
			req.Header.Set("Content-Type", "application/json")

			// setup recorder
			rr := httptest.NewRecorder()

			// test the endpoint
			handler := http.HandlerFunc(server.HandleSaves)
			handler.ServeHTTP(rr, req)

			// validate status code against the expected value in the struct
			if status := rr.Code; status != tc.expectedHTTPStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedHTTPStatus)
			}

			// decode response
			var resp models.TestResponse
			err = json.NewDecoder(rr.Body).Decode(&resp)
			if err != nil {
				t.Fatal("failed to decode response")
			}
		})
	}

}
