package problem

import "log"

func (h *Handler) fetchTestcases(problemId string, isSample bool) ([]Testcase, error) {
	query := `
		SELECT id, problem_id, input, expected_output, is_sample, created_at
		FROM testcases
		WHERE problem_id = $1`
	args := []any{problemId}

	if isSample {
		query += ` AND is_sample = TRUE`
	}

	query += ` ORDER BY is_sample DESC, id ASC`

	var testcases []Testcase
	if err := h.db.Select(&testcases, query, args...); err != nil {
		log.Println(err)
		return nil, err
	}

	return testcases, nil
}
