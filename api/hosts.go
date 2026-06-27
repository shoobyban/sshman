package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/shoobyban/sshman/backend"
)

// HostsHandler is a struct for handling hosts.
type HostsHandler struct {
	Prefix string
}

// Config returns a loaded configuration for the handler.
func (h HostsHandler) Config(r *http.Request) backend.Config {
	ctx := r.Context()
	if cfg, ok := ctx.Value(ConfigKey).(*backend.Data); ok {
		return cfg
	}
	return backend.DefaultConfig()
}

// AddRoutes adds hosthandler specific routes to the router.
func (h HostsHandler) AddRoutes(router *chi.Mux) {
	router.Get(h.Prefix, h.GetAllHosts)
	router.Get(h.Prefix+"/{id}", h.GetHostDetails)
	router.Delete(h.Prefix+"/{id}", h.DeleteHost)
	router.Put(h.Prefix+"/{id}", h.UpdateHost)
	router.Post(h.Prefix, h.CreateHost)
	router.Post(h.Prefix+"/sync", h.SyncHandler)
	router.Delete(h.Prefix+"/sync", h.StopSyncHandler)
}

// GetAllHosts returns all hosts.
func (h HostsHandler) GetAllHosts(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(hostList(h.Config(r).Hosts()))
}

// GetHostDetails returns host details.
func (h HostsHandler) GetHostDetails(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "sync" {
		JSONError(w, "Host not found.", "sync is a reserved endpoint", http.StatusNotFound, map[string]interface{}{"host": id}, true)
		return
	}
	json.NewEncoder(w).Encode(h.Config(r).GetHost(id))
}

// CreateHost creates a new host.
func (h HostsHandler) CreateHost(w http.ResponseWriter, r *http.Request) {
	h.UpdateHost(w, r)
}

// UpdateHost updates a host.
func (h HostsHandler) UpdateHost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var host backend.Host
	cfg := h.Config(r)
	err := json.NewDecoder(r.Body).Decode(&host)
	if err != nil {
		details := err.Error()
		cfg.Log().Errorf("Can't decode host %s", details)
		JSONError(w, "Invalid request body.", details, http.StatusBadRequest, nil, true)
		return
	}
	if host.Alias == "" {
		cfg.Log().Errorf("Can't create host without alias")
		JSONError(w, "Can't create host without alias.", "missing alias field in host payload", http.StatusBadRequest, nil, true)
		return
	}
	if host.Alias == "sync" {
		JSONError(w, "Can't use reserved host alias.", "sync is reserved for the sync endpoint", http.StatusBadRequest, nil, true)
		return
	}
	oldHost := cfg.GetHost(id)
	if dataCfg, ok := cfg.(*backend.Data); ok {
		host.Config = dataCfg
	}
	err = cfg.UpdateHost(&host)
	if err != nil {
		details := err.Error()
		cfg.Log().Errorf("Can't read host users, %s", details)
		JSONError(w, "Failed to update host.", details, http.StatusInternalServerError, map[string]interface{}{"host": host.Alias}, true)
		return
	}
	cfg.SetHost(host.Alias, &host)
	if oldHost != nil {
		if !host.UpdateGroups(cfg, oldHost.Groups) {
			JSONError(w, "Failed to propagate host group changes.", "host group propagation returned partial failure", http.StatusConflict, map[string]interface{}{"host": host.Alias}, true)
			return
		}
		if host.Alias != id {
			cfg.DeleteHost(id)
		}
	} else {
		if !host.UpdateGroups(cfg, []string{}) {
			cfg.DeleteHost(host.Alias)
			JSONError(w, "Failed to propagate host group changes.", "host group propagation returned partial failure", http.StatusConflict, map[string]interface{}{"host": host.Alias}, true)
			return
		}
	}
	if oldHost == nil {
		cfg.Write()
	}
	json.NewEncoder(w).Encode(host)
}

// DeleteHost deletes a host.
func (h HostsHandler) DeleteHost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	cfg := h.Config(r)
	cfg.Log().Infof("Deleting host by alias: %s", id)
	if cfg.DeleteHost(id) {
		cfg.Log().Infof("Deleted host %s", id)
		w.WriteHeader(http.StatusNoContent)
	} else {
		cfg.Log().Errorf("No such host: %s", id)
		w.WriteHeader(http.StatusNotFound)
	}
}

func (h HostsHandler) SyncHandler(_ http.ResponseWriter, r *http.Request) {
	cfg := h.Config(r)
	cfg.Log().Infof("Syncing hosts")
	cfg.Update()
	cfg.Log().Infof("Done syncing")
}

func (h HostsHandler) StopSyncHandler(_ http.ResponseWriter, r *http.Request) {
	cfg := h.Config(r)
	cfg.Log().Infof("Stopping sync")
	cfg.StopUpdate()
	cfg.Log().Infof("Stopped")
}
