package contest_problems

import (
	"github.com/jmoiron/sqlx"
)

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{
		db: db,
	}
}
