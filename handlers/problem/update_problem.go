package problem

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/judgenot0/judge-backend/middlewares"
	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) UpdateProblem(w http.ResponseWriter, r *http.Request) {
	// Only Setter can update the problem
	payload, ok := r.Context().Value("user").(*middlewares.Payload)
	if !ok {
		utils.SendResponse(w, http.StatusUnauthorized, "User information not found")
		return
	}

	var reqProblem Problem
	if err := json.NewDecoder(r.Body).Decode(&reqProblem); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// For setters, verify they created this problem
	var createdBy *uint
	err := h.db.Model(&models.Problem{}).Select("created_by").Where("id = ?", reqProblem.Id).Scan(&createdBy).Error
	if err != nil {
		log.Println("Error checking problem creator:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to verify problem creator")
		return
	}

	if createdBy == nil || *createdBy != uint(payload.Sub) {
		utils.SendResponse(w, http.StatusForbidden, "You can only update problems you've created")
		return
	}

	slug := strings.ReplaceAll(strings.ToLower(reqProblem.Title), " ", "-")

	updateData := map[string]interface{}{
		"title":                reqProblem.Title,
		"slug":                 slug,
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
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to update problem")
		return
	}

	r.SetPathValue("problemId", strconv.FormatInt(reqProblem.Id, 10))
	h.GetProblem(w, r)
}
