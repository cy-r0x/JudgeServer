package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/judgenot0/judge-backend/utils"
)

type Payload struct {
	Sub      string `json:"sub"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func (m *Middlewares) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			utils.SendResopnse(w, http.StatusUnauthorized, "Dhur hala tui hocker")
			return
		}
		headerArr := strings.Split(header, " ")
		if len(headerArr) != 2 {
			utils.SendResopnse(w, http.StatusUnauthorized, "Token koi beda")
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
			utils.SendResopnse(w, http.StatusUnauthorized, "Invalid Token")
			return
		}

		if !token.Valid {
			utils.SendResopnse(w, http.StatusUnauthorized, "Invalid Token")
			return
		}

		// Store payload in context
		ctx := context.WithValue(r.Context(), "user", payload)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
