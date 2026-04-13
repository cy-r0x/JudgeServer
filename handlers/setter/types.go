package setter

import (
	"time"

	"gorm.io/gorm"
)

type Problem struct {
	Id        string    `json:"id" gorm:"column:id"`
	Title     string    `json:"title" gorm:"column:title"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
}

type Handler struct {
	db *gorm.DB
}
