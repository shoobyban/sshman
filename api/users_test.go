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

func TestGetAllUsers(t *testing.T) {
	// test Users.GetAllUsers method
	cfg := backend.NewTestStorage()
	testUsers := map[string]*backend.User{
		"user1": {Email: "sam@test.com", KeyType: "dummy", Key: "key1", Groups: []string{"group1", "group2"}, Config: cfg},
	}
	cfg.AddUser(testUsers["user1"])
	u := UsersHandler{Prefix: "users"}
	// mock http request
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users", nil)
	r = r.WithContext(context.WithValue(r.Context(), ConfigKey, cfg))
	u.GetAllUsers(w, r)
	// check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d %s", http.StatusOK, w.Code, w.Body.String())
		return
	}
	var users map[string]backend.User
	json.NewDecoder(w.Body).Decode(&users)
	if len(users) == 0 {
		t.Errorf("Expected %d users, got %d", len(testUsers), len(users))
		return
	}
}

func TestGetUserDetails(t *testing.T) {
	// test Users.GetUserDetails method
	cfg := backend.NewTestStorage()
	testUsers := map[string]*backend.User{
		"u1": {Email: "sam@test1.com", KeyType: "dummy", Key: "key1", Groups: []string{"group1", "group2"}},
	}
	cfg.AddUser(testUsers["u1"])
	u := UsersHandler{Prefix: "users"}
	// mock http request
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users/EHOrbNpLmRzSn56DowfzQASukyc=", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "EHOrbNpLmRzSn56DowfzQASukyc=")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	r = r.WithContext(context.WithValue(r.Context(), ConfigKey, cfg))
	u.GetUserDetails(w, r)
	// check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	var user backend.User
	json.NewDecoder(w.Body).Decode(&user)
	if user.Email != testUsers["u1"].Email {
		t.Errorf("Expected %s, got %s", testUsers["u1"].Email, user.Email)
	}
}

func TestCreateUser(t *testing.T) {
	// test Users.CreateUser method
	testUsers := map[string]*backend.User{}
	cfg := backend.NewTestStorage()
	cfg.AddUser(testUsers["user1"])
	u := UsersHandler{Prefix: "users"}
	// mock http request
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users",
		strings.NewReader(`{"email": "sam@test1.com","keyfile": "dummy key info", "groups": ["group1", "group2"]}`),
	)
	r = r.WithContext(context.WithValue(r.Context(), ConfigKey, cfg))
	u.CreateUser(w, r)
	// check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	var user backend.User
	json.NewDecoder(w.Body).Decode(&user)
	if user.Email != "sam@test1.com" {
		t.Errorf("Expected sam@test1.com, got %s", user.Email)
		return
	}
	if user.Key != "key" {
		t.Errorf("Expected dummy key, got %s", user.Key)
	}
}

func TestUpdateUser(t *testing.T) {
	// test Users.UpdateUser method
	cfg := backend.NewTestStorage()
	testUsers := map[string]*backend.User{
		"user1": {Email: "sam@test1.com", KeyType: "dummy", Key: "key1", Groups: []string{"group1", "group2"}, Config: cfg},
	}
	cfg.AddUser(testUsers["user1"])
	u := UsersHandler{Prefix: "users"}
	// mock http request
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/users/EHOrbNpLmRzSn56DowfzQASukyc=",
		strings.NewReader(`{"email": "sam@test1.com","keyfile": "dummy key info", "groups": ["group1", "group2"]}`),
	)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "EHOrbNpLmRzSn56DowfzQASukyc=")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	r = r.WithContext(context.WithValue(r.Context(), ConfigKey, cfg))
	u.UpdateUser(w, r)
	// check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d %s", http.StatusOK, w.Code, w.Body.String())
		return
	}
	var user backend.User
	json.NewDecoder(w.Body).Decode(&user)
	if user.Email != "sam@test1.com" {
		t.Errorf("Expected sam@test1.com got %s", user.Email)
		return
	}
	if user.Key != "key" {
		t.Errorf("Expected dummy key, got %s", user.Key)
	}
}

func TestDeleteUser(t *testing.T) {
	// test Users.DeleteUser method
	cfg := backend.NewTestStorage()
	testUsers := []backend.User{
		{Email: "sam@test1.com", KeyType: "dummy", Key: "key1", Groups: []string{"group1", "group2"}},
	}
	for _, user := range testUsers {
		cfg.AddUser(&user)
	}
	u := UsersHandler{Prefix: "users"}
	// mock http request
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/users/EHOrbNpLmRzSn56DowfzQASukyc=", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "EHOrbNpLmRzSn56DowfzQASukyc=")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	r = r.WithContext(context.WithValue(r.Context(), ConfigKey, cfg))
	u.DeleteUser(w, r)
	// check response
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, w.Code)
		return
	}
	if len(cfg.Users()) != 0 {
		t.Errorf("Expected %d users, got %d", 0, len(cfg.Users()))
	}
}
