package usercsv

import (
	"fmt"
)

func (h *Handler) GenerateUser(record []string, idx int) (User, error) {
	if len(record) < 4 {
		return User{}, fmt.Errorf("record length mismatch") // skip invalid rows
	}

	username := fmt.Sprintf("%s_%d", h.Writer.prefix, idx+1)
	password := generatePassword() // Generate password once

	return User{
		FullName:         record[0],
		Clan:             &record[1],
		RoomNo:           &record[2],
		PcNo:             &record[3],
		AllowedContest:   h.Writer.contestId,
		Username:         username,
		Password:         "",
		UnHashedPassword: password,
		Role:             "user",
	}, nil
}
