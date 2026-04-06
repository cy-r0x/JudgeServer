package standings

import (
	"time"

	"gorm.io/gorm"
)

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		db: db,
		Last_standings: make(map[int64]struct {
			timestamp *time.Time
			standings *StandingsResponse
		}),
	}
}
