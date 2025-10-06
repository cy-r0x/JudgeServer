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
		utils.SendResponse(w, http.StatusBadRequest, "Bad Request")
		return
	}

	// fetch user from DB
	var dbUser User
	query := `SELECT id, full_name, username, password, role, room_no, pc_no, allowed_contest FROM users WHERE username=$1 LIMIT 1`
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

	// Set the JWT token cookie for authentication
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true, // keep it secure from JS
		SameSite: http.SameSiteLaxMode,
		Secure:   false, // set to true for HTTPS in production
		Expires:  time.Now().Add(3 * time.Hour),
	})
	// Will Implement Later apatoto ei thik ache -.-
	// userInfo, _ := json.Marshal(map[string]any{
	// 	"id":              payload.Sub,
	// 	"username":        payload.Username,
	// 	"full_name":       payload.FullName,
	// 	"role":            payload.Role,
	// 	"room_no":         payload.RoomNo,
	// 	"pc_no":           payload.PcNo,
	// 	"allowed_contest": payload.AllowedContest,
	// })

	// log.Println(string(userInfo))

	// // Encode the JSON as base64 to make it safe for cookie values
	// encodedUserInfo := base64.StdEncoding.EncodeToString(userInfo)

	// http.SetCookie(w, &http.Cookie{
	// 	Name:     "user",
	// 	Value:    encodedUserInfo,
	// 	Path:     "/",
	// 	HttpOnly: false, // allow JS access for user info
	// 	SameSite: http.SameSiteLaxMode,
	// 	Secure:   false, // set to true for HTTPS in production
	// 	Expires:  time.Now().Add(3 * time.Hour),
	// })

	// success response
	utils.SendResponse(w, http.StatusOK, payload)
}
