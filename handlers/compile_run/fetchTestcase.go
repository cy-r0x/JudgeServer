package compilerun

import (
	"log"

	"github.com/judgenot0/judge-backend/models"
)

func (h *Handler) fetchTestcases(problemId string, isSample bool) ([]Testcase, error) {

	query := h.db.Model(&models.Testcase{}).Select("input", "expected_output").Where("problem_id = ?", problemId)

	if isSample {
		query = query.Where("is_sample = ?", true)
	}

	query = query.Order("is_sample DESC, id ASC")

	var testcases []Testcase
	if err := query.Scan(&testcases).Error; err != nil {
		log.Println(err)
		return nil, err
	}

	if testcases == nil {
		testcases = []Testcase{}
	}

	return testcases, nil
}
