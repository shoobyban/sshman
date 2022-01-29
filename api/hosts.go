package api

import (
	"encoding/json"
	"log"
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
	router.Delete(prefix+"/{id}", h.DeleteHost)
	router.Put(prefix+"/{id}", h.UpdateHost)
	router.Post(prefix, h.CreateHost)

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
	var host backend.Host
	err := json.NewDecoder(r.Body).Decode(&host)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var oldHost *backend.Host
	var exists bool
	if oldHost, exists = h.Config.Hosts[host.Alias]; !exists {
		oldHost = &backend.Host{}
	}
	h.Config.Hosts[host.Alias] = &host
	host.UpdateGroups(h.Config, oldHost.Groups)
	h.Config.Write()
	json.NewEncoder(w).Encode(host)
}

func (h *Hosts) DeleteHost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if h.Config.DeleteHost(id) {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
