package api

import (
	"encoding/json"
	"log"
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

func (h *Users) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user backend.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Error decoding user: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user.File != "" {
		log.Printf("New keyfile: %s", user.File)
		parts, err := backend.SplitParts(user.File)
		if err != nil {
			log.Printf("Error splitting key: %v", err)
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
		log.Printf("User %s does not exist: %v", user.Email, h.Config.Users)
		http.Error(w, "user does not exist", http.StatusBadRequest)
		return
	}
	user.UpdateGroups(h.Config, oldUser.Groups)
	h.Config.UpdateUser(user)
	h.Config.Write()
	json.NewEncoder(w).Encode(user)
}

func (h *Users) DeleteUser(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	if h.Config.DeleteUser(email) {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
