package submissions

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

func (h *Handler) CreateSubmission(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) ListUserSubmissions(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) ListAllSubmissions(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) GetSubmission(w http.ResponseWriter, r *http.Request) {

}
