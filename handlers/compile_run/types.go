package compilerun

import "github.com/judgenot0/judge-backend/config"

type Handler struct {
	config *config.Config
}

type Testcase struct {
	Input          string `json:"input" db:"input"`
	ExpectedOutput string `json:"expected_output" db:"expected_output"`
}

type Submission struct {
	Language    string     `json:"language"`
	SourceCode  string     `json:"source_code"`
	Testcases   []Testcase `json:"testcases"`
	Timelimit   float32    `json:"time_limit"`
	MemoryLimit float32    `json:"memory_limit"`
}
