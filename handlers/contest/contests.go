package contest

import (
	"net/http"

	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) GetContests(w http.ResponseWriter, r *http.Request) {
	data, err := utils.GetContests()
	if err != nil {
		utils.SendResopnse(w, http.StatusInternalServerError, "Internal Server Error")
	}
	if err != nil {
		utils.SendResopnse(w, http.StatusInternalServerError, "Internal Server Error")
	}
	utils.SendResopnse(w, http.StatusOK, data)
}
