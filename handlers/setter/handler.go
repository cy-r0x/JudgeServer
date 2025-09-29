package setter

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/judgenot0/judge-backend/middlewares"
	"github.com/judgenot0/judge-backend/utils"
)

type Handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) ListSetterProblems(w http.ResponseWriter, r *http.Request) {
	payload, ok := r.Context().Value("user").(*middlewares.Payload)
	if !ok {
		utils.SendResopnse(w, http.StatusUnauthorized, "User information not found")
		return
	}
	setterId := payload.Sub
	log.Println(setterId)
	//TODO: Get All Setter Problems -> Send as Response
}
