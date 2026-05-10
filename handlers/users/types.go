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
	Name           *string `json:"fullName"`
	Password       *string `json:"password"`
	AdditionalInfo *string `json:"additionalInfo"`
	RoomNo         *string `json:"roomNo"`
	PcNo           *string `json:"pcNo"`
	AllowedContest *string `json:"allowedContest"`
}

type Payload struct {
	jwt.RegisteredClaims
	Sub            string  `json:"sub"`
	Name           string  `json:"fullName"`
	Username       string  `json:"username"`
	Role           string  `json:"role"`
	AdditionalInfo *string `json:"additionalInfo,omitempty"`
	RoomNo         *string `json:"roomNo,omitempty"`
	PcNo           *string `json:"pcNo,omitempty"`
	AllowedContest *string `json:"allowedContest,omitempty"`
}

type UserResponse struct {
	Id             string  `json:"id" db:"id"`
	Name           string  `json:"fullName" db:"full_name"`
	Username       string  `json:"username" db:"username"`
	AdditionalInfo *string `json:"additionalInfo,omitempty" db:"additional_info"`
}

type Handler struct {
	config *config.Config
	db     *gorm.DB
}
