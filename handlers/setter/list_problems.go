package setter

import (
	"log"
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) ListSetterProblems(w http.ResponseWriter, r *http.Request) {
	payload, ok := r.Context().Value("user").(*middlewares.Payload)
	if !ok {
		utils.SendResponse(w, http.StatusUnauthorized, "User information not found")
		return
	}
	setterId := payload.Sub

	problems := []Problem{}
	query := `
		SELECT id, title, created_at
		FROM problems
		WHERE created_by = $1
		ORDER BY created_at DESC
	`

	err := h.db.Select(&problems, query, setterId)
	if err != nil {
		log.Println("Error fetching setter problems:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch problems")
		return
	}

	utils.SendResponse(w, http.StatusOK, problems)
}
