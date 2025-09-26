package problem

import (
	"net/http"

	"github.com/jmoiron/sqlx"
)

type Testcase struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type Problem struct {
	ProblemId         int        `json:"problemId"`
	ProblemName       string     `json:"problemName"`
	Author            string     `json:"author"`
	BodyDescription   string     `json:"bodyDescription"`
	InputDescription  string     `json:"inputDescription"`
	OutputDescription string     `json:"outputDescription"`
	SampleTestcases   []Testcase `json:"sampleTestcase"`
	RegularTestcases  []Testcase `json:"regularTestcase"`
	TimeLimit         float64    `json:"timeLimit"`
	MemoryLimit       float64    `json:"memoryLimit"`
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
