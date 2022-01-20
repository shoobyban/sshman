package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/shoobyban/sshman/backend"
)

type Hosts struct {
	Prefix string
	Config *backend.Storage
}

func (h *Hosts) Routers(prefix string, router *chi.Mux) *chi.Mux {
	router.Get(prefix, h.GetAllHosts)
	router.Get(prefix+"/{id}", h.GetHostDetails)
	router.Post(prefix, h.CreateHost)
	router.Put(prefix+"/{id}", h.UpdateHost)
	router.Delete(prefix+"/{id}", h.DeleteHost)

	return router
}

func (h *Hosts) GetAllHosts(w http.ResponseWriter, r *http.Request) {
	hosts := h.Config.Hosts
	json.NewEncoder(w).Encode(hosts)
}

func (h *Hosts) GetHostDetails(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	host := h.Config.Hosts[id]
	json.NewEncoder(w).Encode(host)
}

func (h *Hosts) CreateHost(w http.ResponseWriter, r *http.Request) {
	h.UpdateHost(w, r)
}

func (h *Hosts) UpdateHost(w http.ResponseWriter, r *http.Request) {
	var host []string
	json.NewDecoder(r.Body).Decode(&host)
	var oldHost *backend.Host
	var exists bool
	if oldHost, exists = h.Config.Hosts[host[0]]; !exists {
		oldHost = &backend.Host{}
	}
	newHost, err := h.Config.AddHost(host...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	newHost.UpdateGroups(h.Config, oldHost.Groups)
	json.NewEncoder(w).Encode(newHost)
}

func (h *Hosts) DeleteHost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.Config.DeleteHost(id)
	json.NewEncoder(w).Encode(id)
}
