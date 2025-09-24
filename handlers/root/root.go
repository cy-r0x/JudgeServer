package root

import (
	"net/http"

	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) Root(w http.ResponseWriter, r *http.Request) {
	utils.SendResopnse(w, http.StatusOK, "Noice")
}
