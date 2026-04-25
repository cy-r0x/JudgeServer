package usercsv

import (
	"fmt"
)

func (h *Handler) GenerateUser(record []string, idx int, prefix string, contestId string) (User, error) {
	if len(record) < 4 {
		return User{}, fmt.Errorf("record length mismatch") // skip invalid rows
	}

	username := fmt.Sprintf("%s_%d", prefix, idx+1)
	password := generatePassword() // Generate password once

	return User{
		Name:             record[0],
		AdditionalInfo:   &record[1],
		RoomNo:           &record[2],
		PcNo:             &record[3],
		AllowedContest:   &contestId,
		Username:         username,
		Password:         "",
		UnHashedPassword: password,
		Role:             "user",
	}, nil
}
