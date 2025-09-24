package middlewares

import (
	"net/http"

	"github.com/judgenot0/judge-backend/config"
)

type Middleware func(next http.Handler) http.Handler
type Handler func(w http.ResponseWriter, r *http.Request)
type Manager struct {
	globalMiddlewares []Middleware
	config            *config.Config
}

func NewManager(config *config.Config) *Manager {
	return &Manager{
		globalMiddlewares: make([]Middleware, 0),
		config:            config,
	}
}

// To Add Default Middlewares that are needed for all request-respone (e.g. CORS, Preflight Request, Handle Auth, Logger etc)
func (mngr *Manager) Use(middlewares ...Middleware) {
	mngr.globalMiddlewares = append(mngr.globalMiddlewares, middlewares...)
}

// To Chain all the middleware with middleware parametes for clean code
func (mngr *Manager) With(next Handler, middlewares ...Middleware) http.Handler {
	n := http.Handler(http.HandlerFunc(next))

	for i := len(middlewares) - 1; i >= 0; i-- {
		n = middlewares[i](n)
	}

	return n
}

func (mngr *Manager) WrapMux(handler http.Handler) http.Handler {
	h := handler

	for i := len(mngr.globalMiddlewares) - 1; i >= 0; i-- {
		h = mngr.globalMiddlewares[i](h)
	}

	return h
}
