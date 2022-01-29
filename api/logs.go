package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/shoobyban/sshman/backend"
)

type Logs struct {
	Prefix string
	Config *backend.Storage
}

func (h Logs) AddRoutes(router *chi.Mux) {
	router.Get(h.Prefix, h.GetLogs)
}

func (h Logs) GetLogs(w http.ResponseWriter, r *http.Request) {
	for {
		logLine := h.Config.Log.Pop()
		err := json.NewEncoder(w).Encode(logLine)
		if err != nil {
			return
		}
	}
}
