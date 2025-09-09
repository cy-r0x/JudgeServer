package handlers

import "github.com/judgenot0/judge-backend/structs"

func GetProblems(contestId int) ([]structs.ProblemList, error) {
	problemList := []structs.ProblemList{}
	problemList = append(problemList, structs.ProblemList{
		Id:    1234,
		Title: "Adding Two Numbers",
	})
	return problemList, nil
}
