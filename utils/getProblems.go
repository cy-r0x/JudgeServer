package utils

type ProblemList struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

func GetProblems(contestId int) ([]ProblemList, error) {
	problemList := []ProblemList{}
	//TODO: Add Dynamic DB fetch of contest problem set
	problemList = append(problemList, ProblemList{
		Id:    1234,
		Title: "Adding Two Numbers",
	})
	return problemList, nil
}
