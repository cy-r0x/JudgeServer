package users

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"github.com/judgenot0/judge-backend/config"
)

type User struct {
	Id       string `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
	Role     string `json:"role" db:"role"`
}

type UserCreds struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password"`
}

type Payload struct {
	Sub         string `json:"sub"`
	Username    string `json:"username"`
	Role        string `json:"role"`
	AccessToken string `json:"accessToken"`
	jwt.RegisteredClaims
}

type Handler struct {
	config *config.Config
	db     *sqlx.DB
}

func NewHandler(config *config.Config, db *sqlx.DB) *Handler {
	return &Handler{
		config: config,
		db:     db,
	}
}
