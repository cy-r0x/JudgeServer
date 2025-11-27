package users

import (
	"github.com/jmoiron/sqlx"
	"github.com/judgenot0/judge-backend/config"
)

func NewHandler(config *config.Config, db *sqlx.DB) *Handler {
	return &Handler{
		config: config,
		db:     db,
	}
}
