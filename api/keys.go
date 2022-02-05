package api

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/shoobyban/sshman/backend"
)

type Keys struct {
	Prefix string
}

func (h Keys) Config(r *http.Request) *backend.Storage {
	ctx := r.Context()
	if cfg, ok := ctx.Value("config").(*backend.Storage); ok {
		return cfg
	}
	return &backend.Storage{}
}

func (h Keys) AddRoutes(router *chi.Mux) {
	router.Get(h.Prefix, h.GetAllKeys)
	// router.Get(h.Prefix+"/{filename}", h.GetKeyDetails)
	// router.Delete(h.Prefix+"/{filename}", h.DeleteKey)
	// router.Put(h.Prefix+"/{filename}", h.UpdateKey)
	router.Post(h.Prefix, h.CreateKey)
}

func (h Keys) GetAllKeys(w http.ResponseWriter, r *http.Request) {

	t := r.URL.Query().Get("type")
	if t == "" {
		t = "all"
	}
	home, _ := os.UserHomeDir()
	files, err := os.ReadDir(home + "/.ssh")
	if err != nil {
		h.Config(r).Log.Errorf("Can't read ~/.ssh: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

func (h Keys) CreateKey(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Filename string `json:"filename"`
		File     string `json:"file"`
	}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		h.Config(r).Log.Errorf("Can't decode key data %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if payload.Filename == "" {
		http.Error(w, "No filename provided", http.StatusBadRequest)
		return
	}

	if payload.File == "" {
		http.Error(w, "No file provided", http.StatusBadRequest)
		return
	}
	// if file exists
	if _, err := os.Stat(os.Getenv("HOME") + "/.ssh/" + payload.Filename); err == nil {
		http.Error(w, "File already exists", http.StatusBadRequest)
		return
	}
	os.WriteFile(".ssh/"+payload.Filename, []byte(payload.File), 0600)
	w.WriteHeader(http.StatusNoContent)
}
