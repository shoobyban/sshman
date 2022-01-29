package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/shoobyban/sshman/backend"
)

type Groups struct {
	Prefix string
	Config *backend.Storage
}

func (h Groups) AddRoutes(router *chi.Mux) {
	router.Get(h.Prefix, h.GetAllGroups)
	router.Get(h.Prefix+"/{id}", h.GetGroupDetails)
	router.Delete(h.Prefix+"/{id}", h.DeleteGroup)
	router.Put(h.Prefix+"/{id}", h.UpdateGroup)
	router.Post(h.Prefix, h.CreateGroup)
}

func (h Groups) GetAllGroups(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(h.Config.GetGroups())
}

func (h Groups) GetGroupDetails(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	json.NewEncoder(w).Encode(h.Config.Groups[id])
}

// format: label => [servers: [server1, server2], users: [user1, user2]]
func (h Groups) CreateGroup(w http.ResponseWriter, r *http.Request) {
	h.UpdateGroup(w, r)
}

// format: label => [servers: [server1, server2], users: [user1, user2]]
func (h Groups) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	var group struct {
		Label   string
		Users   []string
		Servers []string
	}
	json.NewDecoder(r.Body).Decode(&group)
	h.Config.UpdateGroup(group.Label, group.Users, group.Servers)
	json.NewEncoder(w).Encode(group)
}

func (h Groups) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if h.Config.DeleteGroup(id) {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
