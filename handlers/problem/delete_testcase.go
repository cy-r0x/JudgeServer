package problem

import (
	"log"
	"net/http"
	"strconv"

	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) DeleteTestcase(w http.ResponseWriter, r *http.Request) {
	testcase_id, err := strconv.Atoi(r.PathValue("testcaseId"))
	if err != nil {
		log.Println(err)
		utils.SendResponse(w, http.StatusBadRequest, "Invalid testcase ID")
		return
	}

	query := `DELETE FROM testcases WHERE id=$1`

	tx, err := h.db.Beginx()
	if err != nil {
		log.Println("Error starting transaction:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to delete testcase")
		return
	}

	result, err := tx.Exec(query, testcase_id)
	if err != nil {
		tx.Rollback()
		log.Println("Error deleting testcase:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to delete testcase")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		log.Println("Error verifying deletion:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to verify deletion")
		return
	}

	if rowsAffected == 0 {
		tx.Rollback()
		utils.SendResponse(w, http.StatusNotFound, "Testcase not found")
		return
	}

	if err := tx.Commit(); err != nil {
		log.Println("Error committing transaction:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to delete testcase")
		return
	}

	utils.SendResponse(w, http.StatusOK, map[string]string{"message": "Testcase deleted successfully"})
}
