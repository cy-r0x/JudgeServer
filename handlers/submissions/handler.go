package submissions

import (
	"github.com/judgenot0/judge-backend/config"
	"github.com/judgenot0/judge-backend/infra/queue"
	"gorm.io/gorm"
)

type Handler struct {
	db          *gorm.DB
	config      *config.Config
	queueClient *queue.Queue
}

func NewHandler(db *gorm.DB, config *config.Config, queueClient *queue.Queue) *Handler {
	return &Handler{
		db:          db,
		config:      config,
		queueClient: queueClient,
	}
}
