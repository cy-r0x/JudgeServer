package usercsv

import "net/http"

func (h *Handler) GetCSV(w http.ResponseWriter, r *http.Request) {
	contestId := r.URL.Query().Get("contestId")
	if contestId == "" {
		http.Error(w, "contest_id is required", http.StatusBadRequest)
		return
	}

	query := `SELECT file_path FROM filepath WHERE contest_id = $1`
	var filePath string
	err := h.db.QueryRow(query, contestId).Scan(&filePath)
	if err != nil {
		http.Error(w, "error fetching file path", http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, filePath)
}
