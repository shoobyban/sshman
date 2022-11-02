package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/shoobyban/sshman/backend"
)

// UsersHandler is the handler for users
type UsersHandler struct {
	Prefix string
}

// Config returns the config for the handler
func (h UsersHandler) Config(r *http.Request) *backend.Storage {
	ctx := r.Context()
	if cfg, ok := ctx.Value(ConfigKey).(*backend.Storage); ok {
		return cfg
	}
	return &backend.Storage{}
}

// AddRoutes adds the routes for the handler
func (h UsersHandler) AddRoutes(router *chi.Mux) {
	router.Get(h.Prefix, h.GetAllUsers)
	router.Get(h.Prefix+"/{id}", h.GetUserDetails)
	router.Delete(h.Prefix+"/{id}", h.DeleteUser)
	router.Put(h.Prefix+"/{id}", h.UpdateUser)
	router.Post(h.Prefix, h.CreateUser)
}

// GetAllUsers returns all users
func (h UsersHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users := h.Config(r).Users()
	json.NewEncoder(w).Encode(users)
}

// GetUserDetails returns the details of a user
func (h UsersHandler) GetUserDetails(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := h.Config(r).GetUser(id)
	json.NewEncoder(w).Encode(user)
}

// CreateUser creates a new user
func (h UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user backend.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cfg := h.Config(r)
	if user.File != "" {
		parts, err := backend.SplitParts(user.File)
		if err != nil {
			cfg.Log().Errorf("Invalid key format: %v", err)
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
	if user.Email == "" || user.Key == "" || user.KeyType == "" || user.Name == "" {
		cfg.Log().Errorf("Missing required fields: email: '%s' key: '%s' keytype: '%s' name: '%s'", user.Email, user.Key, user.KeyType, user.Name)
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}
	_, oldUser := cfg.GetUserByEmail(user.Email)
	if oldUser != nil {
		http.Error(w, "user already exists", http.StatusBadRequest)
		return
	}
	cfg.AddUser(&user, "")
	user.UpdateGroups(cfg, []string{})
	json.NewEncoder(w).Encode(user)
}

// UpdateUser updates a user
func (h UsersHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
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
		cfg.Log().Infof("Error decoding user: %v (%s)", err, string(bodyBytes))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user.File != "" {
		cfg.Log().Infof("New keyfile: %s", user.File)
		parts, err := backend.SplitParts(user.File)
		if err != nil {
			cfg.Log().Infof("Error splitting key: %v", err)
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
		cfg.Log().Infof("User %s does not exist: %v", user.Email, cfg.Users())
		http.Error(w, "user does not exist", http.StatusBadRequest)
		return
	}
	user.UpdateGroups(cfg, oldUser.Groups)
	oldUser.Email = user.Email
	oldUser.Name = user.Name
	oldUser.Key = user.Key
	oldUser.KeyType = user.KeyType
	oldUser.Groups = user.Groups
	diff := backend.Difference(oldUser.Hosts, user.Hosts)
	for _, removedAlias := range diff[0] {
		removedHost := cfg.GetHost(removedAlias)
		removedHost.RemoveUser(&user)
	}
	for _, addedAlias := range diff[1] {
		addedHost := cfg.GetHost(addedAlias)
		addedHost.AddUser(&user)
	}
	oldUser.Hosts = user.Hosts
	cfg.UpdateUser(oldUser)
	cfg.Write()
	json.NewEncoder(w).Encode(user)
}

// DeleteUser deletes a user
func (h UsersHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	log.Printf("Deleting user %s", id)
	cfg := h.Config(r)
	user := cfg.GetUser(id)
	if user == nil {
		http.Error(w, "user does not exist", http.StatusBadRequest)
		return
	}
	if h.Config(r).DeleteUserByID(id) {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
