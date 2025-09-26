package contest

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/judgenot0/judge-backend/utils"
)

type Contest struct {
	ContestId   int    `json:"contestId"`
	ContestName string `json:"contestName"`
	StartTime   uint64 `json:"startTime"`
	EndTime     uint64 `json:"endTime"`
	Duration    uint64 `json:"duration"`
	Status      string `json:"status"`
}

type Handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) CreateContest(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) UpdateContest(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) ListContest(w http.ResponseWriter, r *http.Request) {
	contests := []Contest{}
	//TODO: Add Dynamic DB fetch of contests
	utils.SendResopnse(w, http.StatusOK, contests)
}

func (h *Handler) GetContest(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) DeleteContest(w http.ResponseWriter, r *http.Request) {

}
