package problem

import (
	"log/slog"
	"net/http"

	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) DeleteTestcase(w http.ResponseWriter, r *http.Request) {
	testcaseId := r.PathValue("testcaseId")
	if testcaseId == "" {
		slog.Warn("missing testcase ID")
		utils.SendResponse(w, http.StatusBadRequest, "Invalid testcase ID", nil)
		return
	}

	result := h.db.Delete(&models.Testcase{}, "id = ?", testcaseId)
	if result.Error != nil {
		slog.Error("Error deleting testcase", "error", result.Error)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to delete testcase", nil)
		return
	}

	if result.RowsAffected == 0 {
		utils.SendResponse(w, http.StatusNotFound, "Testcase not found", nil)
		return
	}

	utils.SendResponse(w, http.StatusOK, nil, map[string]string{"message": "Testcase deleted successfully"})
}
