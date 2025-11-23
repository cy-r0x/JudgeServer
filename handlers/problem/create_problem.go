package problem

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/judgenot0/judge-backend/middlewares"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) CreateProblem(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	payload, ok := r.Context().Value("user").(*middlewares.Payload)
	if !ok {
		utils.SendResponse(w, http.StatusUnauthorized, "User information not found")
		return
	}

	var problem Problem
	err := decoder.Decode(&problem)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	problem.Statement = ""
	problem.InputStatement = ""
	problem.OutputStatement = ""
	problem.TimeLimit = 1
	problem.MemoryLimit = 256
	problem.CheckerStrictSpace = false
	problem.CheckerType = "string"
	problem.CheckerPrecision = nil
	problem.CreatedBy = payload.Sub
	problem.Slug = strings.ReplaceAll(strings.ToLower(problem.Title), " ", "-")
	problem.CreatedAt = time.Now()

	// Insert the Problem into DB and get Problem Id
	tx, err := h.db.Beginx()
	if err != nil {
		log.Println("Error starting transaction:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to create problem")
		return
	}

	var problemID int64
	err = tx.QueryRow(
		`INSERT INTO problems 
		(title, slug, created_by, statement, input_statement, output_statement, time_limit, memory_limit, checker_type, checker_strict_space, checker_precision)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
		RETURNING id`,
		problem.Title, problem.Slug, problem.CreatedBy, problem.Statement, problem.InputStatement, problem.OutputStatement, problem.TimeLimit, problem.MemoryLimit, problem.CheckerType, problem.CheckerStrictSpace, problem.CheckerPrecision,
	).Scan(&problemID)

	if err != nil {
		tx.Rollback()
		log.Println("Error creating problem:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to create problem")
		return
	}

	if err := tx.Commit(); err != nil {
		log.Println("Error committing transaction:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to create problem")
		return
	}

	problem.Id = problemID
	utils.SendResponse(w, http.StatusCreated, problem)
}
