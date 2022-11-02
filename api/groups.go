package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/shoobyban/sshman/backend"
)

// GroupsHandler handler
type GroupsHandler struct {
	Prefix string
}

// Config returns a loaded configuration from context
func (h GroupsHandler) Config(r *http.Request) *backend.Storage {
	ctx := r.Context()
	if cfg, ok := ctx.Value(ConfigKey).(*backend.Storage); ok {
		return cfg
	}
	return &backend.Storage{}
}

// AddRoutes adds group specific routes to the router
func (h GroupsHandler) AddRoutes(router *chi.Mux) {
	router.Get(h.Prefix, h.GetAllGroups)
	router.Get(h.Prefix+"/{id}", h.GetGroupDetails)
	router.Delete(h.Prefix+"/{id}", h.DeleteGroup)
	router.Put(h.Prefix+"/{id}", h.UpdateGroup)
	router.Post(h.Prefix, h.CreateGroup)
}

// GetAllGroups handler returns all groups
func (h GroupsHandler) GetAllGroups(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(h.Config(r).GetGroups())
}

// GetGroupDetails handler returns group details
func (h GroupsHandler) GetGroupDetails(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	json.NewEncoder(w).Encode(h.Config(r).GetGroup(id))
}

// CreateGroup handler creates a group and adds users and hosts to it
// format: label => [hosts: [host1, host2], users: [user1, user2]]
func (h GroupsHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	h.UpdateGroup(w, r)
}

// UpdateGroup handler updates a group and updates users and hosts binding
// format: label => [hosts: [host1, host2], users: [user1, user2]]
func (h GroupsHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	var group struct {
		Label string
		Users []string
		Hosts []string
	}
	id := chi.URLParam(r, "id")
	json.NewDecoder(r.Body).Decode(&group)
	h.Config(r).UpdateGroup(group.Label, group.Hosts, group.Users)
	if id != "" && group.Label != id {
		h.Config(r).DeleteGroup(id)
	}
	json.NewEncoder(w).Encode(group)
}

// DeleteGroup handler deletes a group
func (h GroupsHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if h.Config(r).DeleteGroup(id) {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
