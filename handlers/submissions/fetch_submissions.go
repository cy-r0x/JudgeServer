package submissions

import "time"

func (h *Handler) fetchSubmissions(params SubmissionListParams) ([]SubmissionResponse, int, error) {
	// Normalize pagination
	if params.Limit <= 0 || params.Limit > 100 {
		params.Limit = 20
	}
	if params.Page < 1 {
		params.Page = 1
	}
	offset := (params.Page - 1) * params.Limit

	type SubmissionRow struct {
		ID         int64     `json:"id"`
		UserID     string    `json:"userId"`
		Name       string    `json:"name"`
		Username   string    `json:"username"`
		ProblemID  string    `json:"problemId"`
		ProblemIdx string    `json:"problemIndex"`
		ContestID  *string   `json:"contestId"`
		Language   string    `json:"language"`
		Status     string    `json:"status"`
		ExecTime   *float64  `json:"execTime"`
		ExecMemory *float64  `json:"execMemory"`
		CreatedAt  time.Time `json:"createdAt"`
		TotalCount int       `json:"-"`
	}

	var rows []SubmissionRow
	query := `
		SELECT 
			s.id,
			s.user_id,
			u.name,
			u.username,
			s.problem_id,
			cp.index AS problem_index,
			s.contest_id,
			s.language,
			s.status,
			s.exec_time,
			s.exec_memory,
			s.created_at,
			COUNT(*) OVER() AS total_count
		FROM submissions s
		JOIN users u ON s.user_id = u.id
		LEFT JOIN contest_problems cp 
			ON s.contest_id = cp.contest_id AND s.problem_id = cp.problem_id
		WHERE s.contest_id = ?
	`
	args := []any{params.ContestID}

	// 🔐 Optional: filter by user (only if UserID is provided)
	if params.UserID != nil {
		query += " AND s.user_id = ?"
		args = append(args, *params.UserID)
	}

	// 🔍 Apply filters *if provided* (no scope restrictions)
	if params.Status != "" {
		query += " AND s.status = ?"
		args = append(args, params.Status)
	}
	if params.SearchName != "" {
		query += " AND u.name ILIKE ?"
		args = append(args, "%"+params.SearchName+"%")
	}
	if params.SearchUsername != "" {
		query += " AND u.username ILIKE ?"
		args = append(args, "%"+params.SearchUsername+"%")
	}

	query += " ORDER BY s.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, params.Limit, offset)

	if err := h.db.Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	// Map to response DTOs (🔒 SourceCode excluded in list view)
	submissions := make([]SubmissionResponse, 0, len(rows))
	totalCount := 0
	for _, row := range rows {
		totalCount = row.TotalCount
		submissions = append(submissions, SubmissionResponse{
			Id:         row.ID,
			UserId:     row.UserID,
			Name:       row.Name,
			Username:   row.Username,
			ProblemId:  row.ProblemID,
			ContestId:  row.ContestID,
			Language:   row.Language,
			Status:     row.Status,
			ExecTime:   row.ExecTime,
			ExecMemory: row.ExecMemory,
			CreatedAt:  row.CreatedAt,
			// SourceCode intentionally omitted
		})
	}

	return submissions, totalCount, nil
}
