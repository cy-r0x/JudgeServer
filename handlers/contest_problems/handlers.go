package contest_problems

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/judgenot0/judge-backend/utils"
)

type ContestProblem struct {
	ContestId     int64  `json:"contest_id" db:"contest_id"`
	ProblemId     int64  `json:"problem_id" db:"problem_id"`
	Index         int    `json:"index" db:"index"`
	ProblemName   string `json:"problem_name,omitempty" db:"problem_name"`
	ProblemAuthor string `json:"problem_author,omitempty" db:"problem_author"`
}

type Handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) GetContestProblems(w http.ResponseWriter, r *http.Request) {
	contestId := r.PathValue("contestId")
	
	// Convert contestId to int64
	contestIdInt, err := strconv.ParseInt(contestId, 10, 64)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid contest ID")
		return
	}

	// Check if contest exists
	var contestExists bool
	if err = h.db.Get(&contestExists, `SELECT EXISTS(SELECT 1 FROM contests WHERE id=$1)`, contestIdInt); err != nil {
		log.Println("Failed to check contest existence:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to get contest problems")
		return
	}
	if !contestExists {
		utils.SendResponse(w, http.StatusNotFound, "Contest does not exist")
		return
	}

	// Get all contest problems with problem details and author information in one optimized query
	var contestProblems []ContestProblem
	query := `
		SELECT 
			cp.contest_id,
			cp.problem_id,
			cp.index,
			p.title as problem_name,
			COALESCE(u.full_name, 'Unknown') as problem_author
		FROM contest_problems cp
		JOIN problems p ON cp.problem_id = p.id
		LEFT JOIN users u ON p.created_by = u.id
		WHERE cp.contest_id = $1
		ORDER BY cp.index ASC
	`

	if err = h.db.Select(&contestProblems, query, contestIdInt); err != nil {
		log.Println("Failed to get contest problems:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to get contest problems")
		return
	}

	utils.SendResponse(w, http.StatusOK, contestProblems)
}

func (h *Handler) AssignContestProblems(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var contestProblem ContestProblem
	err := decoder.Decode(&contestProblem)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	if contestProblem.ContestId == 0 || contestProblem.ProblemId == 0 {
		utils.SendResponse(w, http.StatusBadRequest, "Contest ID and Problem ID are required")
		return
	}

	// Check if contest exists
	var contestExists bool
	if err = h.db.Get(&contestExists, `SELECT EXISTS(SELECT 1 FROM contests WHERE id=$1)`, contestProblem.ContestId); err != nil {
		log.Println("Failed to check contest existence:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to assign contest problem")
		return
	}
	if !contestExists {
		utils.SendResponse(w, http.StatusBadRequest, "Contest does not exist")
		return
	}

	// Get problem details and author information in a single optimized query
	type ProblemDetails struct {
		Title    string `db:"title"`
		FullName string `db:"full_name"`
	}

	var problemDetails ProblemDetails
	query := `
		SELECT p.title, u.full_name 
		FROM problems p 
		LEFT JOIN users u ON p.created_by = u.id 
		WHERE p.id = $1
	`

	if err = h.db.Get(&problemDetails, query, contestProblem.ProblemId); err != nil {
		if err.Error() == "sql: no rows in result set" {
			utils.SendResponse(w, http.StatusBadRequest, "Problem does not exist")
			return
		}
		log.Println("Failed to get problem details:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to assign contest problem")
		return
	}

	// Set problem details in the response struct
	contestProblem.ProblemName = problemDetails.Title
	contestProblem.ProblemAuthor = problemDetails.FullName

	tx, err := h.db.Beginx()
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to assign contest problem")
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Ensure the problem is not already assigned
	var exists bool
	if err = tx.Get(&exists, `SELECT EXISTS(SELECT 1 FROM contest_problems WHERE contest_id=$1 AND problem_id=$2)`, contestProblem.ContestId, contestProblem.ProblemId); err != nil {
		log.Println("Failed to check existing assignment:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to assign contest problem")
		return
	}
	if exists {
		tx.Rollback()
		utils.SendResponse(w, http.StatusConflict, "Problem already assigned to this contest")
		return
	}

	// Count existing problems to determine next index
	var count int
	if err = tx.Get(&count, `SELECT COUNT(*) FROM contest_problems WHERE contest_id=$1`, contestProblem.ContestId); err != nil {
		log.Println("Failed to count contest problems:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to assign contest problem")
		return
	}

	contestProblem.Index = count + 1

	if _, err = tx.Exec(`INSERT INTO contest_problems (contest_id, problem_id, index) VALUES ($1, $2, $3)`, contestProblem.ContestId, contestProblem.ProblemId, contestProblem.Index); err != nil {
		log.Println("Failed to insert contest problem:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to assign contest problem")
		return
	}

	if err = tx.Commit(); err != nil {
		log.Println("Failed to commit contest problem assignment:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to assign contest problem")
		return
	}

	utils.SendResponse(w, http.StatusOK, contestProblem)
}

func (h *Handler) DeleteContestProblem(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var contestProblem ContestProblem
	err := decoder.Decode(&contestProblem)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	if contestProblem.ContestId == 0 || contestProblem.ProblemId == 0 {
		utils.SendResponse(w, http.StatusBadRequest, "Contest ID and Problem ID are required")
		return
	}

	tx, err := h.db.Beginx()
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to delete contest problem")
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	result, execErr := tx.Exec(`DELETE FROM contest_problems WHERE contest_id=$1 AND problem_id=$2`, contestProblem.ContestId, contestProblem.ProblemId)
	if execErr != nil {
		err = execErr
		log.Println("Failed to delete contest problem:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to delete contest problem")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		tx.Rollback()
		utils.SendResponse(w, http.StatusNotFound, "Problem not assigned to this contest")
		return
	}

	var remaining []ContestProblem
	if err = tx.Select(&remaining, `SELECT contest_id, problem_id, index FROM contest_problems WHERE contest_id=$1 ORDER BY index ASC`, contestProblem.ContestId); err != nil {
		log.Println("Failed to fetch remaining contest problems:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to delete contest problem")
		return
	}

	for i, cp := range remaining {
		newIndex := i + 1
		if cp.Index == newIndex {
			continue
		}
		if _, err = tx.Exec(`UPDATE contest_problems SET index=$1 WHERE contest_id=$2 AND problem_id=$3`, newIndex, cp.ContestId, cp.ProblemId); err != nil {
			log.Println("Failed to update contest problem index:", err)
			utils.SendResponse(w, http.StatusInternalServerError, "Failed to delete contest problem")
			return
		}
	}

	if err = tx.Commit(); err != nil {
		log.Println("Failed to commit contest problem deletion:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to delete contest problem")
		return
	}

	utils.SendResponse(w, http.StatusOK, map[string]any{"message": "Contest problem removed"})
}
