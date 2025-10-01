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

type Solve struct {
	ProblemId    int64     `json:"problem_id"`
	ProblemIndex int       `json:"problem_index"`
	SolvedAt     time.Time `json:"solved_at"`
	Penalty      int       `json:"penalty"`
}

type Standing struct {
	ContestId    int64      `json:"contest_id" db:"contest_id"`
	UserId       int64      `json:"user_id" db:"user_id"`
	Username     string     `json:"username" db:"username"`
	Penalty      int        `json:"penalty" db:"penalty"`
	SolvedCount  int        `json:"solved_count" db:"solved_count"`
	Solved       []Solve    `json:"solved"`
	LastSolvedAt *time.Time `json:"last_solved_at,omitempty"`
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

	standings := []Standing{}
	err = h.db.Select(&standings, `
		SELECT cs.contest_id, cs.user_id, u.username, cs.penalty, cs.solved_count
		FROM contest_standings cs
	JOIN users u ON cs.user_id = u.id
		WHERE cs.contest_id = $1
	`, contestId)
	if err != nil {
		log.Println("Error fetching contest standings:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch standings")
		return
	}

	solveRows := []struct {
		UserId       int64         `db:"user_id"`
		ProblemId    int64         `db:"problem_id"`
		SolvedAt     time.Time     `db:"solved_at"`
		Penalty      int           `db:"penalty"`
		ProblemIndex sql.NullInt64 `db:"problem_index"`
	}{}

	err = h.db.Select(&solveRows, `
		SELECT cs.user_id,
		       cs.problem_id,
		       cs.solved_at,
		       cs.penalty,
		       cp."index" AS problem_index
		FROM contest_solves cs
		LEFT JOIN contest_problems cp ON cp.contest_id = cs.contest_id AND cp.problem_id = cs.problem_id
		WHERE cs.contest_id = $1
		ORDER BY cs.solved_at ASC
	`, contestId)
	if err != nil {
		log.Println("Error fetching contest solve details:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to fetch standings")
		return
	}

	userSolves := make(map[int64][]Solve)
	for _, row := range solveRows {
		index := 0
		if row.ProblemIndex.Valid {
			index = int(row.ProblemIndex.Int64)
		}
		userSolves[row.UserId] = append(userSolves[row.UserId], Solve{
			ProblemId:    row.ProblemId,
			ProblemIndex: index,
			SolvedAt:     row.SolvedAt,
			Penalty:      row.Penalty,
		})
	}

	for i := range standings {
		solves := userSolves[standings[i].UserId]
		if len(solves) > 0 {
			sort.Slice(solves, func(a, b int) bool {
				return solves[a].SolvedAt.Before(solves[b].SolvedAt)
			})
			standings[i].Solved = solves
			standings[i].SolvedCount = len(solves)
			last := solves[len(solves)-1].SolvedAt
			standings[i].LastSolvedAt = &last
		}
	}

	sort.SliceStable(standings, func(i, j int) bool {
		if standings[i].SolvedCount != standings[j].SolvedCount {
			return standings[i].SolvedCount > standings[j].SolvedCount
		}
		if standings[i].Penalty != standings[j].Penalty {
			return standings[i].Penalty < standings[j].Penalty
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

	utils.SendResponse(w, http.StatusOK, standings)
}

func (h *Handler) UpdateStanding(w http.ResponseWriter, r *http.Request) {

}
