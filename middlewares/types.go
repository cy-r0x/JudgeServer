package middlewares

import "github.com/golang-jwt/jwt/v5"

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
