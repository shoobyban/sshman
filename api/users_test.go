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
	testUsers := map[string]*backend.User{
		"user1": {Email: "sam@test.com", KeyType: "dummy", Key: "key1", Groups: []string{"group1", "group2"}},
	}
	cfg := &backend.Storage{Users: testUsers}
	u := Users{Prefix: "users", Config: cfg}
	// mock http request
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users", nil)
	u.GetAllUsers(w, r)
	// check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	var users []backend.User
	json.NewDecoder(w.Body).Decode(&users)
	if len(users) == 0 {
		t.Errorf("Expected %d users, got %d", len(testUsers), len(users))
		return
	}
	if users[0].Email != testUsers["user1"].Email {
		t.Errorf("Expected %s, got %s", testUsers["user1"].Email, users[0].Email)
	}
}

func TestGetUserDetails(t *testing.T) {
	// test Users.GetUserDetails method
	testUsers := map[string]*backend.User{
		"user1": {Email: "sam@test1.com", KeyType: "dummy", Key: "key1", Groups: []string{"group1", "group2"}},
	}
	cfg := &backend.Storage{Users: testUsers}
	u := Users{Prefix: "users", Config: cfg}
	// mock http request
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users/sam@test1.com", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("email", "sam@test1.com")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	u.GetUserDetails(w, r)
	// check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	var user backend.User
	json.NewDecoder(w.Body).Decode(&user)
	if user.Email != testUsers["user1"].Email {
		t.Errorf("Expected %s, got %s", testUsers["user1"].Email, user.Email)
	}
}

func TestCreateUser(t *testing.T) {
	// test Users.CreateUser method
	testUsers := map[string]*backend.User{}
	cfg := &backend.Storage{Users: testUsers, Conn: &backend.SFTPConn{}}
	u := Users{Prefix: "users", Config: cfg}
	// mock http request
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users",
		strings.NewReader(`["sam@test1.com", "../backend/fixtures/dummy.key", "group1", "group2"]`),
	)
	u.CreateUser(w, r)
	// check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	var user backend.User
	json.NewDecoder(w.Body).Decode(&user)
	if user.Email != "sam@test1.com" {
		t.Errorf("Expected sam@test1.com, got %s", user.Email)
	}
}

func TestUpdateUser(t *testing.T) {
	// test Users.UpdateUser method
	testUsers := map[string]*backend.User{
		"user1": {Email: "sam@test1.com", KeyType: "dummy", Key: "key1", Groups: []string{"group1", "group2"}},
	}
	cfg := &backend.Storage{Users: testUsers, Conn: &backend.SFTPConn{}}
	u := Users{Prefix: "users", Config: cfg}
	// mock http request
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/users/sam@test1.com",
		strings.NewReader(`["sam@test1.com", "../backend/fixtures/dummy.key", "group1", "group2"]`),
	)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("email", "user1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	u.UpdateUser(w, r)
	// check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	var user backend.User
	json.NewDecoder(w.Body).Decode(&user)
	if user.Email != "sam@test1.com" {
		t.Errorf("Expected sam@test1.com got %s", user.Email)
	}
}

func TestDeleteUser(t *testing.T) {
	// test Users.DeleteUser method
	testUsers := map[string]*backend.User{
		"user1": {Email: "sam@test1.com", KeyType: "dummy", Key: "key1", Groups: []string{"group1", "group2"}},
	}
	cfg := &backend.Storage{Users: testUsers, Conn: &backend.SFTPConn{}}
	u := Users{Prefix: "users", Config: cfg}
	// mock http request
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/users/sam@test1.com", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("email", "sam@test1.com")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	u.DeleteUser(w, r)
	// check response
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, w.Code)
	}
	if len(cfg.Users) != 0 {
		t.Errorf("Expected %d users, got %d", 0, len(cfg.Users))
	}
}
