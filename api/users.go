package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/shoobyban/sshman/backend"
)

type Users struct {
	Prefix string
	Config *backend.Storage
}

func (h *Users) Routers(prefix string, router *chi.Mux) *chi.Mux {
	router.Get(prefix, h.GetAllUsers)
	router.Get(prefix+"/{email}", h.GetUserDetails)
	router.Delete(prefix+"/{email}", h.DeleteUser)
	router.Put(prefix+"/{email}", h.UpdateUser)
	router.Post(prefix, h.CreateUser)

	return router
}

func (h *Users) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users := h.Config.GetUsers("")
	json.NewEncoder(w).Encode(users)
}

func (h *Users) GetUserDetails(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	_, user := h.Config.GetUserByEmail(email)
	json.NewEncoder(w).Encode(user)
}

func (h *Users) CreateUser(w http.ResponseWriter, r *http.Request) {
	h.UpdateUser(w, r)
}

func (h *Users) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user []string
	json.NewDecoder(r.Body).Decode(&user)
	var oldUser *backend.User
	var exists bool
	if oldUser, exists = h.Config.Users[user[0]]; !exists {
		oldUser = &backend.User{}
	}
	u, err := h.Config.PrepareUser(user...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	u.UpdateGroups(h.Config, oldUser.Groups)
	json.NewEncoder(w).Encode(u)
}

func (h *Users) DeleteUser(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	if h.Config.DeleteUser(email) {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
