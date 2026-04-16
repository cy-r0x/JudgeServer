package problem

import (
	"log"
	"net/http"

	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) DeleteTestcase(w http.ResponseWriter, r *http.Request) {
	testcaseId := r.PathValue("testcaseId")
	if testcaseId == "" {
		log.Println("missing testcase ID")
		utils.SendResponse(w, http.StatusBadRequest, "Invalid testcase ID")
		return
	}

	result := h.db.Delete(&models.Testcase{}, "id = ?", testcaseId)
	if result.Error != nil {
		log.Println("Error deleting testcase:", result.Error)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to delete testcase")
		return
	}

	if result.RowsAffected == 0 {
		utils.SendResponse(w, http.StatusNotFound, "Testcase not found")
		return
	}

	utils.SendResponse(w, http.StatusOK, map[string]string{"message": "Testcase deleted successfully"})
}
