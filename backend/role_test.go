package backend

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestDefaultRolesExist(t *testing.T) {
    cfg := DefaultConfig()
    roles := cfg.Roles()
    if _, ok := roles["admin"]; !ok {
        t.Fatalf("expected default role 'admin' to exist")
    }
    if _, ok := roles["user"]; !ok {
        t.Fatalf("expected default role 'user' to exist")
    }
}

func TestAssignRoleToUser(t *testing.T) {
    // create a minimal config and user
    cfg := DefaultConfig()
    // use in-memory storage for tests to avoid touching disk
    cfg.Storage = &MemoryStorage{}
    user := &User{Email: "roletest@example.com", KeyType: "ssh-rsa", Key: "keydata", Name: "Role Test"}
    // add user to config
    err := cfg.AddUser(user, "")
    if err != nil {
        t.Fatalf("failed to add user: %v", err)
    }

    // ensure role exists
    roles := cfg.Roles()
    if _, ok := roles["admin"]; !ok {
        t.Fatalf("expected role 'admin' to exist")
    }

    // assign role to user and update
    user.Roles = append(user.Roles, "admin")
    if err := cfg.UpdateUser(user); err != nil {
        t.Fatalf("UpdateUser failed: %v", err)
    }

    // retrieve user and assert role persisted
    _, u := cfg.GetUserByEmail("roletest@example.com")
    if u == nil {
        t.Fatalf("user not found after update")
    }
    assert.Contains(t, u.Roles, "admin")
}
