package users

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/judgenot0/judge-backend/handlers/structs"
	"github.com/judgenot0/judge-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user structs.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// basic validation
	if user.Username == "" || user.Password == "" {
		utils.SendResponse(w, http.StatusBadRequest, "username and password are required")
		return
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		utils.SendResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()

	// default role if not provided
	if user.Role == "" {
		user.Role = "user"
	}

	// insert into DB
	query := `
		INSERT INTO users (full_name, username, password, clan, room_no, pc_no, role, allowed_contest, created_at)
		VALUES (:full_name,:username, :password, :clan, :room_no, :pc_no, :role, :allowed_contest, :created_at);
	`

	rows, err := h.db.NamedQuery(query, user)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&user.Id); err != nil {
			utils.SendResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	// don't send password back
	user.Password = ""

	utils.SendResponse(w, http.StatusCreated, user)
}
