package api

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/shoobyban/sshman/backend"
)

// KeysHandler returns a list of ssh key filenames from ~/.ssh.
type KeysHandler struct {
	Prefix string
}

// Config returns a loaded configuration for the handler.
func (h KeysHandler) Config(r *http.Request) backend.Config {
	ctx := r.Context()
	if cfg, ok := ctx.Value(ConfigKey).(*backend.Data); ok {
		return cfg
	}
	return backend.DefaultConfig()
}

// AddRoutes adds keyhandler specific routes to the router.
func (h KeysHandler) AddRoutes(router *chi.Mux) {
	router.Get(h.Prefix, h.GetAllKeys)
	// router.Get(h.Prefix+"/{filename}", h.GetKeyDetails)
	// router.Delete(h.Prefix+"/{filename}", h.DeleteKey)
	// router.Put(h.Prefix+"/{filename}", h.UpdateKey)
	router.Post(h.Prefix, h.CreateKey)
}

// GetAllKeys returns a list of all ssh key files from ~/.ssh filtered by type={all|public|private}.
func (h KeysHandler) GetAllKeys(w http.ResponseWriter, r *http.Request) {

	t := r.URL.Query().Get("type")
	if t == "" {
		t = "all"
	}
	home, _ := os.UserHomeDir()
	files, err := os.ReadDir(home + "/.ssh")
	if err != nil {
		details := err.Error()
		h.Config(r).Log().Errorf("Can't read ~/.ssh: %s", details)
		JSONError(w, "Failed to read ~/.ssh directory.", details, http.StatusInternalServerError, nil, true)
		return
	}
	var keys []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		switch t {
		case "all":
			keys = append(keys, "~/.ssh/"+file.Name())
		case "public":
			if len(file.Name()) > 4 && file.Name()[len(file.Name())-4:] == ".pub" {
				keys = append(keys, "~/.ssh/"+file.Name())
			}
		case "private":
			if len(file.Name()) < 4 || file.Name()[len(file.Name())-4:] != ".pub" {
				keys = append(keys, "~/.ssh/"+file.Name())
			}
		}
	}
	e := json.NewEncoder(w)
	if r.URL.Query().Get("pretty") != "" {
		e.SetIndent("", "  ")
	}
	e.Encode(keys)
}

// CreateKey creates a new ssh key from an uploaded file.
func (h KeysHandler) CreateKey(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Filename string `json:"filename"`
		File     string `json:"file"`
	}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		details := err.Error()
		h.Config(r).Log().Errorf("Can't decode key data %s", details)
		JSONError(w, "Invalid request body.", details, http.StatusBadRequest, nil, true)
		return
	}
	if payload.Filename == "" {
		JSONError(w, "No filename provided.", "missing filename in payload", http.StatusBadRequest, nil, true)
		return
	}

	if payload.File == "" {
		JSONError(w, "No file provided.", "missing file content in payload", http.StatusBadRequest, nil, true)
		return
	}
	// if file exists
	if _, err := os.Stat(os.Getenv("HOME") + "/.ssh/" + payload.Filename); err == nil {
		JSONError(w, "File already exists.", "destination file already exists in ~/.ssh", http.StatusBadRequest, nil, true)
		return
	}
	os.WriteFile(".ssh/"+payload.Filename, []byte(payload.File), 0600)
	w.WriteHeader(http.StatusNoContent)
}
