package contest

import (
	"github.com/jmoiron/sqlx"
	"github.com/judgenot0/judge-backend/config"
)

func NewHandler(db *sqlx.DB, config *config.Config) *Handler {
	return &Handler{
		db:     db,
		config: config,
	}
}
