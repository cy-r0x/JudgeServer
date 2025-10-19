package problem

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) AddTestCase(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var testcase Testcase
	if err := decoder.Decode(&testcase); err != nil {
		log.Println("Error decoding request body:", err)
		utils.SendResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if testcase.ProblemId == 0 {
		utils.SendResponse(w, http.StatusBadRequest, "Problem ID is required")
		return
	}

	query := `INSERT INTO testcases 
				(problem_id, input, expected_output, is_sample)
				VALUES ($1, $2, $3, $4)
				RETURNING id`

	tx, err := h.db.Beginx()
	if err != nil {
		log.Println("Error starting transaction:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to create testcase")
		return
	}

	var testcase_id int64
	err = tx.QueryRow(query, testcase.ProblemId, testcase.Input, testcase.ExpectedOutput, testcase.IsSample).Scan(&testcase_id)
	if err != nil {
		tx.Rollback()
		log.Println("Error creating testcase:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to create testcase")
		return
	}

	if err := tx.Commit(); err != nil {
		log.Println("Error committing transaction:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to create testcase")
		return
	}

	testcase.Id = testcase_id
	utils.SendResponse(w, http.StatusOK, testcase)
}
