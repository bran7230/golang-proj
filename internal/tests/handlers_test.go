package server

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"golang-proj/internal/models"
	"golang-proj/internal/server"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

type stubTycoonService struct{}

func (s stubTycoonService) ProcessTycoonData(*models.TycoonRequest) error {
	return nil
}

// generateHMAC creates an HMAC-SHA256 signature for testing
func generateHMAC(bodyBytes []byte, secretKey string) string {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write(bodyBytes)
	return hex.EncodeToString(mac.Sum(nil))
}

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
			expectedHTTPStatus: http.StatusAccepted,
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
			expectedHTTPStatus: http.StatusAccepted,
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
			expectedHTTPStatus: http.StatusAccepted,
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
			expectedHTTPStatus: http.StatusAccepted,
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
			expectedHTTPStatus: http.StatusAccepted,
		},
		{
			name: "Player has no stats(should NEVER BE A NULL OBJ!!).",
			reqBody: models.TycoonRequest{
				ServerId:  "awsawswas",
				Timestamp: time.Now(),
				Players: []models.Player{
					{
						PlayerId: 1234566,
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
			name: "Json file is too large(more than 1mb)",
			reqBody: models.TycoonRequest{
				ServerId:  strings.Repeat("A", 1<<26),
				Timestamp: time.Now(),
				Players: []models.Player{
					{
						PlayerId: 1234566,
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
			expectedHTTPStatus: http.StatusRequestEntityTooLarge,
		},
	}

	secretKey := os.Getenv("HMAC_KEY")

	if secretKey == "" {
		secretKey = "some-cool-secret"
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

			// Generate the HMAC
			hmacSig := generateHMAC(bodyBytes, secretKey)
			req.Header.Set("HMAC-Signature", hmacSig)

			// setup recorder
			rr := httptest.NewRecorder()

			// test the endpoint with the secret key from env
			handler := server.HandleSaves(stubTycoonService{}, []byte(secretKey))
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
