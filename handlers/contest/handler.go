package contest

import (
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/judgenot0/judge-backend/utils"
)

type Contest struct {
	Id              int64     `json:"id" db:"id"`
	Title           string    `json:"title" db:"title"`
	Description     string    `json:"description" db:"description"`
	StartTime       time.Time `json:"start_time" db:"start_time"`
	DurationSeconds int64     `json:"duration_seconds" db:"duration_seconds"`
	Status          string    `json:"status" db:"status"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

type ContestProblem struct {
	ContestId int64  `json:"contest_id" db:"contest_id"`
	ProblemId int64  `json:"problem_id" db:"problem_id"`
	Index     string `json:"index" db:"index"`
}

type Handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) CreateContest(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) UpdateContest(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) ListContests(w http.ResponseWriter, r *http.Request) {
	contests := []Contest{}
	//TODO: Add Dynamic DB fetch of contests
	utils.SendResponse(w, http.StatusOK, contests)
}

func (h *Handler) GetContest(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) DeleteContest(w http.ResponseWriter, r *http.Request) {

}
