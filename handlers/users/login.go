package users

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	const maxBodySize = 1024 // 1 KB

	// Limit the size of the request body
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	var creds UserCreds
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Bad Request", nil)
		return
	}

	// fetch user from DB
	var dbUser models.User
	err := h.db.
		Where("username = ?", creds.Username).
		First(&dbUser).Error

	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, "Invalid username or password", nil)
		return
	}
	// compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(creds.Password)); err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, "Invalid username or password", nil)
		return
	}

	var allowedContest *string
	if dbUser.AllowedContest != nil {
		ac := *dbUser.AllowedContest
		allowedContest = &ac
	}

	// build payload
	payload := &Payload{
		Sub:            dbUser.Id,
		Name:           dbUser.Name,
		Username:       dbUser.Username,
		AdditionalInfo: dbUser.AdditionalInfo,
		Role:           string(dbUser.Role),
		RoomNo:         dbUser.RoomNo,
		PcNo:           dbUser.PcNo,
		AllowedContest: allowedContest,
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
		slog.Error("error signing jwt", "error", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Login failed", nil)
		return
	}

	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(3 * time.Hour),
		HttpOnly: true,
		Secure:   true, // Set to true if running over HTTPS
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)

	// Build user response (without password)
	userResponse := map[string]any{
		"id":              dbUser.Id,
		"fullName":        dbUser.Name,
		"username":        dbUser.Username,
		"role":            string(dbUser.Role),
		"additionalInfo":  dbUser.AdditionalInfo,
		"roomNo":          dbUser.RoomNo,
		"pcNo":            dbUser.PcNo,
		"allowedContest":  dbUser.AllowedContest,
		"accessToken":     accessToken,
	}

	// success response
	utils.SendResponse(w, http.StatusOK, "Login Successful", userResponse)
}
