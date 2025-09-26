package setter

import (
	"net/http"

	"github.com/jmoiron/sqlx"
)

type Handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) ListSetterProblems(w http.ResponseWriter, r *http.Request) {

}
