package problem

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/judgenot0/judge-backend/middlewares"
	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) CreateProblem(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	payload, ok := r.Context().Value("user").(*middlewares.Payload)
	if !ok {
		utils.SendResponse(w, http.StatusUnauthorized, "User information not found")
		return
	}

	var reqProblem Problem
	err := decoder.Decode(&reqProblem)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	reqProblem.Statement = ""
	reqProblem.InputStatement = ""
	reqProblem.OutputStatement = ""
	reqProblem.TimeLimit = 1.0
	reqProblem.MemoryLimit = 256.0
	reqProblem.CheckerStrictSpace = false
	reqProblem.CheckerType = "string"
	reqProblem.CheckerPrecision = nil
	reqProblem.CreatedBy = payload.Sub
	reqProblem.Slug = strings.ReplaceAll(strings.ToLower(reqProblem.Title), " ", "-")
	reqProblem.CreatedAt = time.Now()

	newProblem := models.Problem{
		Title:              reqProblem.Title,
		Slug:               reqProblem.Slug,
		Statement:          reqProblem.Statement,
		InputStatement:     reqProblem.InputStatement,
		OutputStatement:    reqProblem.OutputStatement,
		TimeLimit:          float64(reqProblem.TimeLimit),
		MemoryLimit:        float64(reqProblem.MemoryLimit),
		CheckerType:        reqProblem.CheckerType,
		CheckerStrictSpace: reqProblem.CheckerStrictSpace,
		CheckerPrecision:   reqProblem.CheckerPrecision,
		CreatedByID:        &reqProblem.CreatedBy,
		CreatedAt:          reqProblem.CreatedAt,
	}

	err = h.db.Create(&newProblem).Error
	if err != nil {
		log.Println("Error creating problem:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to create problem")
		return
	}

	reqProblem.Id = newProblem.ID
	reqProblem.CreatedAt = newProblem.CreatedAt
	reqProblem.Testcases = []Testcase{}

	utils.SendResponse(w, http.StatusCreated, reqProblem)
}
