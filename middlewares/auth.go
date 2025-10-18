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
	Role           string  `json:"role"`
	RoomNo         *string `json:"room_no"`
	PcNo           *int    `json:"pc_no"`
	AllowedContest *int64  `json:"allowed_contest"`
	AccessToken    string  `json:"accessToken"`
	jwt.RegisteredClaims
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

		payload := &Payload{}

		token, err := jwt.ParseWithClaims(accessToken, payload, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(m.config.SecretKey), nil
		})

		if err != nil {
			log.Println(err)
			utils.SendResponse(w, http.StatusUnauthorized, "Invalid Token")
			return
		}

		if !token.Valid {
			utils.SendResponse(w, http.StatusUnauthorized, "Invalid Token")
			return
		}

		// Store payload in context
		ctx := context.WithValue(r.Context(), "user", payload)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
