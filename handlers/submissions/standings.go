package submissions

import (
	"database/sql"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	PenaltyPerWrongSubmission = 15 // minutes penalty for each wrong submission
)

type submissionInfo struct {
	UserID      int64         `db:"user_id"`
	ContestID   sql.NullInt64 `db:"contest_id"`
	ProblemID   int64         `db:"problem_id"`
	SubmittedAt time.Time     `db:"submitted_at"`
}

func (h *Handler) updateStandingsForAccepted(submissionID int64) {
	info, err := h.fetchSubmissionContext(submissionID)
	if err != nil {
		log.Println("standings context error:", err)
		return
	}
	if info == nil || !info.ContestID.Valid {
		return
	}

	contestID := info.ContestID.Int64

	// Check if already solved
	tx, err := h.db.Beginx()
	if err != nil {
		log.Println("standings tx begin error:", err)
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var alreadySolved bool
	err = tx.Get(&alreadySolved, `SELECT EXISTS (SELECT 1 FROM contest_solves WHERE contest_id=$1 AND user_id=$2 AND problem_id=$3)`, contestID, info.UserID, info.ProblemID)
	if err != nil {
		log.Println("standings check exists error:", err)
		return
	}
	if alreadySolved {
		err = tx.Commit()
		if err != nil {
			log.Println("standings commit error:", err)
		}
		return
	}

	// Get problem index
	var problemIndex int
	err = tx.Get(&problemIndex, `SELECT index FROM contest_problems WHERE contest_id=$1 AND problem_id=$2`, contestID, info.ProblemID)
	if err != nil {
		log.Println("standings get problem index error:", err)
		return
	}

	// Check if this is the first AC for this problem in this contest (first blood)
	var isFirstBlood bool
	err = tx.Get(&isFirstBlood, `SELECT NOT EXISTS (
		SELECT 1 FROM submissions 
		WHERE contest_id=$1 AND problem_id=$2 
		AND verdict = 'ac'
		AND id < $3
	)`, contestID, info.ProblemID, submissionID)
	if err != nil {
		log.Println("standings check first blood error:", err)
		return
	}

	// Mark the submission as first blood if applicable
	if isFirstBlood {
		_, err = tx.Exec(`UPDATE submissions SET first_blood = true WHERE id = $1`, submissionID)
		if err != nil {
			log.Println("standings mark first blood error:", err)
			return
		}
	}

	penalty, err := h.calculatePenalty(tx, contestID, info)
	if err != nil {
		log.Println("standings penalty error:", err)
		return
	}

	// Count total attempts for this problem by this user
	var attemptCount int
	err = tx.Get(&attemptCount, `SELECT COUNT(*) FROM submissions 
		WHERE contest_id=$1 AND user_id=$2 AND problem_id=$3 AND submitted_at <= $4`,
		contestID, info.UserID, info.ProblemID, info.SubmittedAt)
	if err != nil {
		log.Println("standings count attempts error:", err)
		return
	}

	// Insert into contest_solves (keep for backward compatibility)
	_, err = tx.Exec(`INSERT INTO contest_solves (contest_id, user_id, problem_id, solved_at, penalty, attempt_count, first_blood) VALUES ($1, $2, $3, $4, $5, $6, $7)`, contestID, info.UserID, info.ProblemID, info.SubmittedAt, penalty, attemptCount, isFirstBlood)
	if err != nil {
		log.Println("standings insert solve error:", err)
		return
	}

	// Update contest_user_problems (new optimized table)
	_, err = tx.Exec(`
		INSERT INTO contest_user_problems (contest_id, user_id, problem_id, problem_index, is_solved, solved_at, penalty, attempt_count, first_blood)
		VALUES ($1, $2, $3, $4, TRUE, $5, $6, $7, $8)
		ON CONFLICT (contest_id, user_id, problem_id)
		DO UPDATE SET 
			is_solved = TRUE,
			solved_at = EXCLUDED.solved_at,
			penalty = EXCLUDED.penalty,
			attempt_count = EXCLUDED.attempt_count,
			first_blood = EXCLUDED.first_blood
	`, contestID, info.UserID, info.ProblemID, problemIndex, info.SubmittedAt, penalty, attemptCount, isFirstBlood)
	if err != nil {
		log.Println("standings upsert user_problem error:", err)
		return
	}

	// Update contest_standings with enhanced fields
	_, err = tx.Exec(`
		INSERT INTO contest_standings (contest_id, user_id, penalty, solved_count, last_solved_at)
		VALUES ($1, $2, $3, 1, $4)
		ON CONFLICT (contest_id, user_id)
		DO UPDATE SET 
			penalty = contest_standings.penalty + EXCLUDED.penalty,
			solved_count = contest_standings.solved_count + 1,
			last_solved_at = GREATEST(contest_standings.last_solved_at, EXCLUDED.last_solved_at)
	`, contestID, info.UserID, penalty, info.SubmittedAt)
	if err != nil {
		log.Println("standings upsert error:", err)
		return
	}

	// Update contest_problem_stats
	_, err = tx.Exec(`
		INSERT INTO contest_problem_stats (contest_id, problem_id, problem_index, solved_count, attempted_users)
		VALUES ($1, $2, $3, 1, 1)
		ON CONFLICT (contest_id, problem_id)
		DO UPDATE SET solved_count = contest_problem_stats.solved_count + 1
	`, contestID, info.ProblemID, problemIndex)
	if err != nil {
		log.Println("standings update problem stats error:", err)
		return
	}

	if err = tx.Commit(); err != nil {
		log.Println("standings commit error:", err)
	}
}

func (h *Handler) fetchSubmissionContext(submissionID int64) (*submissionInfo, error) {
	var info submissionInfo
	err := h.db.Get(&info, `SELECT user_id, contest_id, problem_id, submitted_at FROM submissions WHERE id=$1`, submissionID)
	if err != nil {
		return nil, err
	}

	if !info.ContestID.Valid {
		return &info, nil
	}

	return &info, nil
}

func (h *Handler) calculatePenalty(tx *sqlx.Tx, contestID int64, info *submissionInfo) (int, error) {
	var wrongCount int
	err := tx.Get(&wrongCount, `
		SELECT COUNT(*) 
		FROM submissions 
		WHERE contest_id=$1 AND user_id=$2 AND problem_id=$3 AND submitted_at < $4 
		AND verdict IN ('wa','tle','re','mle')
	`, contestID, info.UserID, info.ProblemID, info.SubmittedAt)
	if err != nil {
		return 0, err
	}

	var contestStart time.Time
	if err := tx.Get(&contestStart, `SELECT start_time FROM contests WHERE id=$1`, contestID); err != nil {
		return 0, err
	}

	elapsed := info.SubmittedAt.Sub(contestStart)
	if elapsed < 0 {
		elapsed = 0
	}

	elapsedMinutes := int(elapsed.Minutes())
	if elapsedMinutes < 0 {
		elapsedMinutes = 0
	}

	return elapsedMinutes + wrongCount*PenaltyPerWrongSubmission, nil
}

func (h *Handler) updateStandingsForNonAccepted(submissionID int64, verdict string) {
	info, err := h.fetchSubmissionContext(submissionID)
	if err != nil {
		log.Println("standings context error:", err)
		return
	}
	if info == nil || !info.ContestID.Valid {
		return
	}

	contestID := info.ContestID.Int64

	// Only track certain wrong verdicts
	lowerVerdict := verdict
	if lowerVerdict != "wa" && lowerVerdict != "tle" && lowerVerdict != "re" && lowerVerdict != "mle" {
		return
	}

	tx, err := h.db.Beginx()
	if err != nil {
		log.Println("standings tx begin error:", err)
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Get problem index
	var problemIndex int
	err = tx.Get(&problemIndex, `SELECT index FROM contest_problems WHERE contest_id=$1 AND problem_id=$2`, contestID, info.ProblemID)
	if err != nil {
		log.Println("standings get problem index error:", err)
		return
	}

	// Update contest_user_problems with attempt count
	_, err = tx.Exec(`
		INSERT INTO contest_user_problems (contest_id, user_id, problem_id, problem_index, is_solved, attempt_count)
		VALUES ($1, $2, $3, $4, FALSE, 1)
		ON CONFLICT (contest_id, user_id, problem_id)
		DO UPDATE SET attempt_count = contest_user_problems.attempt_count + 1
	`, contestID, info.UserID, info.ProblemID, problemIndex)
	if err != nil {
		log.Println("standings upsert user_problem error:", err)
		return
	}

	// Update contest_standings wrong_attempts
	_, err = tx.Exec(`
		INSERT INTO contest_standings (contest_id, user_id, penalty, solved_count, wrong_attempts)
		VALUES ($1, $2, 0, 0, 1)
		ON CONFLICT (contest_id, user_id)
		DO UPDATE SET wrong_attempts = contest_standings.wrong_attempts + 1
	`, contestID, info.UserID)
	if err != nil {
		log.Println("standings update wrong_attempts error:", err)
		return
	}

	// Update contest_problem_stats - increment attempted_users if this is first attempt
	_, err = tx.Exec(`
		INSERT INTO contest_problem_stats (contest_id, problem_id, problem_index, solved_count, attempted_users)
		VALUES ($1, $2, $3, 0, 1)
		ON CONFLICT (contest_id, problem_id)
		DO UPDATE SET attempted_users = (
			SELECT COUNT(DISTINCT user_id)
			FROM contest_user_problems
			WHERE contest_id = $1 AND problem_id = $2
		)
	`, contestID, info.ProblemID, problemIndex)
	if err != nil {
		log.Println("standings update problem stats error:", err)
		return
	}

	if err = tx.Commit(); err != nil {
		log.Println("standings commit error:", err)
	}
}
