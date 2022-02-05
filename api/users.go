package api

import (
	"encoding/json"
	"io"
	"log"
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
	router.Get(h.Prefix+"/{id}", h.GetUserDetails)
	router.Delete(h.Prefix+"/{id}", h.DeleteUser)
	router.Put(h.Prefix+"/{id}", h.UpdateUser)
	router.Post(h.Prefix, h.CreateUser)
}

func (h Users) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users := h.Config.Users()
	json.NewEncoder(w).Encode(users)
}

func (h Users) GetUserDetails(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := h.Config.GetUser(id)
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

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(bodyBytes) == 0 {
		log.Printf("no body")
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		h.Config.Log.Infof("Error decoding user: %v (%s)", err, string(bodyBytes))
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
	id := chi.URLParam(r, "id")
	oldUser := h.Config.GetUser(id)
	if oldUser == nil {
		h.Config.Log.Infof("User %s does not exist: %v", user.Email, h.Config.Users())
		http.Error(w, "user does not exist", http.StatusBadRequest)
		return
	}
	user.UpdateGroups(h.Config, oldUser.Groups)
	oldUser.Email = user.Email
	oldUser.Name = user.Name
	oldUser.Key = user.Key
	oldUser.KeyType = user.KeyType
	oldUser.Groups = user.Groups
	h.Config.UpdateUser(oldUser)
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
