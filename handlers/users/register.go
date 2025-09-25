package users

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/judgenot0/judge-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       string `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
	Role     string `json:"role" db:"role"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.SendResopnse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// basic validation
	if user.Username == "" || user.Email == "" || user.Password == "" {
		utils.SendResopnse(w, http.StatusBadRequest, "username, email and password are required")
		return
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		utils.SendResopnse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	user.Password = string(hashedPassword)

	// default role if not provided
	if user.Role == "" {
		user.Role = "user"
	}

	// insert into DB
	query := `
		INSERT INTO users (username, email, password, role)
		VALUES (:username, :email, :password, :role)
		RETURNING id;
	`

	rows, err := h.db.NamedQuery(query, user)
	if err != nil {
		log.Println(err)
		utils.SendResopnse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&user.Id); err != nil {
			log.Println(err)
			utils.SendResopnse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	// don’t send password back
	user.Password = ""

	utils.SendResopnse(w, http.StatusCreated, user)
}
