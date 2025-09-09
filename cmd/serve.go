package cmd

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/judgenot0/judge-backend/handlers"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	data, err := handlers.GetContests()
	if err != nil {

	}
	encoder.Encode(data)
}

func handleProblemList(w http.ResponseWriter, r *http.Request) {
	contestIdStr := r.PathValue("contestId")
	contestId, err := strconv.Atoi(contestIdStr)
	if err != nil {
		http.Error(w, "Invalid contestId", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	data, err := handlers.GetProblems(contestId)
	if err != nil {

	}
	encoder.Encode(data)
}

func Serve() {
	mux := http.NewServeMux()
	routes(mux)
	http.ListenAndServe(":8080", mux)
}
