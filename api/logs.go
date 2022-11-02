package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/shoobyban/sshman/backend"
)

// LogsHandler is a struct for streaming logs to frontend
type LogsHandler struct {
	Prefix string
}

// Config returns a loaded configuration for the handler
func (h LogsHandler) Config(r *http.Request) *backend.Storage {
	ctx := r.Context()
	if cfg, ok := ctx.Value(ConfigKey).(*backend.Storage); ok {
		return cfg
	}
	return &backend.Storage{}
}

// AddRoutes adds logs handler specific routes to the router
func (h LogsHandler) AddRoutes(router *chi.Mux) {
	router.Get(h.Prefix, h.GetLogs)
}

// GetLogs streams logs on a keep-alive connection to the frontend
func (h LogsHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	wo := backend.LogWorker{}
	wo.Source = make(chan interface{}, 10)

	h.Config(r).Log().Open(wo)
	for {
		select {
		case logLine := <-wo.Source:
			var buf bytes.Buffer
			enc := json.NewEncoder(&buf)
			enc.Encode(logLine)
			fmt.Fprintf(w, "data: %v\n", buf.String())
			f.Flush()
		case <-wo.Quit:
			return
		case <-r.Context().Done():
			h.Config(r).Log().Close(wo)
			return
		}
	}
}
