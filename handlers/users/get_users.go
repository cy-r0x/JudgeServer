package users

import (
	"net/http"
	"strconv"

	"github.com/judgenot0/judge-backend/utils"
)

type UserResponse struct {
	UserId   int64  `json:"userId" db:"id"`
	FullName string `json:"full_name" db:"full_name"`
	Username string `json:"username" db:"username"`
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	contestId := r.PathValue("contestId")

	// Convert contestId to int64
	contestIdInt, err := strconv.ParseInt(contestId, 10, 64)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid contest ID")
		return
	}

	// Query to get users where allowed_contest matches contestId
	query := `SELECT id, full_name, username FROM users WHERE allowed_contest = $1`

	var users []UserResponse
	err = h.db.Select(&users, query, contestIdInt)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	utils.SendResponse(w, http.StatusOK, users)
}
