package submissions

import (
	"log"
	"time"

	"gorm.io/gorm"
)

const (
	PenaltyPerWrongSubmission = 15 // minutes penalty for each wrong submission
)

type submissionInfo struct {
	UserID      int64     `gorm:"column:user_id"`
	ContestID   *int64    `gorm:"column:contest_id"`
	ProblemID   int64     `gorm:"column:problem_id"`
	SubmittedAt time.Time `gorm:"column:submitted_at"`
}

func (h *Handler) updateStandingsForAccepted(submissionID int64) {
	info, err := h.fetchSubmissionContext(submissionID)
	if err != nil {
		log.Println("standings context error:", err)
		return
	}
	if info == nil || info.ContestID == nil {
		return
	}

	contestID := *info.ContestID

	tx := h.db.Begin()
	if tx.Error != nil {
		log.Println("standings tx begin error:", tx.Error)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if already solved
	var alreadySolved bool
	err = tx.Raw(`SELECT EXISTS (SELECT 1 FROM contest_solves WHERE contest_id=? AND user_id=? AND problem_id=?)`, contestID, info.UserID, info.ProblemID).Scan(&alreadySolved).Error
	if err != nil {
		tx.Rollback()
		log.Println("standings check exists error:", err)
		return
	}
	if alreadySolved {
		if err := tx.Commit().Error; err != nil {
			log.Println("standings commit error:", err)
		}
		return
	}

	// Get problem index
	var problemIndex int
	err = tx.Raw(`SELECT index FROM contest_problems WHERE contest_id=? AND problem_id=?`, contestID, info.ProblemID).Scan(&problemIndex).Error
	if err != nil {
		tx.Rollback()
		log.Println("standings get problem index error:", err)
		return
	}

	// Check if this is the first AC for this problem in this contest (first blood)
	var isFirstBlood bool
	err = tx.Raw(`SELECT NOT EXISTS (
		SELECT 1 FROM submissions 
		WHERE contest_id=? AND problem_id=? 
		AND verdict = 'ac'
		AND id < ?
	)`, contestID, info.ProblemID, submissionID).Scan(&isFirstBlood).Error
	if err != nil {
		tx.Rollback()
		log.Println("standings check first blood error:", err)
		return
	}

	// Mark the submission as first blood if applicable
	if isFirstBlood {
		if err := tx.Exec(`UPDATE submissions SET first_blood = true WHERE id = ?`, submissionID).Error; err != nil {
			tx.Rollback()
			log.Println("standings mark first blood error:", err)
			return
		}
	}

	penalty, err := h.calculatePenalty(tx, contestID, info)
	if err != nil {
		tx.Rollback()
		log.Println("standings penalty error:", err)
		return
	}

	// Count total attempts for this problem by this user
	var attemptCount int
	err = tx.Raw(`SELECT COUNT(*) FROM submissions 
		WHERE contest_id=? AND user_id=? AND problem_id=? AND submitted_at <= ?`,
		contestID, info.UserID, info.ProblemID, info.SubmittedAt).Scan(&attemptCount).Error
	if err != nil {
		tx.Rollback()
		log.Println("standings count attempts error:", err)
		return
	}

	// Insert into contest_solves (keep for backward compatibility)
	err = tx.Exec(`INSERT INTO contest_solves (contest_id, user_id, problem_id, solved_at, penalty, attempt_count, first_blood) VALUES (?, ?, ?, ?, ?, ?, ?)`, contestID, info.UserID, info.ProblemID, info.SubmittedAt, penalty, attemptCount, isFirstBlood).Error
	if err != nil {
		tx.Rollback()
		log.Println("standings insert solve error:", err)
		return
	}

	// Update contest_user_problems (new optimized table)
	err = tx.Exec(`
		INSERT INTO contest_user_problems (contest_id, user_id, problem_id, problem_index, is_solved, solved_at, penalty, attempt_count, first_blood)
		VALUES (?, ?, ?, ?, TRUE, ?, ?, ?, ?)
		ON CONFLICT (contest_id, user_id, problem_id)
		DO UPDATE SET 
			is_solved = TRUE,
			solved_at = EXCLUDED.solved_at,
			penalty = EXCLUDED.penalty,
			attempt_count = EXCLUDED.attempt_count,
			first_blood = EXCLUDED.first_blood
	`, contestID, info.UserID, info.ProblemID, problemIndex, info.SubmittedAt, penalty, attemptCount, isFirstBlood).Error
	if err != nil {
		tx.Rollback()
		log.Println("standings upsert user_problem error:", err)
		return
	}

	// Update contest_standings with enhanced fields
	err = tx.Exec(`
		INSERT INTO contest_standings (contest_id, user_id, penalty, solved_count, last_solved_at)
		VALUES (?, ?, ?, 1, ?)
		ON CONFLICT (contest_id, user_id)
		DO UPDATE SET 
			penalty = contest_standings.penalty + EXCLUDED.penalty,
			solved_count = contest_standings.solved_count + 1,
			last_solved_at = GREATEST(contest_standings.last_solved_at, EXCLUDED.last_solved_at)
	`, contestID, info.UserID, penalty, info.SubmittedAt).Error
	if err != nil {
		tx.Rollback()
		log.Println("standings upsert error:", err)
		return
	}

	// Update contest_problem_stats
	err = tx.Exec(`
		INSERT INTO contest_problem_stats (contest_id, problem_id, problem_index, solved_count, attempted_users)
		VALUES (?, ?, ?, 1, 1)
		ON CONFLICT (contest_id, problem_id)
		DO UPDATE SET solved_count = contest_problem_stats.solved_count + 1
	`, contestID, info.ProblemID, problemIndex).Error
	if err != nil {
		tx.Rollback()
		log.Println("standings update problem stats error:", err)
		return
	}

	if err = tx.Commit().Error; err != nil {
		log.Println("standings commit error:", err)
	}
}

func (h *Handler) fetchSubmissionContext(submissionID int64) (*submissionInfo, error) {
	var info submissionInfo
	err := h.db.Raw(`SELECT user_id, contest_id, problem_id, submitted_at FROM submissions WHERE id=?`, submissionID).Scan(&info).Error
	if err != nil {
		return nil, err
	}

	if info.ContestID == nil {
		return &info, nil
	}

	return &info, nil
}

func (h *Handler) calculatePenalty(tx *gorm.DB, contestID int64, info *submissionInfo) (int, error) {
	var wrongCount int
	err := tx.Raw(`
		SELECT COUNT(*) 
		FROM submissions 
		WHERE contest_id=? AND user_id=? AND problem_id=? AND submitted_at < ? 
		AND verdict IN ('wa','tle','re','mle')
	`, contestID, info.UserID, info.ProblemID, info.SubmittedAt).Scan(&wrongCount).Error
	if err != nil {
		return 0, err
	}

	var contestStart time.Time
	if err := tx.Raw(`SELECT start_time FROM contests WHERE id=?`, contestID).Scan(&contestStart).Error; err != nil {
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
	if info == nil || info.ContestID == nil {
		return
	}

	contestID := *info.ContestID

	// Only track certain wrong verdicts
	lowerVerdict := verdict
	if lowerVerdict != "wa" && lowerVerdict != "tle" && lowerVerdict != "re" && lowerVerdict != "mle" {
		return
	}

	tx := h.db.Begin()
	if tx.Error != nil {
		log.Println("standings tx begin error:", tx.Error)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get problem index
	var problemIndex int
	err = tx.Raw(`SELECT index FROM contest_problems WHERE contest_id=? AND problem_id=?`, contestID, info.ProblemID).Scan(&problemIndex).Error
	if err != nil {
		tx.Rollback()
		log.Println("standings get problem index error:", err)
		return
	}

	// Update contest_user_problems with attempt count
	err = tx.Exec(`
		INSERT INTO contest_user_problems (contest_id, user_id, problem_id, problem_index, is_solved, attempt_count)
		VALUES (?, ?, ?, ?, FALSE, 1)
		ON CONFLICT (contest_id, user_id, problem_id)
		DO UPDATE SET attempt_count = contest_user_problems.attempt_count + 1
	`, contestID, info.UserID, info.ProblemID, problemIndex).Error
	if err != nil {
		tx.Rollback()
		log.Println("standings upsert user_problem error:", err)
		return
	}

	// Update contest_standings wrong_attempts
	err = tx.Exec(`
		INSERT INTO contest_standings (contest_id, user_id, penalty, solved_count, wrong_attempts)
		VALUES (?, ?, 0, 0, 1)
		ON CONFLICT (contest_id, user_id)
		DO UPDATE SET wrong_attempts = contest_standings.wrong_attempts + 1
	`, contestID, info.UserID).Error
	if err != nil {
		tx.Rollback()
		log.Println("standings update wrong_attempts error:", err)
		return
	}

	// Update contest_problem_stats - increment attempted_users if this is first attempt
	err = tx.Exec(`
		INSERT INTO contest_problem_stats (contest_id, problem_id, problem_index, solved_count, attempted_users)
		VALUES (?, ?, ?, 0, 1)
		ON CONFLICT (contest_id, problem_id)
		DO UPDATE SET attempted_users = (
			SELECT COUNT(DISTINCT user_id)
			FROM contest_user_problems
			WHERE contest_id = ? AND problem_id = ?
		)
	`, contestID, info.ProblemID, problemIndex, contestID, info.ProblemID).Error
	if err != nil {
		tx.Rollback()
		log.Println("standings update problem stats error:", err)
		return
	}

	if err = tx.Commit().Error; err != nil {
		log.Println("standings commit error:", err)
	}
}
