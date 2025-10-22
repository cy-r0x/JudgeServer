package usercsv

import (
	"fmt"
	"strings"

	"github.com/judgenot0/judge-backend/handlers/structs"
)

func (h *Handler) WriteUser(record []string, idx int) (structs.User, error) {
	if len(record) < 4 {
		return structs.User{}, fmt.Errorf("record length mismatch") // skip invalid rows
	}

	row := []string{record[0]} // Name
	clan_data := strings.Split(record[1], ",")
	if len(clan_data) != h.clanLength {
		return structs.User{}, fmt.Errorf("clan data length mismatch")
	}

	for i := 0; i < h.clanLength; i++ {
		clan_data[i] = strings.TrimSpace(clan_data[i])
	}

	row = append(row, clan_data...)
	row = append(row, record[2]) // Room No
	row = append(row, record[3]) // PC No
	username := fmt.Sprintf("%s_%d", h.prefix, idx+1)
	password := generatePassword() // Generate password once
	row = append(row, username)    // Username
	row = append(row, password)    // Password

	// Lock the CSV writer for thread-safe writing
	h.mu.Lock()
	err := h.writer.Write(row)
	h.writer.Flush()
	h.mu.Unlock()

	if err != nil {
		fmt.Println("Error writing row to CSV:", err)
		return structs.User{}, err
	}

	return structs.User{
		FullName:       record[0],
		Clan:           &record[1],
		RoomNo:         &record[2],
		PcNo:           &record[3],
		AllowedContest: h.contestId,
		Username:       username,
		Password:       password, // Use the same password
		Role:           "user",
	}, nil
}
