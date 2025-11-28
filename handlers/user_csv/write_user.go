package usercsv

import (
	"fmt"
	"strings"
)

func (h *Handler) WriteUser(user User) error {

	defer h.Writer.writer.Flush()
	clan := ""
	if user.Clan != nil {
		clan = *user.Clan
	}

	roomNo := ""
	if user.RoomNo != nil {
		roomNo = *user.RoomNo
	}
	pcNo := ""
	if user.PcNo != nil {
		pcNo = *user.PcNo
	}

	row := []string{user.FullName}
	clan_data := strings.Split(clan, ",")
	if len(clan_data) != h.Writer.clanLength {
		return fmt.Errorf("clan data length mismatch")
	}

	for i := 0; i < h.Writer.clanLength; i++ {
		clan_data[i] = strings.TrimSpace(clan_data[i])
	}

	row = append(row, clan_data...)
	row = append(row, roomNo)                // Room No
	row = append(row, pcNo)                  // PC No
	row = append(row, user.Username)         // Username
	row = append(row, user.UnHashedPassword) // Password

	h.mu.Lock()
	err := h.Writer.writer.Write(row)
	h.Writer.writer.Flush()
	h.mu.Unlock()

	if err != nil {
		fmt.Println("Error writing row to CSV:", err)
		return fmt.Errorf("error writing row to CSV: %w", err)
	}
	return nil

}
