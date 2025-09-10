package handlers

import (
	"encoding/json"
	"net/http"
)

func HandleContests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	data, err := GetContests()
	if err != nil {

	}
	encoder.Encode(data)
}
