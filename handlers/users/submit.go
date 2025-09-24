package users

import (
	"net/http"

	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) HandleSubmit(w http.ResponseWriter, r *http.Request) {
	utils.SendResopnse(w, http.StatusCreated, "Hoise beda lo")
}
