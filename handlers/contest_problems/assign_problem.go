package contest_problems

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) AssignContestProblems(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var contestProblem ContestProblem
	err := decoder.Decode(&contestProblem)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	if contestProblem.ContestId == 0 || contestProblem.ProblemId == 0 {
		utils.SendResponse(w, http.StatusBadRequest, "Contest ID and Problem ID are required")
		return
	}

	// Check if contest exists
	var countContest int64
	if err = h.db.Model(&models.Contest{}).Where("id = ?", contestProblem.ContestId).Count(&countContest).Error; err != nil {
		log.Println("Failed to check contest existence:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to assign contest problem")
		return
	}
	if countContest == 0 {
		utils.SendResponse(w, http.StatusBadRequest, "Contest does not exist")
		return
	}

	// Get problem details and author information in a single optimized query using GORM Joins
	type ProblemDetails struct {
		Title    string `gorm:"column:title"`
		FullName string `gorm:"column:full_name"`
	}

	var problemDetails ProblemDetails
	query := `
		SELECT p.title, u.full_name 
		FROM problems p 
		LEFT JOIN users u ON p.created_by = u.id 
		WHERE p.id = ?
	`

	if err = h.db.Raw(query, contestProblem.ProblemId).Scan(&problemDetails).Error; err != nil {
		log.Println("Failed to get problem details:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to assign contest problem")
		return
	}
	if problemDetails.Title == "" {
		utils.SendResponse(w, http.StatusBadRequest, "Problem does not exist")
		return
	}

	// Set problem details in the response struct
	contestProblem.ProblemName = problemDetails.Title
	contestProblem.ProblemAuthor = problemDetails.FullName

	tx := h.db.Begin()
	if tx.Error != nil {
		log.Println("Failed to begin transaction:", tx.Error)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to assign contest problem")
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Ensure the problem is not already assigned
	var existsCount int64
	if err = tx.Model(&models.ContestProblem{}).Where("contest_id = ? AND problem_id = ?", contestProblem.ContestId, contestProblem.ProblemId).Count(&existsCount).Error; err != nil {
		log.Println("Failed to check existing assignment:", err)
		tx.Rollback()
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to assign contest problem")
		return
	}
	if existsCount > 0 {
		tx.Rollback()
		utils.SendResponse(w, http.StatusConflict, "Problem already assigned to this contest")
		return
	}

	// Count existing problems to determine next index
	var count int64
	if err = tx.Model(&models.ContestProblem{}).Where("contest_id = ?", contestProblem.ContestId).Count(&count).Error; err != nil {
		log.Println("Failed to count contest problems:", err)
		tx.Rollback()
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to assign contest problem")
		return
	}

	contestProblem.Index = int(count) + 1

	newCP := models.ContestProblem{
		ContestID: uint(contestProblem.ContestId),
		ProblemID: uint(contestProblem.ProblemId),
		Index:     contestProblem.Index,
	}

	if err = tx.Create(&newCP).Error; err != nil {
		log.Println("Failed to insert contest problem:", err)
		tx.Rollback()
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to assign contest problem")
		return
	}

	if err = tx.Commit().Error; err != nil {
		log.Println("Failed to commit contest problem assignment:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to assign contest problem")
		return
	}

	utils.SendResponse(w, http.StatusOK, contestProblem)
}
