package setter

import (
	"github.com/jmoiron/sqlx"
)

type Problem struct {
	Id        int64  `json:"id" db:"id"`
	Title     string `json:"title" db:"title"`
	CreatedAt string `json:"created_at" db:"created_at"`
}

type Handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{
		db: db,
	}
}
