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

type UpdateUserPayload struct {
	Name           *string `json:"name"`
	Password       *string `json:"password"`
	AdditionalInfo *string `json:"additional_info"`
	RoomNo         *string `json:"room_no"`
	PcNo           *string `json:"pc_no"`
	AllowedContest *string `json:"allowed_contest"`
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
	Id             string  `json:"id" db:"id"`
	Name           string  `json:"full_name" db:"full_name"`
	Username       string  `json:"username" db:"username"`
	AdditionalInfo *string `json:"additional_info,omitempty" db:"additional_info"`
}

type Handler struct {
	config *config.Config
	db     *gorm.DB
}
