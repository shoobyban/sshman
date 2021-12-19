package backend

import (
	"fmt"
	"testing"
)

func TestGetUsers(t *testing.T) {
	sftp := &SFTPConn{mock: true}
	cfg := testConfig("foo", map[string]*Hostentry{
		"a": {Alias: "a", Host: "a:22", User: "aroot", Groups: []string{"a", "b"}, Users: []string{"foo@email"}},
		"b": {Alias: "b", Host: "b:22", User: "aroot", Groups: []string{"a", "c"}},
	}, map[string]*User{
		"asdfasdf": {Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot", Groups: []string{"a"}},
	}, sftp)
	h := cfg.Hosts["a"]
	if len(h.GetUsers()) != 1 {
		t.Errorf("host not returning users")
	}
	if !h.HasUser("foo@email") {
		t.Errorf("HasUser doesn't work")
	}
	h = cfg.Hosts["b"]
	if len(h.GetUsers()) != 0 {
		t.Errorf("host returning users when")
	}
	if h.HasUser("foo@email") {
		t.Errorf("HasUser returns true when no user")
	}
}

func TestUpdateHostGroups(t *testing.T) {
	sftp := &SFTPConn{mock: true}
	cfg := testConfig("foo", map[string]*Hostentry{
		"hosta": {Alias: "hosta", Host: "a:22", User: "aroot", Groups: []string{"groupa", "groupb"}},
		"hostb": {Alias: "hostb", Host: "b:22", User: "aroot", Groups: []string{"groupa", "groupc"}},
	}, map[string]*User{
		"asdfasdf0": {Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot", Groups: []string{"groupa", "groupb"}},
		"asdfasdf1": {Email: "bar@email", KeyType: "ssh-rsa", Key: "keydata", Name: "broot", Groups: []string{"groupa", "groupc"}},
	}, sftp)
	h := cfg.Hosts["hosta"]
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
		t.Errorf("could update server with read error")
	}
}

func TestServerMoveGroup(t *testing.T) {
	cfg := testConfig("foo", map[string]*Hostentry{
		"hosta": {Alias: "hosta", Host: "a:22", User: "aroot", Groups: []string{"groupb"}, Users: []string{"foo@email"}},
		"hostb": {Alias: "hostb", Host: "b:22", User: "aroot", Groups: []string{"groupc"}, Users: []string{"bar@email"}},
	}, map[string]*User{
		"C-7Hteo_D9vJXQ3UfzxbwnXaijM=": {Email: "foo@email", KeyType: "ssh-rsa", Key: "foo", Name: "aroot", Groups: []string{"groupa", "groupb"}},
		"djZ11qHY0KOijeymK7aKvYuvhvM=": {Email: "bar@email", KeyType: "ssh-rsa", Key: "bar1", Name: "broot", Groups: []string{"groupa", "groupc"}},
		"AzxIRrUGpKSOMs31RRXJHTSZrbM=": {Email: "bar2@email", KeyType: "ssh-rsa", Key: "bar2", Name: "buser", Groups: []string{"groupa"}},
	}, &SFTPConn{mock: true, testServers: map[string]SFTPMockServer{
		"a:22": {Host: "a:22", User: "test", File: "ssh-rsa foo foo\nssh-rsa bar2 bar2\n"},
		"b:22": {Host: "b:22", User: "test", File: "ssh-rsa bar1 bar\nssh-rsa bar2 bar2\n"},
	}})
	fmt.Println("A")
	h := cfg.Hosts["hosta"]
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
