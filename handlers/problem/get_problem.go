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
		utils.SendResponse(w, http.StatusUnauthorized, "User information not found")
		return
	}

	problemId := r.PathValue("problemId")
	if problemId == "" {
		utils.SendResponse(w, http.StatusBadRequest, "Problem ID is required")
		return
	}

	var problem Problem
	isSampleOnly := false

	switch payload.Role {
	case "user":
		if payload.AllowedContest == nil || *payload.AllowedContest == "" {
			utils.SendResponse(w, http.StatusForbidden, "You don't have access to this problem")
			return
		}

		// Check if the user has access to this problem through their allowed contest
		var count int64
		err := h.db.Model(&models.ContestProblem{}).Where("problem_id = ? AND contest_id = ?", problemId, *payload.AllowedContest).Count(&count).Error
		if err != nil {
			log.Println("Error checking problem access:", err)
			utils.SendResponse(w, http.StatusInternalServerError, "Failed to verify problem access")
			return
		}

		if count == 0 {
			utils.SendResponse(w, http.StatusForbidden, "You don't have access to this problem")
			return
		}
		isSampleOnly = true
		// Fall through to fetch problem data
	case "setter":
		// Check if the problem was created by this setter
		var createdBy *string
		err := h.db.Model(&models.Problem{}).Select("created_by").Where("id = ?", problemId).Scan(&createdBy).Error
		if err != nil {
			log.Println("Error checking problem creator:", err)
			utils.SendResponse(w, http.StatusInternalServerError, "Failed to verify problem creator")
			return
		}

		if createdBy == nil || *createdBy != payload.Sub {
			utils.SendResponse(w, http.StatusForbidden, "You don't have access to this problem")
			return
		}

		// Fall through to fetch problem data
	case "admin":
		// Admin has access to all problems
	default:
		utils.SendResponse(w, http.StatusForbidden, "Invalid role")
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

		if err != nil || result.ID == "" {
			log.Println("Error fetching problem:", err)
			utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch problem")
			return
		}

		createdBy := ""
		if result.CreatedByID != nil {
			createdBy = *result.CreatedByID
		}

		problem = Problem{
			Id:                 result.ID,
			Title:              result.Title,
			Slug:               result.Slug,
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
			CreatedBy:          createdBy,
			CreatedAt:          result.CreatedAt,
		}

	} else {
		var dbProblem models.Problem
		err := h.db.Where("id = ?", problemId).First(&dbProblem).Error
		if err != nil {
			log.Println("Error fetching problem:", err)
			utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch problem")
			return
		}

		createdBy := ""
		if dbProblem.CreatedByID != nil {
			createdBy = *dbProblem.CreatedByID
		}

		problem = Problem{
			Id:                 dbProblem.ID,
			Title:              dbProblem.Title,
			Slug:               dbProblem.Slug,
			Statement:          dbProblem.Statement,
			InputStatement:     dbProblem.InputStatement,
			OutputStatement:    dbProblem.OutputStatement,
			TimeLimit:          float32(dbProblem.TimeLimit),
			MemoryLimit:        float32(dbProblem.MemoryLimit),
			CheckerType:        dbProblem.CheckerType,
			CheckerStrictSpace: dbProblem.CheckerStrictSpace,
			CheckerPrecision:   dbProblem.CheckerPrecision,
			CreatedBy:          createdBy,
			CreatedAt:          dbProblem.CreatedAt,
		}
	}

	// Fetch testcases and last submission concurrently
	var wg sync.WaitGroup
	var testcases []Testcase
	var lastSubmission *LastSubmissionData

	// Fetch testcases
	wg.Add(1)
	go func() {
		defer wg.Done()
		tc, tcErr := h.fetchTestcases(problemId, isSampleOnly)
		if tcErr != nil {
			log.Println("Error fetching testcases:", tcErr)
			// Continue anyway as we at least have the problem data
		} else {
			testcases = tc
		}
	}()

	// Fetch last submission for users
	if payload.Role == "user" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var ls LastSubmissionData
			err := h.db.Model(&models.Submission{}).
				Select("source_code", "language").
				Where("user_id = ? AND problem_id = ?", payload.Sub, problemId).
				Order("submitted_at DESC").
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
	problem.Testcases = testcases
	problem.LastSubmission = lastSubmission

	utils.SendResponse(w, http.StatusOK, problem)
}
