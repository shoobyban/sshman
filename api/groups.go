package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/shoobyban/sshman/backend"
)

// GroupsHandler is a handler for group-related operations.
type GroupsHandler struct {
	Prefix string
}

// Config returns a loaded configuration from the context.
func (h GroupsHandler) Config(r *http.Request) *backend.Data {
	ctx := r.Context()
	if cfg, ok := ctx.Value(ConfigKey).(*backend.Data); ok {
		return cfg
	}
	return backend.DefaultConfig()
}

// AddRoutes adds group-specific routes to the router.
func (h GroupsHandler) AddRoutes(router *chi.Mux) {
	router.Get(h.Prefix, h.GetAllGroups)
	router.Get(h.Prefix+"/{id}", h.GetGroupDetails)
	router.Delete(h.Prefix+"/{id}", h.DeleteGroup)
	router.Put(h.Prefix+"/{id}", h.UpdateGroup)
	router.Post(h.Prefix, h.CreateGroup)
}

// GetAllGroups is a handler that returns all groups.
func (h GroupsHandler) GetAllGroups(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(h.Config(r).GetGroups()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetGroupDetails is a handler that returns group details.
func (h GroupsHandler) GetGroupDetails(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := json.NewEncoder(w).Encode(h.Config(r).GetGroup(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CreateGroup is a handler that creates a group and adds users and hosts to it.
// format: label => [hosts: [host1, host2], users: [user1, user2]]
func (h GroupsHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	h.UpdateGroup(w, r)
}

// UpdateGroup is a handler that updates a group and its user and host bindings.
// format: label => [hosts: [host1, host2], users: [user1, user2]]
func (h GroupsHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	var group struct {
		Label string
		Users []string
		Hosts []string
	}
	id := chi.URLParam(r, "id")
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if group.Label == "" {
		http.Error(w, "Label cannot be empty", http.StatusBadRequest)
		return
	}
	h.Config(r).UpdateGroup(group.Label, group.Hosts, group.Users)
	if id != "" && group.Label != id {
		h.Config(r).DeleteGroup(id)
	}
	if err := json.NewEncoder(w).Encode(group); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// DeleteGroup is a handler that deletes a group.
func (h GroupsHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if h.Config(r).DeleteGroup(id) {
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(w, "Group not found", http.StatusNotFound)
	}
}
