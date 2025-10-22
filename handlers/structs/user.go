package structs

import "time"

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
