package problem

import (
	"log"
	"net/http"
	"strconv"

	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) DeleteTestcase(w http.ResponseWriter, r *http.Request) {
	testcaseId, err := strconv.ParseInt(r.PathValue("testcaseId"), 10, 64)
	if err != nil {
		log.Println(err)
		utils.SendResponse(w, http.StatusBadRequest, "Invalid testcase ID")
		return
	}

	result := h.db.Delete(&models.Testcase{}, testcaseId)
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
