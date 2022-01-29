package api

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	h.Config.Log.Open(wo)
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
			h.Config.Log.Close(wo)
			return
		}
	}
}
