package middlewares

import "github.com/golang-jwt/jwt/v5"

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
