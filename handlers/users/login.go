package users

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/judgenot0/judge-backend/utils"
)

type UserCreds struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Payload struct {
	Sub      string `json:"sub"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var user UserCreds
	err := decoder.Decode(&user)
	if err != nil {
		utils.SendResopnse(w, http.StatusBadRequest, "Wrong Structure")
		return
	}
	//TODO: check if the usercreds is correct or not
	payload := &Payload{
		Sub:      "123456",
		Username: "Prantor",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Hour)),
		},
	}
	secret := h.config.SecretKey
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	accessToken, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatalln(err)
		return
	}
	utils.SendResopnse(w, http.StatusAccepted, accessToken)
}
