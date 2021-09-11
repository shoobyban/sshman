package backend

import "testing"

func TestGetUsers(t *testing.T) {
	cfg := testConfig("foo", map[string]Hostentry{
		"a": {Alias: "a", Host: "a:22", User: "aroot", Groups: []string{"a", "b"}, Users: []string{"foo@email"}},
		"b": {Alias: "b", Host: "b:22", User: "aroot", Groups: []string{"a", "c"}},
	}, map[string]User{
		"asdfasdf": {Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot", Groups: []string{"a"}},
	}, &SFTPConn{mock: true})
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
	cfg := testConfig("foo", map[string]Hostentry{
		"hosta": {Alias: "hosta", Host: "a:22", User: "aroot", Groups: []string{"groupa", "groupb"}},
		"hostb": {Alias: "hostb", Host: "b:22", User: "aroot", Groups: []string{"groupa", "groupc"}},
	}, map[string]User{
		"asdfasdf0": {Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot", Groups: []string{"groupa", "groupb"}},
		"asdfasdf1": {Email: "bar@email", KeyType: "ssh-rsa", Key: "keydata", Name: "broot", Groups: []string{"groupa", "groupc"}},
	}, &SFTPConn{mock: true})
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
}
