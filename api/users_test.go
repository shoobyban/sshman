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
	cfg := backend.NewData(&backend.MemoryStorage{})
	testUsers := map[string]*backend.User{
		"user1": {Email: "sam@test.com", KeyType: "dummy", Key: "key1", Groups: []string{"group1", "group2"}, Config: cfg},
	}
	cfg.AddUser(testUsers["user1"], "")
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
	var users []backend.User
	json.NewDecoder(w.Body).Decode(&users)
	if len(users) == 0 {
		t.Errorf("Expected %d users, got %d", len(testUsers), len(users))
		return
	}
}

func TestGetAllUsersNormalizesNilCollections(t *testing.T) {
	cfg := backend.NewData(&backend.MemoryStorage{})
	cfg.AddUser(&backend.User{
		Email:   "discovered@test.com",
		KeyType: "ssh-rsa",
		Key:     "key1",
		Name:    "Discovered User",
		Hosts:   []string{"host1"},
		Config:  cfg,
	}, "")

	u := UsersHandler{Prefix: "users"}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users", nil)
	r = r.WithContext(context.WithValue(r.Context(), ConfigKey, cfg))
	u.GetAllUsers(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d %s", http.StatusOK, w.Code, w.Body.String())
	}

	var users []backend.User
	if err := json.NewDecoder(w.Body).Decode(&users); err != nil {
		t.Fatalf("Failed to decode users response: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("Expected 1 user, got %d", len(users))
	}
	if users[0].Groups == nil {
		t.Fatalf("Expected groups to be normalized to an empty array")
	}
	if users[0].Roles == nil {
		t.Fatalf("Expected roles to be normalized to an empty array")
	}
}

func TestGetUserDetails(t *testing.T) {
	// test Users.GetUserDetails method
	cfg := backend.NewData(&backend.MemoryStorage{})
	testUsers := map[string]*backend.User{
		"u1": {Email: "sam@test1.com", KeyType: "dummy", Key: "key1", Groups: []string{"group1", "group2"}},
	}
	cfg.AddUser(testUsers["u1"], "")
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
	cfg := backend.NewData(&backend.MemoryStorage{})
	cfg.AddUser(testUsers["user1"], "")
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

func TestCreateUserRollsBackOnPropagationFailure(t *testing.T) {
	cfg := backend.NewData(&backend.MemoryStorage{})
	conn := backend.MockConn(map[string]backend.SFTPMockHost{
		"a:22": {Host: "a:22", User: "test", File: "ssh-rsa existing existing\n"},
	})
	conn.SetError(true)
	cfg.Conn = conn
	if err := cfg.AddHost(&backend.Host{Alias: "hosta", Host: "a:22", User: "root", Groups: []string{"group1"}}, false); err != nil {
		t.Fatalf("failed to add host: %v", err)
	}

	u := UsersHandler{Prefix: "users"}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users",
		strings.NewReader(`{"email":"sam@test1.com","name":"Sam Test","type":"ssh-rsa","key":"key1","groups":["group1"]}`),
	)
	r = r.WithContext(context.WithValue(r.Context(), ConfigKey, cfg))
	u.CreateUser(w, r)

	if w.Code != http.StatusConflict {
		t.Fatalf("Expected status code %d, got %d %s", http.StatusConflict, w.Code, w.Body.String())
	}
	if len(cfg.Users()) != 0 {
		t.Fatalf("Expected failed create to roll back stored user, got %d users", len(cfg.Users()))
	}

	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if resp.Error.Message == "" {
		t.Fatalf("expected structured JSON error message")
	}
}

func TestUpdateUser(t *testing.T) {
	// test Users.UpdateUser method
	cfg := backend.NewData(&backend.MemoryStorage{})
	testUsers := map[string]*backend.User{
		"user1": {Email: "sam@test1.com", KeyType: "dummy", Key: "key1", Groups: []string{"group1", "group2"}, Config: cfg},
	}
	cfg.AddUser(testUsers["user1"], "")
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

func TestUpdateUserRejectsEmptyBodyWithJSONError(t *testing.T) {
	cfg := backend.NewData(&backend.MemoryStorage{})
	testUser := &backend.User{Email: "sam@test1.com", KeyType: "dummy", Key: "key1", Name: "Sam Test", Config: cfg}
	if err := cfg.AddUser(testUser, ""); err != nil {
		t.Fatalf("failed to add user: %v", err)
	}

	u := UsersHandler{Prefix: "users"}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/users/sam@test1.com", strings.NewReader(""))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "sam@test1.com")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	r = r.WithContext(context.WithValue(r.Context(), ConfigKey, cfg))
	u.UpdateUser(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status code %d, got %d %s", http.StatusBadRequest, w.Code, w.Body.String())
	}

	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode JSON error response: %v", err)
	}
	if resp.Error.Message != "Invalid request body." {
		t.Fatalf("Expected JSON error message, got %#v", resp)
	}
}

func TestDeleteUser(t *testing.T) {
	// test Users.DeleteUser method
	cfg := backend.NewData(&backend.MemoryStorage{})
	testUsers := []backend.User{
		{Email: "sam@test1.com", KeyType: "dummy", Key: "key1", Groups: []string{"group1", "group2"}},
	}
	for _, user := range testUsers {
		cfg.AddUser(&user, "")
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
