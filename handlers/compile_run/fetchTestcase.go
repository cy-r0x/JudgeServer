package compilerun

import "log"

func (h *Handler) fetchTestcases(problemId int64, isSample bool) ([]Testcase, error) {
	query := `
		SELECT input, expected_output
		FROM testcases
		WHERE problem_id = $1`

	if isSample {
		query += ` AND is_sample = TRUE`
	}

	query += ` ORDER BY is_sample DESC, id ASC`

	var testcases []Testcase
	if err := h.db.Select(&testcases, query, problemId); err != nil {
		log.Println(err)
		return nil, err
	}

	return testcases, nil
}
