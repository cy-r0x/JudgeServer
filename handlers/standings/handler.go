package standings

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type Standings struct {
	ContestId  int64 `json:"contest_id" db:"contest_id"`
	UserId     int64 `json:"user_id" db:"user_id"`
	Penalty    int   `json:"penalty" db:"penalty"`
	SolveCount int8  `json:"solve_count" db:"solve_count"`
}

type Handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) GetStandings(w http.ResponseWriter, r *http.Request) {
	contestId := r.PathValue("contestId")
	log.Println(contestId)
	// TODO: Get Standings data from DB, Sort according to submission count(desc) if count is same the sort as penality (asc)
}
