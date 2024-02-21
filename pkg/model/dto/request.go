package dto

import (
	"context"
	"net/http"
)

type ClientResponseDto struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

func NewClientResponseDto(w http.ResponseWriter, statusCode int, message string, payload interface{}) {
	response := ClientResponseDto{
		Status:  statusCode,
		Message: message,
		Payload: payload,
	}
	sendData(w, response)
}

func NewSuccessClientResponseDto(ctx context.Context, w http.ResponseWriter, payload interface{}) {
	NewClientResponseDto(w, 200, "success", payload)
}

func NewErrorClientResponseDto(ctx context.Context, w http.ResponseWriter, statusCode int, message string) {
	NewClientResponseDto(w, statusCode, message, "")
}

func sendData(w http.ResponseWriter, response ClientResponseDto) {
	responseJSON, err := response.MarshalJSON()
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
