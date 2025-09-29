package submissions

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/judgenot0/judge-backend/middlewares"
	"github.com/judgenot0/judge-backend/utils"
)

type Submission struct {
	Id            int64     `json:"id" db:"id"`
	UserId        int64     `json:"user_id" db:"user_id"`
	ProblemId     int64     `json:"problem_id" db:"problem_id"`
	ContestId     int64     `json:"contest_id" db:"contest_id"`
	Language      string    `json:"language" db:"language"`
	SourceCode    string    `json:"source_code" db:"source_code"`
	Verdict       string    `json:"verdict" db:"verdict"`
	ExecutionTime int       `json:"execution_time" db:"execution_time"`
	MemoryUsed    int       `json:"memory_used" db:"memory_used"`
	SubmittedAt   time.Time `json:"submitted_at" db:"submitted_at"`
}

type UserSubmission struct {
	ProblemId  int64  `json:"problem_id"`
	ContestId  int64  `json:"contest_id"`
	Language   string `json:"language"`
	SourceCode string `json:"source_code"`
}

type Handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) CreateSubmission(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	payload, ok := r.Context().Value("user").(*middlewares.Payload)
	if !ok {
		utils.SendResopnse(w, http.StatusUnauthorized, "Invalid Token")
		return
	}
	userId := payload.Sub
	log.Println(userId)
	var submission UserSubmission
	decoder.Decode(&submission)
	//TODO: Add to DB -> get submission ID -> Submit to Queue
}

func (h *Handler) ListUserSubmissions(w http.ResponseWriter, r *http.Request) {
	payload, ok := r.Context().Value("user").(*middlewares.Payload)
	if !ok {
		utils.SendResopnse(w, http.StatusUnauthorized, "Invalid Token")
		return
	}
	userId := payload.Sub
	log.Println(userId)
	//TODO: Get all user submisison -> Send Response
}

func (h *Handler) ListAllSubmissions(w http.ResponseWriter, r *http.Request) {
	//TODO: Get All the submissions of current contest -> Send Response
}

func (h *Handler) GetSubmission(w http.ResponseWriter, r *http.Request) {
	payload, ok := r.Context().Value("user").(*middlewares.Payload)
	if !ok {
		utils.SendResopnse(w, http.StatusUnauthorized, "Invalid Token")
		return
	}
	userId := payload.Sub
	log.Println(userId)
	//TODO: Add admin support to get this data, get data from DB -> Send Response
}
