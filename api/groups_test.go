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

func TestCreateGroup(t *testing.T) {
	cfg := backend.NewTestStorage()
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

	cfg := backend.NewTestStorage()
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

	cfg := backend.NewTestStorage()
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
