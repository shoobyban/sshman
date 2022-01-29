package backend

import "testing"

func TestMoveUserToGroup(t *testing.T) {
	cfg := testConfig("foo", map[string]*Host{
		"hosta": {Alias: "hosta", Host: "a:22", User: "aroot", Groups: []string{"groupb"}, Users: []string{"foo@email"}},
		"hostb": {Alias: "hostb", Host: "b:22", User: "aroot", Groups: []string{"groupc"}, Users: []string{"bar@email"}},
	}, []*User{
		{Email: "foo@email", KeyType: "ssh-rsa", Key: "foo", Name: "aroot", Groups: []string{"groupa", "groupb"}},
		{Email: "bar@email", KeyType: "ssh-rsa", Key: "bar1", Name: "broot", Groups: []string{"groupa", "groupc"}},
		{Email: "bar2@email", KeyType: "ssh-rsa", Key: "bar2", Name: "buser", Groups: []string{"groupa"}},
	}, &SFTPConn{mock: true, testHosts: map[string]SFTPMockHost{
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
