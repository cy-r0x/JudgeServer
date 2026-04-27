package contest_problems

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) DeleteContestProblem(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var contestProblem ContestProblem
	err := decoder.Decode(&contestProblem)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid JSON", nil)
		return
	}
	if contestProblem.ContestId == "" || contestProblem.ProblemId == "" {
		utils.SendResponse(w, http.StatusBadRequest, "Contest ID and Problem ID are required", nil)
		return
	}

	tx := h.db.Begin()
	if tx.Error != nil {
		log.Println("Failed to begin transaction:", tx.Error)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to delete contest problem", nil)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	result := tx.Where("contest_id = ? AND problem_id = ?", contestProblem.ContestId, contestProblem.ProblemId).Delete(&models.ContestProblem{})
	if result.Error != nil {
		tx.Rollback()
		log.Println("Failed to delete contest problem:", result.Error)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to delete contest problem", nil)
		return
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		utils.SendResponse(w, http.StatusNotFound, "Problem not assigned to this contest", nil)
		return
	}

	var remaining []models.ContestProblem
	if err = tx.Where("contest_id = ?", contestProblem.ContestId).Order("index ASC").Find(&remaining).Error; err != nil {
		tx.Rollback()
		log.Println("Failed to fetch remaining contest problems:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to delete contest problem", nil)
		return
	}

	for i, cp := range remaining {
		newIndex := i + 1
		if cp.Index == newIndex {
			continue
		}
		if err = tx.Model(&models.ContestProblem{}).Where("contest_id = ? AND problem_id = ?", cp.ContestID, cp.ProblemID).Update("index", newIndex).Error; err != nil {
			tx.Rollback()
			log.Println("Failed to update contest problem index:", err)
			utils.SendResponse(w, http.StatusInternalServerError, "Failed to delete contest problem", nil)
			return
		}
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Println("Failed to commit contest problem deletion:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to delete contest problem", nil)
		return
	}

	utils.SendResponse(w, http.StatusOK, map[string]any{"message": "Contest problem removed"}, nil)
}
