package utils

type Contest struct {
	ContestId   int    `json:"contestId"`
	ContestName string `json:"contestName"`
	StartTime   uint64 `json:"startTime"`
	EndTime     uint64 `json:"endTime"`
	Duration    uint64 `json:"duration"`
	Status      string `json:"status"`
}

func GetContests() ([]Contest, error) {
	contests := []Contest{}
	//TODO: Add Dynamic DB fetch of contests
	contests = append(contests, Contest{
		ContestId:   10,
		ContestName: "UTA-2025",
		StartTime:   1234,
		EndTime:     4567,
		Duration:    1234,
		Status:      "Finished",
	})
	return contests, nil
}
