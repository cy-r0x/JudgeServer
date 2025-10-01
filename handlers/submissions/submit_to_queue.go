package submissions

import "log"

func (h *Handler) submitToQueue(submissionId int64, submission *UserSubmission) error {

	// Fetch problem limits
	var limits struct {
		TimeLimit   float32 `db:"time_limit"`
		MemoryLimit float32 `db:"memory_limit"`
	}
	if err := h.db.Get(&limits, `SELECT time_limit, memory_limit FROM problems WHERE id=$1`, submission.ProblemId); err != nil {
		log.Println("Error fetching problem limits for queue:", err)
		return err
	}

	// Fetch all testcases for the problem
	var testcases []Testcase
	if err := h.db.Select(&testcases, `
		SELECT input, expected_output
		FROM testcases
		WHERE problem_id = $1
		ORDER BY is_sample DESC, id ASC
	`, submission.ProblemId); err != nil {
		log.Println("Error fetching testcases for queue:", err)
		return err
	}

	queueData := QueueSubmission{
		SubmissionId: submissionId,
		Language:     submission.Language,
		SourceCode:   submission.SourceCode,
		Testcases:    testcases,
		Timelimit:    limits.TimeLimit,
		MemoryLimit:  limits.MemoryLimit,
	}
	log.Println(queueData)
	// TODO: Send queueData to the judging queue/engine
	return nil
}
