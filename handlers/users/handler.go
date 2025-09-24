package users

import "github.com/judgenot0/judge-backend/config"

type Handler struct {
	config *config.Config
}

func NewHandler(config *config.Config) *Handler {
	return &Handler{
		config: config,
	}
}
