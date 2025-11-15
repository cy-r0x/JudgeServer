package cluster

import (
	"net/http"

	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) getNodes(w http.ResponseWriter, r *http.Request) {
	utils.SendResponse(w, 200, h.nodes)
}
