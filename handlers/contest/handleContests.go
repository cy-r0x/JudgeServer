package contest

import (
	"encoding/json"
	"net/http"

	"github.com/judgenot0/judge-backend/utils"
)

func HandleContests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	data, err := utils.GetContests()
	if err != nil {

	}
	encoder.Encode(data)
}
