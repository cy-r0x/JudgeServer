package compilerun

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/judgenot0/judge-backend/middlewares"
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
		utils.SendResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
		return
	}

	var problem Problem
	err := h.db.Get(&problem, `SELECT time_limit, memory_limit,checker_type, checker_strict_space, checker_precision FROM problems WHERE id=$1`, submission.ProblemId)

	if err != nil {
		log.Printf("Error fetching problem details: %v", err)
		utils.SendResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Problem not found",
		})
		return
	}

	testcases, err := h.fetchTestcases(submission.ProblemId, isSampleOnly)
	if err != nil {
		log.Printf("Error fetching testcases for problem ID %d: %v", submission.ProblemId, err)
		utils.SendResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch testcases for this problem",
		})
		return
	}

	problem.SourceCode = submission.SourceCode
	problem.Language = submission.Language
	problem.Testcases = testcases

	fmt.Println(problem)

	url := h.config.EngineUrl + "/run"

	runReq, err := json.Marshal(&problem)
	if err != nil {
		log.Printf("Error marshaling problem data: %v", err)
		utils.SendResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to prepare execution request",
		})
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(runReq))
	if err != nil {
		// Check if the error is due to payload size limit
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			log.Printf("Payload too large: exceeded %d bytes", maxBodySize)
			utils.SendResponse(w, http.StatusRequestEntityTooLarge, map[string]string{
				"error": "Request payload too large",
			})
			return
		}

		log.Printf("Error forwarding request to engine: %v", err)
		utils.SendResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to connect to execution engine",
		})
		return
	}
	defer resp.Body.Close()

	// Decode response from engine
	var result struct {
		Result string `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Error decoding engine response: %v", err)
		utils.SendResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to parse execution result",
		})
		return
	}

	// Send successful response
	utils.SendResponse(w, http.StatusOK, result)
}
