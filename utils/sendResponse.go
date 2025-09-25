package utils

import (
	"encoding/json"
	"net/http"
)

func SendResopnse(w http.ResponseWriter, statusCode int, message any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encoder := json.NewEncoder(w)
	encoder.Encode(message)
}
