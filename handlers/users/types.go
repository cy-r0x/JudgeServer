package users

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/judgenot0/judge-backend/config"
	"gorm.io/gorm"
)

type UserCreds struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password"`
}

type Payload struct {
	Sub            string  `json:"sub"`
	FullName       string  `json:"full_name"`
	Username       string  `json:"username"`
	Role           string  `json:"role"`
	Clan           *string `json:"clan,omitempty"`
	RoomNo         *string `json:"room_no,omitempty"`
	PcNo           *string `json:"pc_no,omitempty"`
	AllowedContest *string `json:"allowed_contest,omitempty"`
	AccessToken    string  `json:"access_token"`
	jwt.RegisteredClaims
}

type UserResponse struct {
	Id       string  `json:"id" db:"id"`
	FullName string  `json:"full_name" db:"full_name"`
	Username string  `json:"username" db:"username"`
	Clan     *string `json:"clan,omitempty" db:"clan"`
}

type User struct {
	Id             string    `json:"id" db:"id"`
	FullName       string    `json:"full_name" db:"full_name"`
	Username       string    `json:"username" db:"username"`
	Password       string    `json:"password,omitempty" db:"password"`
	Role           string    `json:"role,omitempty" db:"role"`
	Clan           *string   `json:"clan,omitempty" db:"clan"`
	RoomNo         *string   `json:"room_no,omitempty" db:"room_no"`
	PcNo           *string   `json:"pc_no,omitempty" db:"pc_no"`
	AllowedContest *string   `json:"allowed_contest,omitempty" db:"allowed_contest"`
	CreatedAt      time.Time `json:"created_at,omitempty" db:"created_at"`
}

type Handler struct {
	config *config.Config
	db     *gorm.DB
}
