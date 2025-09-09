package structs

type ProblemData struct {
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
