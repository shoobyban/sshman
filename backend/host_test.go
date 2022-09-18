package backend

import (
	"fmt"
	"testing"
	"time"
)

func TestGetUsers(t *testing.T) {
	sftp := &SFTPConn{mock: true}
	users := []*User{
		{Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot", Groups: []string{"a"}},
	}
	cfg := testConfig("foo", map[string]*Host{
		"a": {Alias: "a", Host: "a:22", User: "aroot", Groups: []string{"a", "b"}, Users: users},
		"b": {Alias: "b", Host: "b:22", User: "aroot", Groups: []string{"a", "c"}},
	}, users, sftp)
	h := cfg.GetHost("a")
	if len(h.GetUsers()) != 1 {
		t.Errorf("host not returning users")
	}
	if !h.HasUser("foo@email") {
		t.Errorf("HasUser doesn't work")
	}
	h = cfg.GetHost("b")
	if len(h.GetUsers()) != 0 {
		t.Errorf("host returning users when")
	}
	if h.HasUser("foo@email") {
		t.Errorf("HasUser returns true when no user")
	}
}

func TestUpdateHostGroups(t *testing.T) {
	sftp := &SFTPConn{mock: true}
	cfg := testConfig("foo", map[string]*Host{
		"hosta": {Alias: "hosta", Host: "a:22", User: "aroot", Groups: []string{"groupa", "groupb"}, LastUpdated: time.Now()},
		"hostb": {Alias: "hostb", Host: "b:22", User: "aroot", Groups: []string{"groupa", "groupc"}, LastUpdated: time.Now()},
	}, []*User{
		{Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot", Groups: []string{"groupa", "groupb"}},
		{Email: "bar@email", KeyType: "ssh-rsa", Key: "keydata", Name: "broot", Groups: []string{"groupa", "groupc"}},
	}, sftp)
	h := cfg.GetHost("hosta")
	old := h.Groups
	h.SetGroups([]string{"groupb"})
	h.UpdateGroups(cfg, old)
	g := cfg.GetGroups()
	if len(g) != 3 {
		t.Errorf("GetGroups did not work, groups: %#v", g)
	}
	old = h.Groups
	h.SetGroups([]string{"groupa", "groupb"})
	h.UpdateGroups(cfg, old)
	_, u := cfg.GetUserByEmail("foo@email")
	if u == nil {
		t.Errorf("error finding user by email")
		return
	}
	sftp.SetError(true)
	old = h.Groups
	h.SetGroups([]string{})
	if h.UpdateGroups(cfg, old) {
		t.Errorf("could update host with read error")
	}
}

func TestHostMoveGroup(t *testing.T) {
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
	fmt.Println("A")
	h := cfg.GetHost("hosta")
	old := h.Groups
	h.SetGroups([]string{"groupc"})
	h.UpdateGroups(cfg, old)
	if !h.HasUser("bar@email") {
		t.Errorf("host a doesn't have bar user")
	}
	if h.HasUser("foo@email") {
		t.Errorf("host a have foo user")
	}
}
