package users

import (
	"github.com/judgenot0/judge-backend/config"
	"gorm.io/gorm"
)

func NewHandler(config *config.Config, db *gorm.DB) *Handler {
	return &Handler{
		config: config,
		db:     db,
	}
}
