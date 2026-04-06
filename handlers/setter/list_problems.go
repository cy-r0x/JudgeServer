package setter

import (
	"log"
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) ListSetterProblems(w http.ResponseWriter, r *http.Request) {
	payload, ok := r.Context().Value("user").(*middlewares.Payload)
	if !ok {
		utils.SendResponse(w, http.StatusUnauthorized, "User information not found")
		return
	}
	setterId := payload.Sub

	var problems []Problem
	err := h.db.Model(&models.Problem{}).
		Select("id", "title", "created_at").
		Where("created_by = ?", setterId).
		Order("created_at DESC").
		Scan(&problems).Error

	if err != nil {
		log.Println("Error fetching setter problems:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch problems")
		return
	}

	if problems == nil {
		problems = []Problem{}
	}

	utils.SendResponse(w, http.StatusOK, problems)
}
