package problem

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/judgenot0/judge-backend/middlewares"
	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) GetProblem(w http.ResponseWriter, r *http.Request) {
	payload, ok := r.Context().Value("user").(*middlewares.Payload)
	if !ok {
		utils.SendResponse(w, http.StatusUnauthorized, "User information not found", nil)
		return
	}

	problemId := r.PathValue("problemId")
	if problemId == "" {
		utils.SendResponse(w, http.StatusBadRequest, "Problem ID is required", nil)
		return
	}

	var problem Problem
	isSampleOnly := false

	switch payload.Role {
	case "user":
		if payload.AllowedContest == nil || *payload.AllowedContest == "" {
			utils.SendResponse(w, http.StatusForbidden, "You don't have access to this problem", nil)
			return
		}

		var count int64
		err := h.db.Model(&models.ContestProblem{}).Where("problem_id = ? AND contest_id = ?", problemId, *payload.AllowedContest).Count(&count).Error
		if err != nil {
			log.Println("Error checking problem access:", err)
			utils.SendResponse(w, http.StatusInternalServerError, "Failed to verify problem access", nil)
			return
		}

		if count == 0 {
			utils.SendResponse(w, http.StatusForbidden, "You don't have access to this problem", nil)
			return
		}
		isSampleOnly = true
	case "setter":
		var author string
		err := h.db.Model(&models.Problem{}).Select("author").Where("id = ?", problemId).Scan(&author).Error
		if err != nil {
			log.Println("Error checking problem author:", err)
			utils.SendResponse(w, http.StatusInternalServerError, "Failed to verify problem author", nil)
			return
		}

		if author != payload.Sub {
			utils.SendResponse(w, http.StatusForbidden, "You don't have access to this problem", nil)
			return
		}
	case "admin":
	default:
		utils.SendResponse(w, http.StatusForbidden, "Invalid role", nil)
		return
	}

	if payload.Role == "user" {
		type DBResult struct {
			models.Problem
			StartTime       *time.Time `gorm:"column:start_time"`
			DurationSeconds *int64     `gorm:"column:duration_seconds"`
		}
		var result DBResult

		err := h.db.Raw(`
			SELECT p.*, c.start_time, c.duration_seconds
			FROM problems p
			INNER JOIN contest_problems cp ON p.id = cp.problem_id
			INNER JOIN contests c ON cp.contest_id = c.id
			WHERE p.id = ?`, problemId).Scan(&result).Error

		if err != nil || result.Id == "" {
			log.Println("Error fetching problem:", err)
			utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch problem", nil)
			return
		}

		problem = Problem{
			Id:                 result.Id,
			Title:              result.Title,
			Statement:          result.Statement,
			InputStatement:     result.InputStatement,
			OutputStatement:    result.OutputStatement,
			TimeLimit:          float32(result.TimeLimit),
			MemoryLimit:        float32(result.MemoryLimit),
			CheckerType:        result.CheckerType,
			CheckerStrictSpace: result.CheckerStrictSpace,
			CheckerPrecision:   result.CheckerPrecision,
			StartTime:          result.StartTime,
			DurationSeconds:    result.DurationSeconds,
			Author:             result.Author,
			CreatedAt:          result.CreatedAt,
			UpdatedAt:          result.UpdatedAt,
		}

	} else {
		var dbProblem models.Problem
		err := h.db.Where("id = ?", problemId).First(&dbProblem).Error
		if err != nil {
			log.Println("Error fetching problem:", err)
			utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch problem", nil)
			return
		}

		problem = Problem{
			Id:                 dbProblem.Id,
			Title:              dbProblem.Title,
			Statement:          dbProblem.Statement,
			InputStatement:     dbProblem.InputStatement,
			OutputStatement:    dbProblem.OutputStatement,
			TimeLimit:          float32(dbProblem.TimeLimit),
			MemoryLimit:        float32(dbProblem.MemoryLimit),
			CheckerType:        dbProblem.CheckerType,
			CheckerStrictSpace: dbProblem.CheckerStrictSpace,
			CheckerPrecision:   dbProblem.CheckerPrecision,
			Author:             dbProblem.Author,
			CreatedAt:          dbProblem.CreatedAt,
			UpdatedAt:          dbProblem.UpdatedAt,
		}
	}

	var wg sync.WaitGroup
	var testcases []Testcase
	var lastSubmission *LastSubmissionData

	wg.Add(1)
	go func() {
		defer wg.Done()
		tc, tcErr := h.fetchTestcases(problemId, isSampleOnly)
		if tcErr != nil {
			log.Println("Error fetching testcases:", tcErr)
		} else {
			testcases = tc
		}
	}()

	if payload.Role == "user" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var ls LastSubmissionData
			err := h.db.Model(&models.Submission{}).
				Select("source_code", "language").
				Where("user_id = ? AND problem_id = ?", payload.Sub, problemId).
				Order("created_at DESC").
				Limit(1).
				Scan(&ls).Error

			if err == nil && ls.SourceCode != "" {
				lastSubmission = &ls
			}
		}()
	}

	wg.Wait()

	if testcases == nil {
		testcases = []Testcase{}
	}
	// problem.Testcases = testcases
	problem.LastSubmission = lastSubmission

	utils.SendResponse(w, http.StatusOK, nil, problem)
}
