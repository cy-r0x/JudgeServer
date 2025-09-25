package users

import (
	"github.com/jmoiron/sqlx"
	"github.com/judgenot0/judge-backend/config"
)

type Handler struct {
	config *config.Config
	db     *sqlx.DB
}

func NewHandler(config *config.Config, db *sqlx.DB) *Handler {
	return &Handler{
		config: config,
		db:     db,
	}
}
