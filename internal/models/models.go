package models

type TestRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type TestResponse struct {
	ErrorCode int    `json:"errorCode"`
	Data      string `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
