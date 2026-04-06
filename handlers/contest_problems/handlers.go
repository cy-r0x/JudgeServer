package contest_problems

import (
	"gorm.io/gorm"
)

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		db: db,
	}
}
