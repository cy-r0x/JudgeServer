package compilerun

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) CompileRun(w http.ResponseWriter, r *http.Request) {
	const maxBodySize = 50 * 1024 // 50 KB
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	isSampleOnly := false

	payload, ok := r.Context().Value("user").(*middlewares.Payload)
	if ok && payload.Role == "user" {
		isSampleOnly = true
	}

	decoder := json.NewDecoder(r.Body)
	var submission UserSubmission
	if err := decoder.Decode(&submission); err != nil {
		log.Printf("Error decoding request body: %v", err)
		utils.SendResponse(w, http.StatusBadRequest, " Invalid request payload", nil)
		return
	}

	var problem Problem
	err := h.db.Model(&models.Problem{}).
		Select("time_limit", "memory_limit", "checker_type", "checker_strict_space", "checker_precision").
		Where("id = ?", submission.ProblemId).
		Scan(&problem).Error

	if err != nil {
		log.Printf("Error fetching problem details: %v", err)
		utils.SendResponse(w, http.StatusBadRequest, "Problem not found", nil)
		return
	}

	testcases, err := h.fetchTestcases(submission.ProblemId, isSampleOnly)
	if err != nil {
		log.Printf("Error fetching testcases for problem ID %s: %v", submission.ProblemId, err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch testcases for this problem", nil)
		return
	}

	problem.SourceCode = submission.SourceCode
	problem.Language = submission.Language
	problem.Testcases = testcases

	url := h.config.EngineUrl + "/run"

	runReq, err := json.Marshal(&problem)
	if err != nil {
		log.Printf("Error marshaling problem data: %v", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to prepare execution request", nil)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(runReq))
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			log.Printf("Payload too large: exceeded %d bytes", maxBodySize)
			utils.SendResponse(w, http.StatusRequestEntityTooLarge, "Request payload too large", nil)
			return
		} else {
			log.Printf("Error sending request to execution engine: %v", err)
			utils.SendResponse(w, http.StatusInternalServerError, "Failed to execute code", nil)
			return
		}
	}

	defer resp.Body.Close()

	// Decode response from engine
	var result struct {
		Result string `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Error decoding engine response: %v", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to parse execution result", nil)
		return
	}

	// Send successful response
	utils.SendResponse(w, http.StatusOK, "Execution completed successfully", result)
}
