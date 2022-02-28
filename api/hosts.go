package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/shoobyban/sshman/backend"
)

// HostsHandler is a struct for handling hosts
type HostsHandler struct {
	Prefix string
}

// Config returns a loaded configuration for the handler
func (h HostsHandler) Config(r *http.Request) *backend.Storage {
	ctx := r.Context()
	if cfg, ok := ctx.Value(ConfigKey).(*backend.Storage); ok {
		return cfg
	}
	return backend.NewStorage(false)
}

// AddRoutes adds hosthandler specific routes to the router
func (h HostsHandler) AddRoutes(router *chi.Mux) {
	router.Get(h.Prefix, h.GetAllHosts)
	router.Get(h.Prefix+"/{id}", h.GetHostDetails)
	router.Delete(h.Prefix+"/{id}", h.DeleteHost)
	router.Put(h.Prefix+"/{id}", h.UpdateHost)
	router.Post(h.Prefix, h.CreateHost)
	router.Get(h.Prefix+"/sync", h.SyncHandler)
	router.Delete(h.Prefix+"/sync", h.StopSyncHandler)
}

// GetAllHosts returns all hosts
func (h HostsHandler) GetAllHosts(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(h.Config(r).Hosts())
}

// GetHostDetails returns host details
func (h HostsHandler) GetHostDetails(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	json.NewEncoder(w).Encode(h.Config(r).GetHost(id))
}

// CreateHost creates a new host
func (h HostsHandler) CreateHost(w http.ResponseWriter, r *http.Request) {
	h.UpdateHost(w, r)
}

// UpdateHost updates a host
func (h HostsHandler) UpdateHost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var host backend.Host
	cfg := h.Config(r)
	err := json.NewDecoder(r.Body).Decode(&host)
	if err != nil {
		cfg.Log.Errorf("Can't decode host %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if host.Alias == "" {
		cfg.Log.Errorf("Can't create host without alias")
		http.Error(w, "Can't create host without alias", http.StatusBadRequest)
		return
	}
	var oldHost *backend.Host
	oldHost = cfg.GetHost(id)
	if oldHost == nil { // for CreateHost handler
		oldHost = &backend.Host{}
	}
	cfg.SetHost(host.Alias, &host)
	host.UpdateGroups(cfg, oldHost.Groups)
	if host.Alias != id {
		cfg.DeleteHost(id)
	}
	cfg.Write()
	err = cfg.UpdateHost(&host)
	if err != nil {
		cfg.Log.Errorf("Can't read host users %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(host)
}

// DeleteHost deletes a host
func (h HostsHandler) DeleteHost(w http.ResponseWriter, r *http.Request) {
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

func (h HostsHandler) SyncHandler(_ http.ResponseWriter, r *http.Request) {
	cfg := h.Config(r)
	cfg.Log.Infof("Syncing hosts")
	cfg.Update()
	cfg.Log.Infof("Done syncing")
}

func (h HostsHandler) StopSyncHandler(_ http.ResponseWriter, r *http.Request) {
	cfg := h.Config(r)
	cfg.Log.Infof("Stopping sync")
	cfg.StopUpdate()
	cfg.Log.Infof("Stopped")
}
