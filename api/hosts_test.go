package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/shoobyban/sshman/backend"
)

// write test for all methods in host.go

func TestGetAllHosts(t *testing.T) {
	// test Hosts.GetAllHosts method
	cfg := backend.NewTestStorage()
	cfg.AddHost(
		&backend.Host{Alias: "host1", Host: "host1.com", User: "user1", Groups: []string{"group1", "group2"}},
		false)
	h := HostsHandler{Prefix: ""}
	// mock http request
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r = r.WithContext(context.WithValue(r.Context(), ConfigKey, cfg))
	h.GetAllHosts(w, r)
	// check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var hosts map[string]backend.Host
	json.NewDecoder(w.Body).Decode(&hosts)
	if len(hosts) != 1 {
		t.Errorf("Expected 1 hosts, got %d: %s", len(hosts), w.Body.String())
	}
	if hosts["host1"].User != "user1" {
		t.Errorf("Expected user1, got %s", hosts["host1"].User)
	}
}

func TestGetHostDetails(t *testing.T) {
	// test Hosts.GetHostDetails method
	cfg := backend.NewTestStorage()
	testHosts := []backend.Host{
		{Alias: "host1", Host: "host1.com", User: "user1", Groups: []string{"group1", "group2"}},
	}
	cfg.AddHost(&testHosts[0], false)
	h := HostsHandler{Prefix: ""}
	// mock http request
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/host1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "host1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	r = r.WithContext(context.WithValue(r.Context(), ConfigKey, cfg))
	h.GetHostDetails(w, r)
	// check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var host backend.Host
	json.NewDecoder(w.Body).Decode(&host)
	if host.User != testHosts[0].User {
		t.Errorf("Expected %s, got %s", testHosts[0].User, host.User)
	}
}

func TestUpdateHost(t *testing.T) {
	// test Hosts.UpdateHost method
	cfg := backend.NewTestStorage()
	cfg.AddHost(
		&backend.Host{Alias: "host1", Host: "host1.com", User: "user1", Groups: []string{"group1", "group2"}},
		false)
	// mock http request
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/host1",
		strings.NewReader(`{"alias": "host1", "host": "host1.com", "user": "user2", "groups": ["group1", "group2"]}`),
	)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "host1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	r = r.WithContext(context.WithValue(r.Context(), ConfigKey, cfg))
	h := HostsHandler{Prefix: ""}
	h.UpdateHost(w, r)
	// check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var host backend.Host
	json.NewDecoder(w.Body).Decode(&host)
	if host.User != "user2" {
		t.Errorf("Expected user2, got %s %#v", host.User, host)
	}
}

// test Hosts.CreateHost method
func TestCreateHost(t *testing.T) {
	cfg := backend.NewTestStorage()
	h := HostsHandler{Prefix: ""}
	// mock http request
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/",
		strings.NewReader(`{"alias": "host1", "host": "host1.com", "keyfile": "ssh dummy key", "user": "user2", "groups": ["group1", "group2"]}`),
	)
	r = r.WithContext(context.WithValue(r.Context(), ConfigKey, cfg))
	h.CreateHost(w, r)
	// check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var host backend.Host
	json.NewDecoder(w.Body).Decode(&host)
	if host.User != "user2" {
		t.Errorf("Expected user2, got %s %#v", host.User, host)
	}
}

// test Hosts.DeleteHost method
func TestDeleteHost(t *testing.T) {
	testHosts := map[string]*backend.Host{
		"host1": {Alias: "host1", Host: "host1.com", User: "user1", Groups: []string{"group1", "group2"}},
	}
	cfg := backend.NewTestStorage()
	cfg.AddHost(testHosts["host1"], false)
	h := HostsHandler{Prefix: ""}
	// mock http request
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/host1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "host1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	r = r.WithContext(context.WithValue(r.Context(), ConfigKey, cfg))
	h.DeleteHost(w, r)
	// check response
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, w.Code)
	}
	if len(cfg.Hosts()) != 0 {
		t.Errorf("Number of hosts is %d after deleting", len(cfg.Hosts()))
	}
}
