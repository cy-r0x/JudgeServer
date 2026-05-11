package submissions

import (
	"log/slog"

	"github.com/judgenot0/judge-backend/models"
	"gorm.io/gorm"
)

const (
	PenaltyPerWrongSubmission = 20 // minutes penalty for each wrong submission
)

func (h *Handler) updateStandingsForAccepted(submissionId int64) error {
	var submission models.Submission
	if err := h.db.Where("id = ?", submissionId).First(&submission).Error; err != nil {
		slog.Error("standings context error", "error", err)
		return err
	}
	if submission.ContestID == nil {
		return nil
	}

	contestID := *submission.ContestID

	tx := h.db.Begin()
	if tx.Error != nil {
		slog.Error("standings tx begin error", "error", tx.Error)
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if already solved
	var existingResult models.ContestProblemResult
	findErr := tx.Where("contest_id = ? AND user_id = ? AND problem_id = ?", contestID, submission.UserID, submission.ProblemID).First(&existingResult).Error
	if findErr == nil && existingResult.IsSolved {
		if err := tx.Commit().Error; err != nil {
			slog.Error("standings commit error", "error", err)
		}
		return nil
	}
	if findErr != nil && findErr != gorm.ErrRecordNotFound {
		tx.Rollback()
		slog.Error("standings check exists error", "error", findErr)
		return findErr
	}

	// Check first blood
	var existingFirstBlood int64
	if err := tx.Model(&models.ContestProblemResult{}).
		Where("contest_id = ? AND problem_id = ? AND is_first_blood = ?", contestID, submission.ProblemID, true).
		Count(&existingFirstBlood).Error; err != nil {
		tx.Rollback()
		slog.Error("standings check first blood error", "error", err)
		return err
	}

	isFirstBlood := false
	if existingFirstBlood == 0 {
		var earlierAC int64
		if err := tx.Model(&models.Submission{}).
			Where("contest_id = ? AND problem_id = ? AND status = ? AND created_at < ?",
				contestID, submission.ProblemID, "ACCEPTED", submission.CreatedAt).
			Count(&earlierAC).Error; err != nil {
			tx.Rollback()
			slog.Error("standings check earlier ac error", "error", err)
			return err
		}
		isFirstBlood = earlierAC == 0
	}

	// Calculate penalty
	penalty, err := h.calculatePenalty(tx, contestID, &submission)
	if err != nil {
		tx.Rollback()
		slog.Error("standings penalty error", "error", err)
		return err
	}

	// Count all non-AC submissions before this one
	var wrongAttempts int64
	if err := tx.Model(&models.Submission{}).
		Where("contest_id = ? AND user_id = ? AND problem_id = ? AND created_at < ? AND status != ?",
			contestID, submission.UserID, submission.ProblemID, submission.CreatedAt, "ACCEPTED").
		Count(&wrongAttempts).Error; err != nil {
		tx.Rollback()
		slog.Error("standings count wrong attempts error", "error", err)
		return err
	}

	if findErr == gorm.ErrRecordNotFound {
		result := models.ContestProblemResult{
			ContestId:            contestID,
			UserId:               submission.UserID,
			ProblemId:            submission.ProblemID,
			IsSolved:             true,
			WrongAttempts:        int(wrongAttempts),
			AcceptedSubmissionId: &submissionId,
			SolvedAt:             &submission.CreatedAt,
			Penalty:              penalty,
			IsFirstBlood:         isFirstBlood,
		}
		if err := tx.Create(&result).Error; err != nil {
			tx.Rollback()
			slog.Error("standings create result error", "error", err)
			return err
		}
	} else {
		if err := tx.Model(&existingResult).Updates(map[string]interface{}{
			"is_solved":              true,
			"wrong_attempts":         int(wrongAttempts),
			"accepted_submission_id": submissionId,
			"solved_at":              submission.CreatedAt,
			"penalty":                penalty,
			"is_first_blood":         isFirstBlood,
		}).Error; err != nil {
			tx.Rollback()
			slog.Error("standings update result error", "error", err)
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error("standings commit error", "error", err)
	}
	return nil
}

func (h *Handler) calculatePenalty(tx *gorm.DB, contestID string, submission *models.Submission) (int, error) {
	var wrongCount int64
	err := tx.Model(&models.Submission{}).
		Where("contest_id = ? AND user_id = ? AND problem_id = ? AND created_at < ? AND status IN ?",
			contestID, submission.UserID, submission.ProblemID, submission.CreatedAt,
			[]string{"WRONG_ANSWER", "TIME_LIMIT_EXCEEDED", "RUNTIME_ERROR", "MEMORY_LIMIT_EXCEEDED"}).
		Count(&wrongCount).Error
	if err != nil {
		return 0, err
	}

	var contest models.Contest
	if err := tx.Where("id = ?", contestID).First(&contest).Error; err != nil {
		return 0, err
	}

	elapsed := submission.CreatedAt.Sub(contest.StartTime)
	if elapsed < 0 {
		elapsed = 0
	}

	elapsedMinutes := int(elapsed.Minutes())
	if elapsedMinutes < 0 {
		elapsedMinutes = 0
	}

	return elapsedMinutes + int(wrongCount)*PenaltyPerWrongSubmission, nil
}

func (h *Handler) updateStandingsForNonAccepted(submissionId int64, verdict string) error {
	var submission models.Submission
	if err := h.db.Where("id = ?", submissionId).First(&submission).Error; err != nil {
		slog.Error("standings context error", "error", err)
		return err
	}
	if submission.ContestID == nil {
		return nil
	}

	contestID := *submission.ContestID

	tx := h.db.Begin()
	if tx.Error != nil {
		slog.Error("standings tx begin error", "error", tx.Error)
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check if already solved
	var existingResult models.ContestProblemResult
	findErr := tx.Where("contest_id = ? AND user_id = ? AND problem_id = ?", contestID, submission.UserID, submission.ProblemID).First(&existingResult).Error
	if findErr == nil && existingResult.IsSolved {
		if err := tx.Commit().Error; err != nil {
			slog.Error("standings commit error", "error", err)
		}
		return nil
	}
	if findErr != nil && findErr != gorm.ErrRecordNotFound {
		tx.Rollback()
		slog.Error("standings check exists error", "error", findErr)
		return findErr
	}

	if findErr == gorm.ErrRecordNotFound {
		result := models.ContestProblemResult{
			ContestId:     contestID,
			UserId:        submission.UserID,
			ProblemId:     submission.ProblemID,
			IsSolved:      false,
			WrongAttempts: 1,
			Penalty:       0,
			IsFirstBlood:  false,
		}
		if err := tx.Create(&result).Error; err != nil {
			tx.Rollback()
			slog.Error("standings create result error", "error", err)
			return err
		}
	} else {
		if err := tx.Model(&existingResult).Update("wrong_attempts", gorm.Expr("wrong_attempts + ?", 1)).Error; err != nil {
			tx.Rollback()
			slog.Error("standings update wrong attempts error", "error", err)
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		slog.Error("standings commit error", "error", err)
	}
	return nil
}
