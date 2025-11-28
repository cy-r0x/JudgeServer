package users

import (
	"net/http"
	"strconv"

	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	contestId := r.PathValue("contestId")

	// Convert contestId to int64
	contestIdInt, err := strconv.ParseInt(contestId, 10, 64)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid contest ID")
		return
	}

	// Query to get users where allowed_contest matches contestId
	query := `SELECT id, full_name, username, clan, room_no, pc_no FROM users WHERE allowed_contest = $1`

	var users []User
	err = h.db.Select(&users, query, contestIdInt)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	utils.SendResponse(w, http.StatusOK, users)
}
