package problem

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) AddTestCase(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var testcase Testcase
	if err := decoder.Decode(&testcase); err != nil {
		log.Println("Error decoding request body:", err)
		utils.SendResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if testcase.ProblemId == "" {
		utils.SendResponse(w, http.StatusBadRequest, "Problem ID is required", nil)
		return
	}

	newTestcase := models.Testcase{
		ProblemID:      testcase.ProblemId,
		Input:          testcase.Input,
		ExpectedOutput: testcase.ExpectedOutput,
		IsSample:       testcase.IsSample,
	}

	err := h.db.Create(&newTestcase).Error
	if err != nil {
		log.Println("Error creating testcase:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to create testcase", nil)
		return
	}

	testcase.Id = newTestcase.Id
	testcase.CreatedAt = newTestcase.CreatedAt

	utils.SendResponse(w, http.StatusOK, nil, testcase)
}
