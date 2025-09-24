package submissions

import (
	"net/http"

	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) GetSubmissions(w http.ResponseWriter, r *http.Request) {
	utils.SendResopnse(w, http.StatusCreated, "Hoise beda lo")
}
