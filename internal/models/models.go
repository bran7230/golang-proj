package models

import "time"

type TestRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type TestResponse struct {
	ErrorCode int `json:"errorCode"`
	Data      any `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type TycoonRequest struct {
	ServerId  string    `json:"serverId"`
	Timestamp time.Time `json:"timestamp"`
	Players   []Player  `json:"players"`
}

type Player struct {
	PlayerId      int64       `json:"playerId"`
	Stats         PlayerStats `json:"stats"`
	PlacedObjects []Objects   `json:"placedObjects"`
}

type PlayerStats struct {
	TotalCurrency int `json:"totalCurrency"`
	Rebirths      int `json:"rebirths"`
}

type Objects struct {
	ItemId   string          `json:"itemId"`
	Position ObjectPositions `json:"position"`
	Rotation ObjectRotation  `json:"rotation"`
}

type ObjectPositions struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type ObjectRotation struct {
	Y float64 `json:"y"`
}
