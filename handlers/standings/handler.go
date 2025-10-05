package standings

import (
	"database/sql"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/judgenot0/judge-backend/utils"
)

type ProblemStatus struct {
	ProblemId     int64      `json:"problem_id"`
	ProblemIndex  int        `json:"problem_index"`
	Solved        bool       `json:"solved"`
	FirstSolvedAt *time.Time `json:"first_solved_at,omitempty"`
	Attempts      int        `json:"attempts"`
	Penalty       int        `json:"penalty"`
	FirstBlood    bool       `json:"first_blood"`
}

type UserStanding struct {
	UserId       int64           `json:"user_id"`
	Username     string          `json:"username"`
	TotalPenalty int             `json:"total_penalty"`
	SolvedCount  int             `json:"solved_count"`
	Problems     []ProblemStatus `json:"problems"`
	LastSolvedAt *time.Time      `json:"last_solved_at,omitempty"`
}

type StandingsResponse struct {
	ContestId         int64          `json:"contest_id"`
	ContestTitle      string         `json:"contest_title"`
	TotalProblemCount int            `json:"total_problem_count"`
	Standings         []UserStanding `json:"standings"`
}

type Handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) GetStandings(w http.ResponseWriter, r *http.Request) {
	contestIdStr := r.PathValue("contestId")
	if contestIdStr == "" {
		utils.SendResponse(w, http.StatusBadRequest, "Contest ID is required")
		return
	}

	contestId, err := strconv.ParseInt(contestIdStr, 10, 64)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid contest ID")
		return
	}

	// Get all problems in the contest
	type ContestProblem struct {
		ProblemId int64 `db:"problem_id"`
		Index     int   `db:"index"`
	}

	var contestProblems []ContestProblem
	err = h.db.Select(&contestProblems, `
		SELECT problem_id, index 
		FROM contest_problems 
		WHERE contest_id = $1 
		ORDER BY index ASC
	`, contestId)
	if err != nil {
		log.Println("Error fetching contest problems:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch contest problems")
		return
	}

	// Fetch contest title
	var contestTitle string
	if err := h.db.Get(&contestTitle, `SELECT title FROM contests WHERE id = $1`, contestId); err != nil {
		log.Println("Error fetching contest title:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch contest title")
		return
	}

	// Aggregate per-user/problem statistics in a single query
	type standingsRow struct {
		UserId       int64         `db:"user_id"`
		Username     string        `db:"username"`
		ProblemId    int64         `db:"problem_id"`
		ProblemIndex int           `db:"problem_index"`
		Attempts     int           `db:"attempts"`
		SolvedAt     sql.NullTime  `db:"solved_at"`
		Penalty      sql.NullInt64 `db:"penalty"`
		FirstBlood   sql.NullBool  `db:"first_blood"`
	}

	rows := []standingsRow{}
	err = h.db.Select(&rows, `
		WITH participants AS (
			SELECT DISTINCT u.id AS user_id, u.username
			FROM submissions s
			JOIN users u ON u.id = s.user_id
			WHERE s.contest_id = $1
		),
		problems AS (
			SELECT problem_id, index
			FROM contest_problems
			WHERE contest_id = $1
		)
		SELECT
			u.user_id,
			u.username,
			p.problem_id,
			p.index AS problem_index,
			cs.solved_at,
			cs.penalty,
			cs.first_blood,
			COALESCE(SUM(CASE WHEN s.id IS NOT NULL AND (cs.solved_at IS NULL OR s.submitted_at <= cs.solved_at) THEN 1 ELSE 0 END), 0) AS attempts
		FROM participants u
		CROSS JOIN problems p
		LEFT JOIN contest_solves cs
			ON cs.contest_id = $1 AND cs.user_id = u.user_id AND cs.problem_id = p.problem_id
		LEFT JOIN submissions s
			ON s.contest_id = $1 AND s.user_id = u.user_id AND s.problem_id = p.problem_id
		GROUP BY u.user_id, u.username, p.problem_id, p.index, cs.solved_at, cs.penalty, cs.first_blood
		ORDER BY u.user_id, p.index
	`, contestId)
	if err != nil {
		log.Println("Error fetching aggregated standings:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch standings")
		return
	}

	userCapacity := 0
	if len(contestProblems) > 0 {
		userCapacity = (len(rows) + len(contestProblems) - 1) / len(contestProblems)
	}
	standings := make([]UserStanding, 0, userCapacity)
	if len(rows) > 0 {
		var current *UserStanding
		var currentUserID int64

		for _, row := range rows {
			if current == nil || row.UserId != currentUserID {
				standings = append(standings, UserStanding{
					UserId:   row.UserId,
					Username: row.Username,
					Problems: make([]ProblemStatus, 0, len(contestProblems)),
				})
				current = &standings[len(standings)-1]
				currentUserID = row.UserId
			}

			penalty := 0
			if row.Penalty.Valid {
				penalty = int(row.Penalty.Int64)
			}

			problemStatus := ProblemStatus{
				ProblemId:    row.ProblemId,
				ProblemIndex: row.ProblemIndex,
				Attempts:     row.Attempts,
				Penalty:      penalty,
				FirstBlood:   row.FirstBlood.Valid && row.FirstBlood.Bool,
			}

			if row.SolvedAt.Valid {
				solvedAt := row.SolvedAt.Time
				problemStatus.Solved = true
				problemStatus.FirstSolvedAt = &solvedAt
				current.TotalPenalty += penalty

				current.SolvedCount++

				if current.LastSolvedAt == nil || solvedAt.After(*current.LastSolvedAt) {
					last := solvedAt
					current.LastSolvedAt = &last
				}
			}

			current.Problems = append(current.Problems, problemStatus)
		}
	}
	// Sort standings: more solves -> lower penalty -> earlier last solve
	sort.SliceStable(standings, func(i, j int) bool {
		if standings[i].SolvedCount != standings[j].SolvedCount {
			return standings[i].SolvedCount > standings[j].SolvedCount
		}
		if standings[i].TotalPenalty != standings[j].TotalPenalty {
			return standings[i].TotalPenalty < standings[j].TotalPenalty
		}
		li, lj := standings[i].LastSolvedAt, standings[j].LastSolvedAt
		switch {
		case li == nil && lj == nil:
			return standings[i].UserId < standings[j].UserId
		case li == nil:
			return false
		case lj == nil:
			return true
		default:
			return li.Before(*lj)
		}
	})

	response := StandingsResponse{
		ContestId:         contestId,
		ContestTitle:      contestTitle,
		TotalProblemCount: len(contestProblems),
		Standings:         standings,
	}

	utils.SendResponse(w, http.StatusOK, response)
}
