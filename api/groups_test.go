package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/shoobyban/sshman/backend"
)

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func TestUpdateGroupPropagatesUserToHost(t *testing.T) {
	conn := backend.MockConn(map[string]backend.SFTPMockHost{
		"host1.com": {Host: "host1.com", User: "root", File: "ssh-rsa existing existing\n"},
	})
	cfg := backend.NewData(&backend.MemoryStorage{})
	cfg.Conn = conn
	if err := cfg.AddHost(&backend.Host{Alias: "host1", Host: "host1.com", User: "root", Groups: []string{}, Config: cfg}, false); err != nil {
		t.Fatalf("failed to add host: %v", err)
	}
	if err := cfg.AddUser(&backend.User{Email: "user1", KeyType: "ssh-rsa", Key: "key1", Name: "User One", Config: cfg}, ""); err != nil {
		t.Fatalf("failed to add user: %v", err)
	}

	gh := GroupsHandler{Prefix: ""}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/group1", strings.NewReader(`{"Hosts":["host1"],"Users":["user1"],"Label":"group1"}`))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "group1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	r = r.WithContext(context.WithValue(r.Context(), ConfigKey, cfg))
	gh.UpdateGroup(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d %s", http.StatusOK, w.Code, w.Body.String())
	}
	if !cfg.GetHost("host1").HasUser("user1") {
		t.Fatalf("expected propagated user on host after group update")
	}
	if !strings.Contains(conn.GetHosts()["host1.com"].File, "ssh-rsa key1 User One") {
		t.Fatalf("expected authorized_keys to include propagated user, got %q", conn.GetHosts()["host1.com"].File)
	}
}

func TestCreateGroup(t *testing.T) {
	cfg := backend.NewData(&backend.MemoryStorage{})
	cfg.AddHost(
		&backend.Host{Alias: "host1", Host: "host1.com", User: "user1", Groups: []string{"group2"}},
		false)
	gh := GroupsHandler{Prefix: ""}
	// mock http request
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/group1", strings.NewReader(`{"Hosts":["host1"], "Label":"group1"}`))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "group1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	r = r.WithContext(context.WithValue(r.Context(), ConfigKey, cfg))
	gh.CreateGroup(w, r)
	// check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	h := cfg.GetHost("host1")
	if !contains(h.Groups, "group1") {
		t.Errorf("Expected host to be in group1, got %v", h.Groups)
	}
}
func TestUpdateGroup(t *testing.T) {

	cfg := backend.NewData(&backend.MemoryStorage{})
	cfg.AddHost(
		&backend.Host{Alias: "host1", Host: "host1.com", User: "user1", Groups: []string{"group2"}},
		false)
	cfg.AddGroup("group1", []string{}, []string{"user1"})
	gh := GroupsHandler{Prefix: ""}
	// mock http request
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/group1", strings.NewReader(`{"Hosts":["host1"], "Label":"group1"}`))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "group1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	r = r.WithContext(context.WithValue(r.Context(), ConfigKey, cfg))
	gh.UpdateGroup(w, r)
	// check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	h := cfg.GetHost("host1")
	if !contains(h.Groups, "group1") {
		t.Errorf("Expected host to be in group1, got %v", h.Groups)
	}
	g := cfg.GetGroup("group1")
	if !g.HasHost("host1") {
		t.Errorf("Expected group to contain host1, got %v", g.Hosts)
	}
}

func TestDeleteGroup(t *testing.T) {

	cfg := backend.NewData(&backend.MemoryStorage{})
	cfg.AddHost(
		&backend.Host{Alias: "host1", Host: "host1.com", User: "user1", Groups: []string{"group2"}},
		false)
	cfg.AddGroup("group1", []string{"host1"}, []string{"user1"})
	gh := GroupsHandler{Prefix: ""}
	// mock http request
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/group1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "group1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	r = r.WithContext(context.WithValue(r.Context(), ConfigKey, cfg))
	gh.DeleteGroup(w, r)
	// check response
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, w.Code)
	}
}
