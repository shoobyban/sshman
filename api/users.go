package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/shoobyban/sshman/backend"
)

// UsersHandler is the handler for user-related operations
type UsersHandler struct {
	Prefix string
}

func resolveUserIdentifier(cfg backend.Config, id string) (string, *backend.User) {
	if user := cfg.GetUser(id); user != nil {
		return id, user
	}
	return cfg.GetUserByEmail(id)
}

// Config returns the config for the handler
func (h UsersHandler) Config(r *http.Request) backend.Config {
	ctx := r.Context()
	if cfg, ok := ctx.Value(ConfigKey).(*backend.Data); ok {
		return cfg
	}
	return backend.DefaultConfig()
}

func rollbackCreatedUser(cfg backend.Config, user *backend.User) {
	if user == nil {
		return
	}
	_ = cfg.RemoveUserFromHosts(user)
	if id, existing := cfg.GetUserByEmail(user.Email); existing != nil {
		cfg.DeleteUserByID(id)
	}
}

func cloneUser(user *backend.User) *backend.User {
	if user == nil {
		return nil
	}
	cloned := *user
	cloned.Groups = append([]string{}, user.Groups...)
	cloned.Hosts = append([]string{}, user.Hosts...)
	cloned.Roles = append([]string{}, user.Roles...)
	return &cloned
}

func rollbackUpdatedUser(cfg backend.Config, original *backend.User, current *backend.User) {
	if original == nil || current == nil {
		return
	}
	if err := syncUserHostDiff(cfg, original, current.Hosts, original.Hosts); err != nil {
		cfg.Log().Errorf("failed to roll back user host assignments for %s: %v", original.Email, err)
	}
}

func syncUserHostDiff(cfg backend.Config, user *backend.User, currentHosts, desiredHosts []string) error {
	changes := backend.Difference(currentHosts, desiredHosts)
	var syncErrors []string

	for _, removedAlias := range changes[0] {
		removedHost := cfg.GetHost(removedAlias)
		if removedHost == nil {
			syncErrors = append(syncErrors, fmt.Sprintf("host %s not found", removedAlias))
			continue
		}
		if err := removedHost.RemoveUser(user); err != nil {
			syncErrors = append(syncErrors, fmt.Sprintf("error removing %s from %s: %v", user.Email, removedAlias, err))
		}
	}

	for _, addedAlias := range changes[1] {
		addedHost := cfg.GetHost(addedAlias)
		if addedHost == nil {
			syncErrors = append(syncErrors, fmt.Sprintf("host %s not found", addedAlias))
			continue
		}
		if err := addedHost.AddUser(user); err != nil {
			syncErrors = append(syncErrors, fmt.Sprintf("error adding %s to %s: %v", user.Email, addedAlias, err))
		}
	}

	if len(syncErrors) > 0 {
		return errors.New(strings.Join(syncErrors, "\n"))
	}
	return nil
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
	json.NewEncoder(w).Encode(userList(h.Config(r).Users()))
}

// GetUserDetails returns the details of a user
func (h UsersHandler) GetUserDetails(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, user := resolveUserIdentifier(h.Config(r), id)
	json.NewEncoder(w).Encode(normalizeUser(user))
}

// CreateUser creates a new user
func (h UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user backend.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		details := err.Error()
		JSONError(w, "Invalid request body.", details, http.StatusBadRequest, nil, true)
		return
	}
	cfg := h.Config(r)
	if user.File != "" {
		parts, err := backend.SplitParts(user.File)
		if err != nil {
			details := err.Error()
			cfg.Log().Errorf("Invalid key format: %v", details)
			JSONError(w, "Invalid key format.", details, http.StatusBadRequest, nil, true)
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
		details := fmt.Sprintf("missing fields email:'%s' key:'%s' keytype:'%s' name:'%s'", user.Email, user.Key, user.KeyType, user.Name)
		cfg.Log().Errorf("Missing required fields: %s", details)
		JSONError(w, "Missing required fields.", details, http.StatusBadRequest, nil, true)
		return
	}
	_, oldUser := cfg.GetUserByEmail(user.Email)
		if oldUser != nil {
		JSONError(w, "User already exists.", "user with this email already exists", http.StatusBadRequest, nil, true)
		return
	}
	if err := cfg.AddUser(&user, ""); err != nil {
		JSONError(w, "Failed to create user.", err.Error(), http.StatusBadRequest, nil, true)
		return
	}
	user.Config = cfg
	if err := user.UpdateGroups(cfg, []string{}); err != nil {
		rollbackCreatedUser(cfg, &user)
		JSONError(w, "Failed to propagate user to hosts.", err.Error(), http.StatusConflict, nil, true)
		return
	}
	json.NewEncoder(w).Encode(normalizeUser(&user))
}

// UpdateUser updates a user
func (h UsersHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user backend.User

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		JSONError(w, "Invalid request body.", err.Error(), http.StatusBadRequest, nil, true)
		return
	}
	if len(bodyBytes) == 0 {
		JSONError(w, "Invalid request body.", "empty body", http.StatusBadRequest, nil, true)
		return
	}
	cfg := h.Config(r)
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		cfg.Log().Infof("Error decoding user: %v (%s)", err, string(bodyBytes))
		JSONError(w, "Invalid request body.", err.Error(), http.StatusBadRequest, nil, true)
		return
	}
	if user.File != "" {
		cfg.Log().Infof("New keyfile: %s", user.File)
		parts, err := backend.SplitParts(user.File)
		if err != nil {
			details := err.Error()
			cfg.Log().Infof("Error splitting key: %v", details)
			JSONError(w, "Invalid key format.", details, http.StatusBadRequest, nil, true)
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
	_, oldUser := resolveUserIdentifier(cfg, id)
	if oldUser == nil {
		cfg.Log().Infof("User %s does not exist: %v", user.Email, cfg.Users())
		JSONError(w, "User does not exist.", "no user with given id", http.StatusBadRequest, nil, true)
		return
	}
	oldSnapshot := cloneUser(oldUser)
	user.Config = cfg
	if err := user.UpdateGroups(cfg, oldUser.Groups); err != nil {
		rollbackUpdatedUser(cfg, oldSnapshot, &user)
		JSONError(w, "Failed to propagate user to hosts.", err.Error(), http.StatusConflict, nil, true)
		return
	}
	if err := syncUserHostDiff(cfg, &user, oldUser.Hosts, user.Hosts); err != nil {
		rollbackUpdatedUser(cfg, oldSnapshot, &user)
		JSONError(w, "Failed to update user host assignments.", err.Error(), http.StatusConflict, nil, true)
		return
	}
	oldUser.Email = user.Email
	oldUser.Name = user.Name
	oldUser.Key = user.Key
	oldUser.KeyType = user.KeyType
	oldUser.Groups = user.Groups
	oldUser.Hosts = user.Hosts
	cfg.UpdateUser(oldUser)
	cfg.Write()
	json.NewEncoder(w).Encode(normalizeUser(&user))
}

// DeleteUser deletes a user.
func (h UsersHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	cfg := h.Config(r)
	cfg.Log().Infof("Deleting user %s", id)
	resolvedID, user := resolveUserIdentifier(cfg, id)
	if user == nil {
		JSONError(w, "User does not exist.", "no user with given id", http.StatusBadRequest, nil, true)
		return
	}
	if h.Config(r).DeleteUserByID(resolvedID) {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
