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
	router.Get(prefix+"/{id}", h.GetUserDetails)
	router.Post(prefix, h.CreateUser)
	router.Put(prefix+"/{id}", h.UpdateUser)
	router.Delete(prefix+"/{id}", h.DeleteUser)

	return router
}

func (h *Users) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users := h.Config.GetUsers("")
	json.NewEncoder(w).Encode(users)
}

func (h *Users) GetUserDetails(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := h.Config.GetUser(id)
	json.NewEncoder(w).Encode(user)
}

func (h *Users) CreateUser(w http.ResponseWriter, r *http.Request) {
	h.UpdateUser(w, r)
}

func (h *Users) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user []string
	json.NewDecoder(r.Body).Decode(&user)
	oldUser := h.Config.GetUser(user[0])
	u, err := h.Config.PrepareUser(user...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	u.UpdateGroups(h.Config, oldUser.Groups)
	json.NewEncoder(w).Encode(user)
}

func (h *Users) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.Config.DeleteUser(id)
	json.NewEncoder(w).Encode(id)
}
