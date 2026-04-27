package problem

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) UpdateProblem(w http.ResponseWriter, r *http.Request) {
	payload, ok := r.Context().Value("user").(*middlewares.Payload)
	if !ok {
		utils.SendResponse(w, http.StatusUnauthorized, "User information not found", nil)
		return
	}

	var reqProblem Problem
	if err := json.NewDecoder(r.Body).Decode(&reqProblem); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	var author string
	err := h.db.Model(&models.Problem{}).Select("author").Where("id = ?", reqProblem.Id).Scan(&author).Error
	if err != nil {
		log.Println("Error checking problem author:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to verify problem author", nil)
		return
	}

	if author != payload.Sub {
		utils.SendResponse(w, http.StatusForbidden, "You can only update problems you've created", nil)
		return
	}

	updateData := map[string]interface{}{
		"title":                reqProblem.Title,
		"statement":            reqProblem.Statement,
		"input_statement":      reqProblem.InputStatement,
		"output_statement":     reqProblem.OutputStatement,
		"time_limit":           float64(reqProblem.TimeLimit),
		"memory_limit":         float64(reqProblem.MemoryLimit),
		"checker_type":         reqProblem.CheckerType,
		"checker_strict_space": reqProblem.CheckerStrictSpace,
		"checker_precision":    reqProblem.CheckerPrecision,
	}

	result := h.db.Model(&models.Problem{}).Where("id = ?", reqProblem.Id).Updates(updateData)
	if result.Error != nil {
		log.Println("Error updating problem:", result.Error)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to update problem", nil)
		return
	}

	r.SetPathValue("problemId", reqProblem.Id)
	h.GetProblem(w, r)
}