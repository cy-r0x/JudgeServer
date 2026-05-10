package standings

import (
	"sync"
	"time"

	"gorm.io/gorm"
)

type ProblemStatus struct {
	Solved        bool       `json:"solved"`
	FirstSolvedAt *time.Time `json:"firstSolvedAt,omitempty"`
	Attempts      int        `json:"attempts"`
	Penalty       int        `json:"penalty"`
	FirstBlood    bool       `json:"firstBlood"`
}

type UserStanding struct {
	UserId       string          `json:"userId"`
	Username     string          `json:"username"`
	Name         string          `json:"name"`
	TotalPenalty int             `json:"totalPenalty"`
	SolvedCount  int             `json:"solvedCount"`
	Problems     []ProblemStatus `json:"problems"`
	LastSolvedAt *time.Time      `json:"lastSolvedAt,omitempty"`
}

type ProblemSolveStatus struct {
	Solved    int `json:"solved"`
	Attempted int `json:"attempted"`
}

type StandingsResponse struct {
	ContestId          string                     `json:"contestId"`
	ContestTitle       string                     `json:"contestTitle"`
	ProblemMapping     map[int]string             `json:"problemMapping"`
	Standings          []UserStanding             `json:"standings"`
	StartTime          time.Time                  `json:"startTime"`
	DurationSeconds    int64                      `json:"durationSeconds"`
	ProblemSolveStatus map[int]ProblemSolveStatus `json:"problemSolveStatus"`
	TotalItem          int                        `json:"totalItem"`
	TotalPages         int                        `json:"totalPages"`
	Limit              int                        `json:"limit"`
	Page               int                        `json:"page"`
}

type Handler struct {
	db             *gorm.DB
	mu             sync.RWMutex
	Last_standings map[string]struct {
		timestamp *time.Time
		standings *StandingsResponse
	}
}