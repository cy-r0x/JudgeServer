package problem

import (
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

type Testcase struct {
	Id             int64     `json:"id" db:"id"`
	ProblemId      int64     `json:"problem_id" db:"problem_id"`
	Input          string    `json:"input" db:"input"`
	ExpectedOutput string    `json:"expected_output" db:"expected_output"`
	IsSample       bool      `json:"is_sample" db:"is_sample"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type Problem struct {
	Id            int64     `json:"id" db:"id"`
	Title         string    `json:"title" db:"title"`
	Slug          string    `json:"slug" db:"slug"`
	Statement     string    `json:"statement" db:"statement"`
	TimeLimitMs   int       `json:"time_limit_ms" db:"time_limit_ms"`
	MemoryLimitMb int       `json:"memory_limit_mb" db:"memory_limit_mb"`
	CreatedBy     int64     `json:"created_by" db:"created_by"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type ProblemList struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type Handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) CreateProblem(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) ListProblems(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) GetProblem(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) UpdateProblem(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) DeleteProblem(w http.ResponseWriter, r *http.Request) {
}
