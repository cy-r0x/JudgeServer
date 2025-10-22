package usercsv

import "fmt"

func (h *Handler) WriteHeader() error {
	defer h.writer.Flush()
	row := []string{"Name"}
	for i := 1; i <= h.clanLength; i++ {
		row = append(row, fmt.Sprintf("Clan_%d", i))
	}
	row = append(row, "Room No", "PC No", "Username", "Password")
	if err := h.writer.Write(row); err != nil {
		fmt.Println("Error writing row to CSV:", err)
		return err
	}
	return nil
}
