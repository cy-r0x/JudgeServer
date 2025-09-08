package main

import (
	"encoding/json"
	"net/http"

	"github.com/judgenot0/judge-backend/handlers"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	data, err := handlers.GetContest()
	if err != nil {

	}
	encoder.Encode(data)
}

func main() {
	mux := http.NewServeMux()
	// mux.HandleFunc("/", handleRoot)
	mux.Handle("GET /", http.HandlerFunc(handleRoot))
	http.ListenAndServe(":8080", mux)
}
