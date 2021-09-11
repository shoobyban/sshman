package backend

import "testing"

func TestUpdateUserGroups(t *testing.T) {
	cfg := testConfig("foo", map[string]Hostentry{
		"hosta": {Alias: "hosta", Host: "a:22", User: "aroot", Groups: []string{"groupa", "groupb"}},
		"hostb": {Alias: "hostb", Host: "b:22", User: "aroot", Groups: []string{"groupa", "groupc"}},
	}, map[string]User{
		"asdfasdf0": {Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot", Groups: []string{"groupa", "groupb"}},
		"asdfasdf1": {Email: "bar@email", KeyType: "ssh-rsa", Key: "keydata", Name: "broot", Groups: []string{"groupa", "groupc"}},
	}, &SFTPConn{mock: true})
	_, u := cfg.GetUserByEmail("foo@email")
	if u == nil {
		t.Errorf("error finding user by email")
		return
	}
	old := u.Groups
	u.SetGroups([]string{"groupa"})
	hosta := cfg.Hosts["hosta"]
	hostb := cfg.Hosts["hostb"]
	if hosta.HasUser("foo@email") {
		t.Errorf("user is not on host")
	}
	if hostb.HasUser("foo@email") {
		t.Errorf("user is still on host")
	}
	_, u = cfg.GetUserByEmail("bar@email")
	if u == nil {
		t.Errorf("error finding user by email")
		return
	}
	old = u.Groups
	u.SetGroups([]string{"groupb"})
	u.UpdateGroups(cfg, old)
	if hostb.HasUser("bar@email") {
		t.Errorf("user is on wrong host")
	}
	if hosta.HasUser("bar@email") {
		t.Errorf("user is still on host b")
	}
	old = u.Groups
	u.Groups = []string{"groupb"}
	u.UpdateGroups(cfg, old)
	if hostb.HasUser("bar@email") {
		t.Errorf("user is on wrong host")
	}
	if hosta.HasUser("bar@email") {
		t.Errorf("user is still on host b")
	}
}
