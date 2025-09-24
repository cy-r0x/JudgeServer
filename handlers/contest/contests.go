package contest

import (
	"encoding/json"
	"net/http"

	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) GetContests(w http.ResponseWriter, r *http.Request) {
	data, err := utils.GetContests()
	if err != nil {
		utils.SendResopnse(w, http.StatusInternalServerError, "Internal Server Error")
	}
	conv, err := json.Marshal(data)
	if err != nil {
		utils.SendResopnse(w, http.StatusInternalServerError, "Internal Server Error")
	}
	utils.SendResopnse(w, http.StatusOK, string(conv))
}
