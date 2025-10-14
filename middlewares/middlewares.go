package middlewares

import "github.com/judgenot0/judge-backend/config"

type Middlewares struct {
	config *config.Config
}

func NewMiddlewares(config *config.Config) *Middlewares {
	return &Middlewares{
		config: config,
	}
}
