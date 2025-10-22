package users

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"github.com/judgenot0/judge-backend/config"
)

type UserCreds struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password"`
}

type Payload struct {
	Sub            int64   `json:"sub"`
	FullName       string  `json:"full_name"`
	Username       string  `json:"username"`
	Role           string  `json:"role"`
	Clan           *string `json:"clan,omitempty"`
	RoomNo         *string `json:"room_no,omitempty"`
	PcNo           *string `json:"pc_no,omitempty"`
	AllowedContest *int64  `json:"allowed_contest,omitempty"`
	AccessToken    string  `json:"access_token"`
	jwt.RegisteredClaims
}

type Handler struct {
	config *config.Config
	db     *sqlx.DB
}

type UserResponse struct {
	Id       int64   `json:"id" db:"id"`
	FullName string  `json:"full_name" db:"full_name"`
	Username string  `json:"username" db:"username"`
	Clan     *string `json:"clan,omitempty" db:"clan"`
}

func NewHandler(config *config.Config, db *sqlx.DB) *Handler {
	return &Handler{
		config: config,
		db:     db,
	}
}
