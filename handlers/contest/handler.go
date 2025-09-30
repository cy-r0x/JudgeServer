package contest

import (
	"encoding/json"
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
	decoder := json.NewDecoder(r.Body)
	var contest Contest
	decoder.Decode(&contest)

	//TODO: Update Contest into DB
}

func (h *Handler) UpdateContestIndex(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var contestProblem ContestProblem
	decoder.Decode(&contestProblem)

	//TODO: Update Problem Index
}

func (h *Handler) ListContests(w http.ResponseWriter, r *http.Request) {
	contests := []Contest{}
	//TODO: Add Dynamic DB fetch of contests
	utils.SendResponse(w, http.StatusOK, contests)
}

func (h *Handler) GetContest(w http.ResponseWriter, r *http.Request) {
	contestId := r.PathValue("contestId")
	if contestId == "" {
		utils.SendResponse(w, http.StatusNotFound, "Contest Not Found")
		return
	}
	//TODO: Get Contest Information -> Get Contest Problems -> Send Response
}
