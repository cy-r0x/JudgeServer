package problem

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/judgenot0/judge-backend/middlewares"
	"github.com/judgenot0/judge-backend/utils"
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
	Id              int64      `json:"id" db:"id"`
	Title           string     `json:"title" db:"title"`
	Slug            string     `json:"slug" db:"slug"`
	Statement       string     `json:"statement" db:"statement"`
	InputStatement  string     `json:"input_statement" db:"input_statement"`
	OutputStatement string     `json:"output_statement" db:"output_statement"`
	TimeLimitMs     int        `json:"time_limit_ms" db:"time_limit_ms"`
	MemoryLimitMb   int        `json:"memory_limit_mb" db:"memory_limit_mb"`
	Testcases       []Testcase `json:"test_cases"`
	CreatedBy       int64      `json:"created_by" db:"created_by"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
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
	decoder := json.NewDecoder(r.Body)

	payload, ok := r.Context().Value("user").(*middlewares.Payload)
	if !ok {
		utils.SendResponse(w, http.StatusUnauthorized, "User information not found")
		return
	}

	var problem Problem
	decoder.Decode(&problem)

	problem.CreatedBy = payload.Sub
	problem.Slug = strings.ReplaceAll(strings.ToLower(problem.Title), " ", "-")
	problem.CreatedAt = time.Now()

	//Push the Problem To DB and get Problem Id

	testcases := problem.Testcases
	log.Println(testcases)

	//Loop over all the testcases and add to Testcases DB with Problem ID

	//Send Response
}

func (h *Handler) ListProblems(w http.ResponseWriter, r *http.Request) {
	payload, ok := r.Context().Value("user").(*middlewares.Payload)
	if !ok {
		utils.SendResponse(w, http.StatusUnauthorized, "User information not found")
		return
	}
	switch payload.Role {
	case "user":
		//Check if the user have access to the contest -> return problem lists
	case "admin":
		//Return the contest problem list
	}
}

func (h *Handler) GetProblem(w http.ResponseWriter, r *http.Request) {
	payload, ok := r.Context().Value("user").(*middlewares.Payload)
	if !ok {
		utils.SendResponse(w, http.StatusUnauthorized, "User information not found")
		return
	}
	switch payload.Role {
	case "user":
		//check if the user have access to the contest -> return the problem
	case "setter":
		//check if the setter owns the problem or not -> return the problem
	case "admin":
		//return the problem data
	}
}

func (h *Handler) UpdateProblem(w http.ResponseWriter, r *http.Request) {
	//Only Setter can update the problem
}
