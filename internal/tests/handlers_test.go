package server

import (
	"bytes"
	"encoding/json"
	"golang-proj/internal/models"
	"golang-proj/internal/server"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// test commit
func TestHandleSaves(t *testing.T) {
	// test cases
	tests := []struct {
		name               string
		reqBody            models.TycoonRequest
		expectedHTTPStatus int
	}{
		{
			name: "Valid request",
			reqBody: models.TycoonRequest{
				ServerId:  "awsawswas",
				Timestamp: time.Now(),
				Players: []models.Player{
					{
						PlayerId: 1234566,
						Stats: &models.PlayerStats{
							TotalCurrency: 123456,
							Rebirths:      1,
						},
						PlacedObjects: []models.Objects{
							{
								ItemId: "awaswswa",
								Position: models.ObjectPositions{
									X: 23,
									Y: 23,
									Z: 0,
								},
								Rotation: models.ObjectRotation{
									Y: 23,
								},
							},
						},
					},
				},
			},
			expectedHTTPStatus: http.StatusOK,
		},
		{
			name: "Missing serverID",
			reqBody: models.TycoonRequest{
				ServerId:  "",
				Timestamp: time.Now(),
				Players: []models.Player{
					{
						PlayerId: 1234566,
						Stats: &models.PlayerStats{
							TotalCurrency: 123456,
							Rebirths:      1,
						},
						PlacedObjects: []models.Objects{
							{
								ItemId: "awaswswa",
								Position: models.ObjectPositions{
									X: 23,
									Y: 23,
									Z: 0,
								},
								Rotation: models.ObjectRotation{
									Y: 23,
								},
							},
						},
					},
				},
			},
			expectedHTTPStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid player ID(Zero)",
			reqBody: models.TycoonRequest{
				ServerId:  "awsawswas",
				Timestamp: time.Now(),
				Players: []models.Player{
					{
						PlayerId: 0,
						Stats: &models.PlayerStats{
							TotalCurrency: 123456,
							Rebirths:      1,
						},
						PlacedObjects: []models.Objects{
							{
								ItemId: "awaswswa",
								Position: models.ObjectPositions{
									X: 23,
									Y: 23,
									Z: 0,
								},
								Rotation: models.ObjectRotation{
									Y: 23,
								},
							},
						},
					},
				},
			},
			expectedHTTPStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid player id(missing completly)",
			reqBody: models.TycoonRequest{
				ServerId:  "awsawswas",
				Timestamp: time.Now(),
				Players: []models.Player{
					{
						Stats: &models.PlayerStats{
							TotalCurrency: 123456,
							Rebirths:      1,
						},
						PlacedObjects: []models.Objects{
							{
								ItemId: "awaswswa",
								Position: models.ObjectPositions{
									X: 23,
									Y: 23,
									Z: 0,
								},
								Rotation: models.ObjectRotation{
									Y: 23,
								},
							},
						},
					},
				},
			},
			expectedHTTPStatus: http.StatusBadRequest,
		},
		{
			name: "Missing timestamp.",
			reqBody: models.TycoonRequest{
				ServerId: "awsawswas",
				Players: []models.Player{
					{
						PlayerId: 1234566,
						Stats: &models.PlayerStats{
							TotalCurrency: 123456,
							Rebirths:      1,
						},
						PlacedObjects: []models.Objects{
							{
								ItemId: "awaswswa",
								Position: models.ObjectPositions{
									X: 23,
									Y: 23,
									Z: 0,
								},
								Rotation: models.ObjectRotation{
									Y: 23,
								},
							},
						},
					},
				},
			},
			expectedHTTPStatus: http.StatusBadRequest,
		},
		{
			name: "No players online.",
			reqBody: models.TycoonRequest{
				ServerId:  "awsawswas",
				Timestamp: time.Now(),
				Players:   []models.Player{},
			},
			expectedHTTPStatus: http.StatusOK,
		},
		{
			name: "Player has no objects placed.",
			reqBody: models.TycoonRequest{
				ServerId:  "awsawswas",
				Timestamp: time.Now(),
				Players: []models.Player{
					{
						PlayerId: 1234566,
						Stats: &models.PlayerStats{
							TotalCurrency: 123456,
							Rebirths:      1,
						},
						PlacedObjects: []models.Objects{},
					},
				},
			},
			expectedHTTPStatus: http.StatusOK,
		},
		{
			name: "Player has no currency.",
			reqBody: models.TycoonRequest{
				ServerId:  "awsawswas",
				Timestamp: time.Now(),
				Players: []models.Player{
					{
						PlayerId: 1234566,
						Stats: &models.PlayerStats{
							TotalCurrency: 0,
							Rebirths:      1,
						},
						PlacedObjects: []models.Objects{
							{
								ItemId: "awaswswa",
								Position: models.ObjectPositions{
									X: 23,
									Y: 23,
									Z: 0,
								},
								Rotation: models.ObjectRotation{
									Y: 23,
								},
							},
						},
					},
				},
			},
			expectedHTTPStatus: http.StatusOK,
		},
		{
			name: "Player has no rebirths.",
			reqBody: models.TycoonRequest{
				ServerId:  "awsawswas",
				Timestamp: time.Now(),
				Players: []models.Player{
					{
						PlayerId: 1234566,
						Stats: &models.PlayerStats{
							TotalCurrency: 123456,
							Rebirths:      0,
						},
						PlacedObjects: []models.Objects{
							{
								ItemId: "awaswswa",
								Position: models.ObjectPositions{
									X: 23,
									Y: 23,
									Z: 0,
								},
								Rotation: models.ObjectRotation{
									Y: 23,
								},
							},
						},
					},
				},
			},
			expectedHTTPStatus: http.StatusOK,
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
