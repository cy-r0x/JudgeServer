package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/judgenot0/judge-backend/config"
	"github.com/judgenot0/judge-backend/utils"
)

type Payload struct {
	Sub      string `json:"sub"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Authenticate(next http.Handler) http.Handler {
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
		config, err := config.GetConfig() //dependecy kemne shorabo :"_) !
		if err != nil {
			utils.SendResopnse(w, http.StatusInternalServerError, err.Error())
			return
		}

		token, err := jwt.ParseWithClaims(accessToken, payload, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(config.SecretKey), nil
		})

		if err != nil {
			utils.SendResopnse(w, http.StatusUnauthorized, err.Error())
			return
		}

		if !token.Valid {
			utils.SendResopnse(w, http.StatusUnauthorized, "Invalid Token")
			return
		}
		fmt.Println("Username from claims:", payload.Username)
		fmt.Println("ExpiresAt:", payload.ExpiresAt)
		next.ServeHTTP(w, r)
	})
}
