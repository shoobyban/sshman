package api

import (
	"context"
	"encoding/json"
	"log"
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
	testHosts := map[string]*backend.Host{
		"host1": {Alias: "host1", Host: "host1.com", User: "user1", Groups: []string{"group1", "group2"}},
	}
	cfg := &backend.Storage{
		Hosts: testHosts}
	h := Hosts{Prefix: "", Config: cfg}
	// mock http request
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
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
	r := httptest.NewRequest(http.MethodGet, "/host1", nil)
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

func TestUpdateHost(t *testing.T) {
	// test Hosts.UpdateHost method
	testHosts := map[string]*backend.Host{
		"host1": {Alias: "host1", Host: "host1.com", User: "user1", Groups: []string{"group1", "group2"}},
	}
	cfg := &backend.Storage{
		Hosts: testHosts,
		Conn:  &backend.SFTPConn{},
	}
	h := Hosts{Prefix: "", Config: cfg}
	// mock http request
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/host1",
		strings.NewReader(`["host1", "host1.com", "user2", "../backend/fixtures/dummy.key", "group1", "group2"]`),
	)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "host1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	h.UpdateHost(w, r)
	// check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	bs := w.Body.String()
	log.Printf("body %s", bs)
	var host backend.Host
	json.NewDecoder(w.Body).Decode(&host)
	if host.User != "user2" {
		t.Errorf("Expected user2, got %s %#v", host.User, host)
	}
}

// test Hosts.CreateHost method
func TestCreateHost(t *testing.T) {
	testHosts := map[string]*backend.Host{}
	cfg := &backend.Storage{
		Hosts: testHosts,
		Conn:  &backend.SFTPConn{},
	}
	h := Hosts{Prefix: "", Config: cfg}
	// mock http request
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/",
		strings.NewReader(`["host1", "host1.com", "user2", "../backend/fixtures/dummy.key", "group1", "group2"]`),
	)
	h.CreateHost(w, r)
	// check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	bs := w.Body.String()
	log.Printf("body %s", bs)
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
	cfg := &backend.Storage{
		Hosts: testHosts,
		Conn:  &backend.SFTPConn{},
	}
	h := Hosts{Prefix: "", Config: cfg}
	// mock http request
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/host1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "host1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	h.DeleteHost(w, r)
	// check response
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, w.Code)
	}
}
