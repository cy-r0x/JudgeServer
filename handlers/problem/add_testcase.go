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
		utils.SendResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if testcase.ProblemId == 0 {
		utils.SendResponse(w, http.StatusBadRequest, "Problem ID is required")
		return
	}

	newTestcase := models.Testcase{
		ProblemID:      uint(testcase.ProblemId),
		Input:          testcase.Input,
		ExpectedOutput: testcase.ExpectedOutput,
		IsSample:       testcase.IsSample,
	}

	err := h.db.Create(&newTestcase).Error
	if err != nil {
		log.Println("Error creating testcase:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to create testcase")
		return
	}

	testcase.Id = int64(newTestcase.ID)
	testcase.CreatedAt = newTestcase.CreatedAt

	utils.SendResponse(w, http.StatusOK, testcase)
}
