package compilerun

import "github.com/judgenot0/judge-backend/config"

func NewHandler(config *config.Config) *Handler {
	return &Handler{
		config: config,
	}
}
