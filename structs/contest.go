package structs

type Contest struct {
	ContestId   int    `json:"contestId"`
	ContestName string `json:"contestName"`
	StartTime   uint64 `json:"startTime"`
	EndTime     uint64 `json:"endTime"`
	Duration    uint64 `json:"duration"`
	Status      string `json:"status"`
}