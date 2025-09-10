package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func HandleProblemList(w http.ResponseWriter, r *http.Request) {
	contestIdStr := r.PathValue("contestId")
	contestId, err := strconv.Atoi(contestIdStr)
	if err != nil {
		http.Error(w, "Invalid contestId", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	data, err := GetProblems(contestId)
	if err != nil {

	}
	encoder.Encode(data)
}
