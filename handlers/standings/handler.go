package leaderboard

type Standings struct {
	ContestId int64 `json:"contest_id" db:"contest_id"`
	UserId    int64 `json:"user_id" db:"user_id"`
	Score     int   `json:"score" db:"score"`
	Penalty   int   `json:"penalty" db:"penalty"`
}
