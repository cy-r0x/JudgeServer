package users

import (
	"time"

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

type UserResponse struct {
	Id       int64   `json:"id" db:"id"`
	FullName string  `json:"full_name" db:"full_name"`
	Username string  `json:"username" db:"username"`
	Clan     *string `json:"clan,omitempty" db:"clan"`
}

type User struct {
	Id             int64     `json:"id" db:"id"`
	FullName       string    `json:"full_name" db:"full_name"`
	Username       string    `json:"username" db:"username"`
	Password       string    `json:"password,omitempty" db:"password"`
	Role           string    `json:"role,omitempty" db:"role"`
	Clan           *string   `json:"clan,omitempty" db:"clan"`
	RoomNo         *string   `json:"room_no,omitempty" db:"room_no"`
	PcNo           *string   `json:"pc_no,omitempty" db:"pc_no"`
	AllowedContest *int64    `json:"allowed_contest,omitempty" db:"allowed_contest"`
	CreatedAt      time.Time `json:"created_at,omitempty" db:"created_at"`
}

type Handler struct {
	config *config.Config
	db     *sqlx.DB
}
