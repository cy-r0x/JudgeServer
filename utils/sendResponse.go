package utils

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool `json:"success"`
	Message any  `json:"message"`
	Data    any  `json:"data"`
}

func SendResponse(w http.ResponseWriter, statusCode int, message any, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	res := APIResponse{
		Success: statusCode >= 200 && statusCode < 300,
		Message: message,
		Data:    data,
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
