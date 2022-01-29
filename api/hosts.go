package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/shoobyban/sshman/backend"
)

type Hosts struct {
	Prefix string
	Config *backend.Storage
}

func (h Hosts) AddRoutes(router *chi.Mux) {
	router.Get(h.Prefix, h.GetAllHosts)
	router.Get(h.Prefix+"/{id}", h.GetHostDetails)
	router.Delete(h.Prefix+"/{id}", h.DeleteHost)
	router.Put(h.Prefix+"/{id}", h.UpdateHost)
	router.Post(h.Prefix, h.CreateHost)
}

func (h Hosts) GetAllHosts(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(h.Config.Hosts())
}

func (h Hosts) GetHostDetails(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	json.NewEncoder(w).Encode(h.Config.GetHost(id))
}

func (h Hosts) CreateHost(w http.ResponseWriter, r *http.Request) {
	h.UpdateHost(w, r)
}

func (h Hosts) UpdateHost(w http.ResponseWriter, r *http.Request) {
	var host backend.Host
	err := json.NewDecoder(r.Body).Decode(&host)
	if err != nil {
		h.Config.Log.Errorf("Can't decode host %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var oldHost *backend.Host
	oldHost = h.Config.GetHost(host.Alias)
	if oldHost == nil { // for CreateHost handler
		oldHost = &backend.Host{}
	}
	h.Config.SetHost(host.Alias, &host)
	host.UpdateGroups(h.Config, oldHost.Groups)
	h.Config.Write()
	json.NewEncoder(w).Encode(host)
}

func (h Hosts) DeleteHost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if h.Config.DeleteHost(id) {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
