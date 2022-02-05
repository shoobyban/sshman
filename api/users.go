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
}

func (h Users) Config(r *http.Request) *backend.Storage {
	ctx := r.Context()
	if cfg, ok := ctx.Value("config").(*backend.Storage); ok {
		return cfg
	}
	return &backend.Storage{}
}

func (h Users) AddRoutes(router *chi.Mux) {
	router.Get(h.Prefix, h.GetAllUsers)
	router.Get(h.Prefix+"/{id}", h.GetUserDetails)
	router.Delete(h.Prefix+"/{id}", h.DeleteUser)
	router.Put(h.Prefix+"/{id}", h.UpdateUser)
	router.Post(h.Prefix, h.CreateUser)
}

func (h Users) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users := h.Config(r).Users()
	json.NewEncoder(w).Encode(users)
}

func (h Users) GetUserDetails(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := h.Config(r).GetUser(id)
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
	cfg := h.Config(r)
	_, oldUser := cfg.GetUserByEmail(user.Email)
	if oldUser != nil {
		http.Error(w, "user already exists", http.StatusBadRequest)
		return
	}
	user.UpdateGroups(cfg, []string{})
	cfg.AddUser(&user)
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
	cfg := h.Config(r)
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		cfg.Log.Infof("Error decoding user: %v (%s)", err, string(bodyBytes))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user.File != "" {
		cfg.Log.Infof("New keyfile: %s", user.File)
		parts, err := backend.SplitParts(user.File)
		if err != nil {
			cfg.Log.Infof("Error splitting key: %v", err)
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
	oldUser := cfg.GetUser(id)
	if oldUser == nil {
		cfg.Log.Infof("User %s does not exist: %v", user.Email, cfg.Users())
		http.Error(w, "user does not exist", http.StatusBadRequest)
		return
	}
	user.UpdateGroups(cfg, oldUser.Groups)
	oldUser.Email = user.Email
	oldUser.Name = user.Name
	oldUser.Key = user.Key
	oldUser.KeyType = user.KeyType
	oldUser.Groups = user.Groups
	cfg.UpdateUser(oldUser)
	cfg.Write()
	json.NewEncoder(w).Encode(user)
}

func (h Users) DeleteUser(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	if h.Config(r).DeleteUser(email) {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
