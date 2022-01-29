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

func (h Users) AddRoutes(router *chi.Mux) {
	router.Get(h.Prefix, h.GetAllUsers)
	router.Get(h.Prefix+"/{email}", h.GetUserDetails)
	router.Delete(h.Prefix+"/{email}", h.DeleteUser)
	router.Put(h.Prefix+"/{email}", h.UpdateUser)
	router.Post(h.Prefix, h.CreateUser)
}

func (h Users) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users := h.Config.GetUsers("")
	json.NewEncoder(w).Encode(users)
}

func (h Users) GetUserDetails(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	_, user := h.Config.GetUserByEmail(email)
	json.NewEncoder(w).Encode(user)
}

func (h Users) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user backend.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user.File != "" {
		parts, err := backend.SplitParts(user.File)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if user.Email == "" {
			user.Email = parts[2]
		}
		user.KeyType = parts[0]
		user.Key = parts[1]
		user.Name = parts[2]
		user.File = ""
	}
	_, oldUser := h.Config.GetUserByEmail(user.Email)
	if oldUser != nil {
		http.Error(w, "user already exists", http.StatusBadRequest)
		return
	}
	user.UpdateGroups(h.Config, []string{})
	h.Config.AddUser(&user)
	json.NewEncoder(w).Encode(user)
}

func (h Users) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user backend.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		h.Config.Log.Infof("Error decoding user: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user.File != "" {
		h.Config.Log.Infof("New keyfile: %s", user.File)
		parts, err := backend.SplitParts(user.File)
		if err != nil {
			h.Config.Log.Infof("Error splitting key: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if user.Email == "" {
			user.Email = parts[2]
		}
		user.KeyType = parts[0]
		user.Key = parts[1]
		user.Name = parts[2]
		user.File = ""
	}
	_, oldUser := h.Config.GetUserByEmail(user.Email)
	if oldUser == nil {
		h.Config.Log.Infof("User %s does not exist: %v", user.Email, h.Config.Users())
		http.Error(w, "user does not exist", http.StatusBadRequest)
		return
	}
	user.UpdateGroups(h.Config, oldUser.Groups)
	h.Config.UpdateUser(&user)
	h.Config.Write()
	json.NewEncoder(w).Encode(user)
}

func (h Users) DeleteUser(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	if h.Config.DeleteUser(email) {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
