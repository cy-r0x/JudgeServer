package problem

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/judgenot0/judge-backend/utils"
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
	data, err := utils.GetProblems(contestId)
	if err != nil {

	}
	encoder.Encode(data)
}
