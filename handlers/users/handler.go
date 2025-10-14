package users

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"github.com/judgenot0/judge-backend/config"
)

type User struct {
	Id             int64     `json:"id" db:"id"`
	FullName       string    `json:"full_name" db:"full_name"`
	Username       string    `json:"username" db:"username"`
	Email          string    `json:"email" db:"email"`
	Password       string    `json:"password" db:"password"`
	Role           string    `json:"role" db:"role"`
	RoomNo         *string   `json:"room_no" db:"room_no"`
	PcNo           *int      `json:"pc_no" db:"pc_no"`
	AllowedContest *int64    `json:"allowed_contest" db:"allowed_contest"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type UserCreds struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password"`
}

type Payload struct {
	Sub            int64   `json:"sub"`
	FullName       string  `json:"full_name"`
	Username       string  `json:"username"`
	Role           string  `json:"role"`
	RoomNo         *string `json:"room_no"`
	PcNo           *int    `json:"pc_no"`
	AllowedContest *int64  `json:"allowed_contest"`
	AccessToken    string  `json:"access_token"`
	jwt.RegisteredClaims
}

type Handler struct {
	config *config.Config
	db     *sqlx.DB
}

type UserResponse struct {
	UserId   int64  `json:"userId" db:"id"`
	FullName string `json:"full_name" db:"full_name"`
	Username string `json:"username" db:"username"`
}

func NewHandler(config *config.Config, db *sqlx.DB) *Handler {
	return &Handler{
		config: config,
		db:     db,
	}
}
