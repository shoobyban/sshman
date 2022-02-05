package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/shoobyban/sshman/backend"
)

type Hosts struct {
	Prefix string
}

func (h Hosts) Config(r *http.Request) *backend.Storage {
	ctx := r.Context()
	if cfg, ok := ctx.Value("config").(*backend.Storage); ok {
		return cfg
	}
	return &backend.Storage{}
}

func (h Hosts) AddRoutes(router *chi.Mux) {
	router.Get(h.Prefix, h.GetAllHosts)
	router.Get(h.Prefix+"/{id}", h.GetHostDetails)
	router.Delete(h.Prefix+"/{id}", h.DeleteHost)
	router.Put(h.Prefix+"/{id}", h.UpdateHost)
	router.Post(h.Prefix, h.CreateHost)
}

func (h Hosts) GetAllHosts(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(h.Config(r).Hosts())
}

func (h Hosts) GetHostDetails(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	json.NewEncoder(w).Encode(h.Config(r).GetHost(id))
}

func (h Hosts) CreateHost(w http.ResponseWriter, r *http.Request) {
	h.UpdateHost(w, r)
}

func (h Hosts) UpdateHost(w http.ResponseWriter, r *http.Request) {
	var host backend.Host
	cfg := h.Config(r)
	err := json.NewDecoder(r.Body).Decode(&host)
	if err != nil {
		cfg.Log.Errorf("Can't decode host %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var oldHost *backend.Host
	oldHost = cfg.GetHost(host.Alias)
	if oldHost == nil { // for CreateHost handler
		oldHost = &backend.Host{}
	}
	cfg.SetHost(host.Alias, &host)
	host.UpdateGroups(cfg, oldHost.Groups)
	cfg.Write()
	json.NewEncoder(w).Encode(host)
}

func (h Hosts) DeleteHost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	cfg := h.Config(r)
	cfg.Log.Infof("Deleting host by alias: %s", id)
	if cfg.DeleteHost(id) {
		cfg.Log.Infof("Deleted host %s", id)
		w.WriteHeader(http.StatusNoContent)
	} else {
		cfg.Log.Errorf("No such host: %s", id)
		w.WriteHeader(http.StatusNotFound)
	}
}
