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
	const maxBodySize = 1024 // 1 KB

	// Limit the size of the request body
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	var creds UserCreds
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Bad Request")
		return
	}

	// fetch user from DB
	var dbUser User
	query := `SELECT id, full_name, username, password, role, clan, room_no, pc_no, allowed_contest FROM users WHERE username=$1 LIMIT 1`
	err := h.db.Get(&dbUser, query, creds.Username)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}
	// compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(creds.Password)); err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	// build payload
	payload := &Payload{
		Sub:            dbUser.Id,
		FullName:       dbUser.FullName,
		Username:       dbUser.Username,
		Clan:           dbUser.Clan,
		Role:           dbUser.Role,
		RoomNo:         dbUser.RoomNo,
		PcNo:           dbUser.PcNo,
		AllowedContest: dbUser.AllowedContest,
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
		utils.SendResponse(w, http.StatusInternalServerError, "Could not login")
		return
	}
	payload.AccessToken = accessToken
	// success response
	utils.SendResponse(w, http.StatusOK, payload)
}
