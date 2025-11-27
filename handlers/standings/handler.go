package standings

import (
	"time"

	"github.com/jmoiron/sqlx"
)

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{
		db: db,
		Last_standings: make(map[int64]struct {
			timestamp *time.Time
			standings *StandingsResponse
		}),
	}
}
