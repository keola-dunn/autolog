package httputil

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Status       string `json:"status"`
	StatusCode   int    `json:"statusCode"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

func RespondWithError(w http.ResponseWriter, statusCode int, errorMessage string) {
	data, _ := json.Marshal(ErrorResponse{
		Status:       http.StatusText(statusCode),
		StatusCode:   statusCode,
		ErrorMessage: errorMessage,
	})
	w.WriteHeader(statusCode)
	w.Write(data)
}

func RespondWithJSON(w http.ResponseWriter, statusCode int, responseBody any) error {
	data, err := json.Marshal(responseBody)
	if err != nil {
		return fmt.Errorf("failed to marshal response body as expected: %w", err)
	}
	w.WriteHeader(statusCode)
	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("failed to write response body: %w", err)
	}
	return nil
}
