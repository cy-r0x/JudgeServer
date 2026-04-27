package problem

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/judgenot0/judge-backend/middlewares"
	"github.com/judgenot0/judge-backend/models"
	"github.com/judgenot0/judge-backend/utils"
)

func (h *Handler) CreateProblem(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	payload, ok := r.Context().Value("user").(*middlewares.Payload)
	if !ok {
		utils.SendResponse(w, http.StatusUnauthorized, "User information not found", nil)
		return
	}

	var reqProblem Problem
	err := decoder.Decode(&reqProblem)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, "Invalid request payload", nil)
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
	reqProblem.Author = payload.Sub
	reqProblem.CreatedAt = time.Now()
	reqProblem.UpdatedAt = time.Now()

	newProblem := models.Problem{
		Title:              reqProblem.Title,
		Statement:          reqProblem.Statement,
		InputStatement:     reqProblem.InputStatement,
		OutputStatement:    reqProblem.OutputStatement,
		TimeLimit:          float64(reqProblem.TimeLimit),
		MemoryLimit:        float64(reqProblem.MemoryLimit),
		CheckerType:        reqProblem.CheckerType,
		CheckerStrictSpace: reqProblem.CheckerStrictSpace,
		CheckerPrecision:   reqProblem.CheckerPrecision,
		Author:             reqProblem.Author,
		CreatedAt:          reqProblem.CreatedAt,
		UpdatedAt:          reqProblem.UpdatedAt,
	}

	err = h.db.Create(&newProblem).Error
	if err != nil {
		log.Println("Error creating problem:", err)
		utils.SendResponse(w, http.StatusInternalServerError, "Failed to create problem", nil)
		return
	}

	reqProblem.Id = newProblem.Id
	reqProblem.CreatedAt = newProblem.CreatedAt
	reqProblem.UpdatedAt = newProblem.UpdatedAt
	reqProblem.Testcases = []Testcase{}

	utils.SendResponse(w, http.StatusCreated, nil, reqProblem)
}