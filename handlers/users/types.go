package users

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/judgenot0/judge-backend/config"
	"gorm.io/gorm"
)

type UserCreds struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password"`
}

type Payload struct {
	jwt.RegisteredClaims
	Sub            string  `json:"sub"`
	Name           string  `json:"full_name"`
	Username       string  `json:"username"`
	Role           string  `json:"role"`
	AdditionalInfo *string `json:"additional_info,omitempty"`
	RoomNo         *string `json:"room_no,omitempty"`
	PcNo           *string `json:"pc_no,omitempty"`
	AllowedContest *string `json:"allowed_contest,omitempty"`
}

type UserResponse struct {
	Id       string  `json:"id" db:"id"`
	FullName string  `json:"full_name" db:"full_name"`
	Username string  `json:"username" db:"username"`
	Clan     *string `json:"clan,omitempty" db:"clan"`
}

type Handler struct {
	config *config.Config
	db     *gorm.DB
}
