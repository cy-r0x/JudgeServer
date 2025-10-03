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

	_, err = tx.Exec(`INSERT INTO contest_solves (contest_id, user_id, problem_id, solved_at, penalty, attempt_count, first_blood) VALUES ($1, $2, $3, $4, $5, $6, $7)`, contestID, info.UserID, info.ProblemID, info.SubmittedAt, penalty, attemptCount, isFirstBlood)
	if err != nil {
		log.Println("standings insert solve error:", err)
		return
	}

	_, err = tx.Exec(`INSERT INTO contest_standings (contest_id, user_id, penalty, solved_count) VALUES ($1, $2, $3, 1)
        ON CONFLICT (contest_id, user_id)
        DO UPDATE SET penalty = contest_standings.penalty + EXCLUDED.penalty,
                      solved_count = contest_standings.solved_count + 1`, contestID, info.UserID, penalty)
	if err != nil {
		log.Println("standings upsert error:", err)
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
