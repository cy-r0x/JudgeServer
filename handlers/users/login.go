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
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type Data struct {
	Sub      string `json:"sub"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
	AccessToken string `json:"accessToken"`
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
		Role:     "user",
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
	data := Data{
		Sub:              payload.Sub,
		Username:         payload.Username,
		Role:             payload.Role,
		RegisteredClaims: payload.RegisteredClaims,
		AccessToken:      accessToken,
	}
	utils.SendResopnse(w, http.StatusAccepted, data)
}
