package utils

import (
	"encoding/json"
	"net/http"
)

func SendResponse(w http.ResponseWriter, statusCode int, message any, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	success := statusCode >= 200 && statusCode < 300
	response := map[string]any{
		"success": success,
		"message": message,
		"data":    data,
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(response)
}
