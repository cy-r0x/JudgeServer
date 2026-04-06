package compilerun

import (
	"github.com/judgenot0/judge-backend/config"
	"gorm.io/gorm"
)

func NewHandler(db *gorm.DB, config *config.Config) *Handler {
	return &Handler{
		db:     db,
		config: config,
	}
}
