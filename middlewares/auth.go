package middlewares

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/judgenot0/judge-backend/utils"
)

type Payload struct {
	Sub            int64   `json:"sub"`
	FullName       string  `json:"full_name"`
	Username       string  `json:"username"`
	Clan           *string `json:"clan"`
	Role           string  `json:"role"`
	RoomNo         *string `json:"room_no"`
	PcNo           *string `json:"pc_no"`
	AllowedContest *int64  `json:"allowed_contest"`
	AccessToken    string  `json:"accessToken"`
	jwt.RegisteredClaims
}

func DecodeToken(tokenStr string, secretKey string) (*Payload, error) {
	payload := &Payload{}
	token, err := jwt.ParseWithClaims(tokenStr, payload, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return payload, nil
}

func (m *Middlewares) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			utils.SendResponse(w, http.StatusUnauthorized, "Authorization header required")
			return
		}
		headerArr := strings.Split(header, " ")
		if len(headerArr) != 2 {
			utils.SendResponse(w, http.StatusUnauthorized, "Invalid token format")
			return
		}

		accessToken := headerArr[1]

		payload, err := DecodeToken(accessToken, m.config.SecretKey)

		if err != nil {
			log.Println(err)
			utils.SendResponse(w, http.StatusUnauthorized, "Invalid Token")
			return
		}

		// Store payload in context
		ctx := context.WithValue(r.Context(), "user", payload)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
