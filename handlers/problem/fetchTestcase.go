package problem

import (
	"log"

	"github.com/judgenot0/judge-backend/models"
)

func (h *Handler) fetchTestcases(problemId string, isSample bool) ([]Testcase, error) {
	var dbTestcases []models.Testcase

	query := h.db.Where("problem_id = ?", problemId)
	if isSample {
		query = query.Where("is_sample = ?", true)
	}

	err := query.Order("is_sample DESC, id ASC").Find(&dbTestcases).Error
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var testcases []Testcase
	for _, tc := range dbTestcases {
		testcases = append(testcases, Testcase{
			Id:             int64(tc.ID),
			ProblemId:      int64(tc.ProblemID),
			Input:          tc.Input,
			ExpectedOutput: tc.ExpectedOutput,
			IsSample:       tc.IsSample,
			CreatedAt:      tc.CreatedAt,
		})
	}

	if testcases == nil {
		testcases = []Testcase{}
	}

	return testcases, nil
}
