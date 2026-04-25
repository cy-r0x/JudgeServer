package usercsv

import (
	"encoding/csv"
	"fmt"
	"net/http"

	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) GetCSV(w http.ResponseWriter, r *http.Request) {
	contestId := r.URL.Query().Get("contestId")
	if contestId == "" {
		utils.SendResponse(w, http.StatusBadRequest, "Contest ID is required", nil)
		return
	}

	var creds []models.UserCreds
	err := h.db.Preload("User").Where("contest_id = ?", contestId).Find(&creds).Error
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "Error fetching credentials", nil)
		return
	}

	if len(creds) == 0 {
		utils.SendResponse(w, http.StatusNotFound, "No credentials found for this contest", nil)
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"contest_%s_credentials.csv\"", contestId))

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write header
	header := []string{"Name", "Username", "Password", "RoomNo", "PcNo"}
	if err := writer.Write(header); err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "Error writing CSV header", nil)
		return
	}

	// Write data
	for _, cred := range creds {
		roomNo := ""
		if cred.User.RoomNo != nil {
			roomNo = *cred.User.RoomNo
		}
		pcNo := ""
		if cred.User.PcNo != nil {
			pcNo = *cred.User.PcNo
		}

		row := []string{
			cred.User.Name,
			cred.User.Username,
			cred.PlainPassword,
			roomNo,
			pcNo,
		}
		if err := writer.Write(row); err != nil {
			utils.SendResponse(w, http.StatusInternalServerError, "Error writing CSV data", nil)
			return
		}
	}
}
