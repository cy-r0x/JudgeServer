package handlers

import (
	"github.com/judgenot0/judge-backend/structs"
)

func GetContest() ([]structs.Contest, error) {
	contests := []structs.Contest{}
	contests = append(contests, structs.Contest{
		ContestId:   10,
		ContestName: "UTA-2025",
		StartTime:   1234,
		EndTime:     4567,
		Duration:    1234,
		Status:      "Finished",
	})
	return contests, nil
}
