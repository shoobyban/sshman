package backend

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMoveUserToGroup(t *testing.T) {
	users := []*User{
		{Email: "foo@email", KeyType: "ssh-rsa", Key: "foo", Name: "aroot", Groups: []string{"groupa", "groupb"}},
		{Email: "bar@email", KeyType: "ssh-rsa", Key: "bar1", Name: "broot", Groups: []string{"groupa", "groupc"}},
		{Email: "bar2@email", KeyType: "ssh-rsa", Key: "bar2", Name: "buser", Groups: []string{"groupa"}},
	}
	cfg := testConfig("foo", map[string]*Host{
		"hosta": {Alias: "hosta", Host: "a:22", User: "aroot", Groups: []string{"groupb"}, Users: []*User{users[0]}, LastUpdated: time.Now()},
		"hostb": {Alias: "hostb", Host: "b:22", User: "aroot", Groups: []string{"groupc"}, Users: []*User{users[1]}, LastUpdated: time.Now()},
	}, users, &SFTPConn{mock: true, testHosts: map[string]SFTPMockHost{
		"a:22": {Host: "a:22", User: "test", File: "ssh-rsa foo foo\nssh-rsa bar2 bar2\n"},
		"b:22": {Host: "b:22", User: "test", File: "ssh-rsa bar1 bar\nssh-rsa bar2 bar2\n"},
	}})
	_, u := cfg.GetUserByEmail("bar@email")
	if u == nil {
		t.Errorf("error finding user by email")
		return
	}
	old := u.Groups
	u.SetGroups([]string{"groupb"})
	u.UpdateGroups(cfg, old)
	hosta := cfg.GetHost("hosta")
	if hosta == nil {
		t.Errorf("error finding host a by alias")
	}
	hostb := cfg.GetHost("hostb")
	if hostb == nil {
		t.Errorf("error finding host b by alias")
	}
	if !hosta.HasUser(u.Email) {
		t.Errorf("user is not on hosta")
	}
	if hostb.HasUser(u.Email) {
		t.Errorf("user is still on hostb")
	}
}

func TestNewUser(t *testing.T) {
	email := "test@example.com"
	keytype := "rsa"
	key := "ssh-rsa AAAAB3..."
	name := "Test User"

	user := NewUser(email, keytype, key, name)

	assert.Equal(t, email, user.Email)
	assert.Equal(t, keytype, user.KeyType)
	assert.Equal(t, key, user.Key)
	assert.Equal(t, name, user.Name)
}

func TestUpdateGroups(t *testing.T) {
	// Use testConfig helper to create a minimal config
	sftp := &SFTPConn{mock: true}
	cfg := testConfig("foo", map[string]*Host{}, []*User{}, sftp)
	user := NewUser("test@example.com", "rsa", "ssh-rsa AAAAB3...", "Test User")
	user.Config = cfg

	oldGroups := []string{"group1", "group2"}
	newGroups := []string{"group2", "group3"}
	user.SetGroups(newGroups)

	err := user.UpdateGroups(cfg, oldGroups)
	assert.NoError(t, err)
}

func TestGetAndSetGroups(t *testing.T) {
	user := NewUser("test@example.com", "rsa", "ssh-rsa AAAAB3...", "Test User")
	groups := []string{"group1", "group2"}

	user.SetGroups(groups)
	assert.Equal(t, groups, user.GetGroups())
}
