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
	if err := json.NewEncoder(w).Encode(groupList(h.Config(r).GetGroups())); err != nil {
		details := err.Error()
		JSONError(w, "Failed to list groups.", details, http.StatusInternalServerError, nil, true)
	}
}

// GetGroupDetails is a handler that returns group details.
func (h GroupsHandler) GetGroupDetails(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := json.NewEncoder(w).Encode(groupDetails(id, h.Config(r).GetGroup(id))); err != nil {
		details := err.Error()
		JSONError(w, "Failed to get group.", details, http.StatusInternalServerError, map[string]interface{}{"group": id}, true)
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
		details := err.Error()
		JSONError(w, "Invalid request body.", details, http.StatusBadRequest, nil, true)
		return
	}
	if group.Label == "" {
		JSONError(w, "Label cannot be empty.", "missing label in group payload", http.StatusBadRequest, nil, true)
		return
	}
	if err := h.Config(r).UpdateGroup(group.Label, group.Hosts, group.Users); err != nil {
		JSONError(w, "Failed to update group.", err.Error(), http.StatusConflict, map[string]interface{}{"group": group.Label}, true)
		return
	}
	if id != "" && group.Label != id {
		if err := h.Config(r).DeleteGroup(id); err != nil {
			JSONError(w, "Failed to rename group.", err.Error(), http.StatusConflict, map[string]interface{}{"group": id}, true)
			return
		}
	}
	if err := json.NewEncoder(w).Encode(group); err != nil {
		details := err.Error()
		JSONError(w, "Failed to encode response.", details, http.StatusInternalServerError, nil, true)
	}
}

// DeleteGroup is a handler that deletes a group.
func (h GroupsHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.Config(r).DeleteGroup(id); err != nil {
		status := http.StatusConflict
		if err.Error() == "group "+id+" not found" {
			status = http.StatusNotFound
		}
		JSONError(w, "Failed to delete group.", err.Error(), status, map[string]interface{}{"group": id}, true)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
