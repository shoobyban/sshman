package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/shoobyban/sshman/backend"
)

// write test for all methods in host.go

func TestGetAllHosts(t *testing.T) {
	// test Hosts.GetAllHosts method
	testHosts := map[string]*backend.Host{
		"host1": {Alias: "host1", Host: "host1.com", User: "user1", Groups: []string{"group1", "group2"}},
	}
	cfg := &backend.Storage{
		Hosts: testHosts}
	h := Hosts{Prefix: "", Config: cfg}
	// mock http request
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	h.GetAllHosts(w, r)
	// check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	bs := w.Body.String()
	log.Printf("body %s", bs)
	var hosts map[string]backend.Host
	json.NewDecoder(w.Body).Decode(&hosts)
	if len(hosts) == 0 {
		t.Errorf("Expected %d hosts, got %d", len(testHosts), len(hosts))
	}
	if hosts["host1"].User != testHosts["host1"].User {
		t.Errorf("Expected %s, got %s", testHosts["host1"].User, hosts["host1"].User)
	}
}

func TestGetHostDetails(t *testing.T) {
	// test Hosts.GetHostDetails method
	testHosts := map[string]*backend.Host{
		"host1": {Alias: "host1", Host: "host1.com", User: "user1", Groups: []string{"group1", "group2"}},
	}
	cfg := &backend.Storage{
		Hosts: testHosts}
	h := Hosts{Prefix: "", Config: cfg}
	// mock http request
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/host1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "host1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	h.GetHostDetails(w, r)
	// check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	bs := w.Body.String()
	log.Printf("body %s", bs)
	var host backend.Host
	json.NewDecoder(w.Body).Decode(&host)
	if host.User != testHosts["host1"].User {
		t.Errorf("Expected %s, got %s", testHosts["host1"].User, host.User)
	}
}
