package backend

import "testing"

func TestMoveUserToGroup(t *testing.T) {
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
	_, u := cfg.GetUserByEmail("bar@email")
	if u == nil {
		t.Errorf("error finding user by email")
		return
	}
	old := u.Groups
	u.SetGroups([]string{"groupb"})
	u.UpdateGroups(cfg, old)
	hosta := cfg.Hosts["hosta"]
	hostb := cfg.Hosts["hostb"]
	if !hosta.HasUser(u.Email) {
		t.Errorf("user is not on hosta")
	}
	if hostb.HasUser(u.Email) {
		t.Errorf("user is still on hostb")
	}
}
