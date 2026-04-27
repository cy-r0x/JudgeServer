package problem

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

type UpdateTestcasePayload struct {
	Input          *string `json:"input,omitempty"`
	ExpectedOutput *string `json:"expected_output,omitempty"`
	IsSample       *bool   `json:"is_sample,omitempty"`
}

func (h *Handler) UpdateTestcase(w http.ResponseWriter, r *http.Request) {
	testcaseId := r.PathValue("testcaseId")
	if testcaseId == "" {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid testcase ID", nil)
		return
	}

	var payload UpdateTestcasePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	updates := map[string]interface{}{}
	if payload.Input != nil {
		updates["input"] = *payload.Input
	}
	if payload.ExpectedOutput != nil {
		updates["expected_output"] = *payload.ExpectedOutput
	}
	if payload.IsSample != nil {
		updates["is_sample"] = *payload.IsSample
	}

	if len(updates) == 0 {
		utils.SendResponse(w, http.StatusBadRequest, "No fields to update", nil)
		return
	}

	result := h.db.Model(&models.Testcase{}).Where("id = ?", testcaseId).Updates(updates)
	if result.Error != nil {
		log.Println("Error updating testcase:", result.Error)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to update testcase", nil)
		return
	}

	if result.RowsAffected == 0 {
		utils.SendResponse(w, http.StatusNotFound, "Testcase not found", nil)
		return
	}

	var updated models.Testcase
	if err := h.db.Where("id = ?", testcaseId).First(&updated).Error; err != nil {
		log.Println("Error fetching updated testcase:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch updated testcase", nil)
		return
	}

	utils.SendResponse(w, http.StatusOK, nil, Testcase{
		Id:             updated.Id,
		ProblemId:      updated.ProblemID,
		Input:          updated.Input,
		ExpectedOutput: updated.ExpectedOutput,
		IsSample:       updated.IsSample,
		CreatedAt:      updated.CreatedAt,
	})
}
