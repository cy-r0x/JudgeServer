package problem

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/judgenot0/judge-backend/middlewares"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) UpdateProblem(w http.ResponseWriter, r *http.Request) {
	// Only Setter can update the problem
	payload, ok := r.Context().Value("user").(*middlewares.Payload)
	if !ok {
		utils.SendResponse(w, http.StatusUnauthorized, "User information not found")
		return
	}

	var updatedProblem Problem
	if err := json.NewDecoder(r.Body).Decode(&updatedProblem); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// For setters, verify they created this problem
	var createdBy int64
	err := h.db.QueryRow(`SELECT created_by FROM problems WHERE id = $1`, updatedProblem.Id).Scan(&createdBy)
	if err != nil {
		log.Println("Error checking problem creator:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to verify problem creator")
		return
	}

	if createdBy != payload.Sub {
		utils.SendResponse(w, http.StatusForbidden, "You can only update problems you've created")
		return
	}

	// Start a transaction for the update operation
	tx, err := h.db.Beginx()
	if err != nil {
		log.Println("Error starting transaction:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to update problem")
		return
	}

	// Update the problem's main details
	_, err = tx.Exec(`
		UPDATE problems 
		SET title = $1, statement = $2, input_statement = $3, output_statement = $4,
		    time_limit = $5, memory_limit = $6, slug = $7
		WHERE id = $8
	`,
		updatedProblem.Title,
		updatedProblem.Statement,
		updatedProblem.InputStatement,
		updatedProblem.OutputStatement,
		updatedProblem.TimeLimit,
		updatedProblem.MemoryLimit,
		strings.ReplaceAll(strings.ToLower(updatedProblem.Title), " ", "-"),
		updatedProblem.Id)

	if err != nil {
		tx.Rollback()
		log.Println("Error updating problem:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to update problem")
		return
	}
	if err := tx.Commit(); err != nil {
		log.Println("Error committing transaction:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to update problem")
		return
	}

	r.SetPathValue("problemId", strconv.FormatInt(updatedProblem.Id, 10))
	h.GetProblem(w, r)
}
