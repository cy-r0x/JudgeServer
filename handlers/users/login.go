package users

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/judgenot0/judge-backend/utils"
	"golang.org/x/crypto/bcrypt"
)



func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var creds UserCreds
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		utils.SendResopnse(w, http.StatusBadRequest, "Bad Request")
		return
	}

	// fetch user from DB
	var dbUser User
	query := `SELECT id, username, password, role FROM users WHERE username=$1 LIMIT 1`
	err := h.db.Get(&dbUser, query, creds.Username)
	if err != nil {
		utils.SendResopnse(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	// compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(creds.Password)); err != nil {
		utils.SendResopnse(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	// build payload
	payload := &Payload{
		Sub:      dbUser.Id,
		Username: dbUser.Username,
		Role:     dbUser.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// sign token
	secret := h.config.SecretKey
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	accessToken, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Println("error signing jwt:", err)
		utils.SendResopnse(w, http.StatusInternalServerError, "Could not login")
		return
	}
	payload.AccessToken = accessToken

	// success response
	utils.SendResopnse(w, http.StatusOK, payload)
}
