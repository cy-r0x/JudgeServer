package usercsv

import (
	"net/http"

	"github.com/judgenot0/judge-backend/models"
)

func (h *Handler) GetCSV(w http.ResponseWriter, r *http.Request) {
	contestId := r.URL.Query().Get("contestId")
	if contestId == "" {
		http.Error(w, "contest_id is required", http.StatusBadRequest)
		return
	}

	var filepathObj models.Filepath
	err := h.db.Where("contest_id = ?", contestId).First(&filepathObj).Error
	if err != nil {
		http.Error(w, "error fetching file path", http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, filepathObj.FilePath)
}
