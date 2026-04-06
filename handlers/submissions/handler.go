package submissions

import (
	"github.com/judgenot0/judge-backend/config"
	"gorm.io/gorm"
)

type Handler struct {
	db     *gorm.DB
	config *config.Config
}

func NewHandler(db *gorm.DB, config *config.Config) *Handler {
	return &Handler{
		db:     db,
		config: config,
	}
}
