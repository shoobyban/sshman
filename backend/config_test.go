package backend

import (
	"strings"
	"testing"
)

func testConfig(key string, hosts map[string]*Host, users map[string]*User, conn SFTP) *Storage {
	c := &Storage{
		Key:   key,
		Hosts: hosts,
		Users: users,
		conn:  conn,
	}
	for a, h := range c.Hosts {
		h.Config = c
		c.Hosts[a] = h
	}
	return c
}

func TestGetUserByEmail(t *testing.T) {
	cfg := testConfig("foo", map[string]*Host{}, map[string]*User{
		"a": {Email: "foo@email", Name: "foo"},
	}, &SFTPConn{mock: true})
	_, u := cfg.GetUserByEmail("foo@email")
	if u == nil || u.Name != "foo" {
		t.Errorf("GetUserByEmail doesn't work, %v", u)
	}
	_, u = cfg.GetUserByEmail("bar@email")
	if u != nil {
		t.Errorf("GetUserByEmail finds weird users, %v", u)
	}
}

func TestAddUser(t *testing.T) {
	sftp := &SFTPConn{mock: true, testHosts: map[string]SFTPMockHost{
		"a:22": {Host: "a:22", User: "test", File: "ssh-rsa foo rootuser\nssh-rsa bar1 user-a.com\n"},
		"b:22": {Host: "b:22", User: "test", File: "ssh-rsa foo rootuser\nssh-rsa bar2 user-b.com\n"},
	}}
	cfg := testConfig("foo", map[string]*Host{
		"a": {Alias: "a", Host: "a:22", User: "aroot", Groups: []string{"a"}},
		"b": {Alias: "b", Host: "b:22", User: "aroot", Groups: []string{"a", "b"}},
	}, map[string]*User{
		"asdfasdf": {Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot", Groups: []string{"a"}},
	}, sftp)
	cfg.AddUserByEmail("foo@email")
	testHosts := sftp.GetHosts()
	if !strings.Contains(testHosts["a:22"].File, "ssh-rsa keydata aroot") {
		t.Errorf("User not upoaded to host a:22")
	}
	if !strings.Contains(testHosts["b:22"].File, "ssh-rsa keydata aroot") {
		t.Errorf("User not upoaded to host b:22")
	}
	if cfg.AddUserByEmail("bar@email") {
		t.Errorf("Adding unknown user worked")
	}
}

func TestDelUser(t *testing.T) {
	sftp := &SFTPConn{mock: true, testHosts: map[string]SFTPMockHost{
		"a:22": {Host: "a:22", User: "test", File: "ssh-rsa keydata rootuser\nssh-rsa bar1 user-a.com\n\n"},
		"b:22": {Host: "b:22", User: "test", File: "ssh-rsa keydata rootuser\nssh-rsa bar2 user-b.com\n"},
	}}
	cfg := testConfig("foo", map[string]*Host{
		"a": {Alias: "a", Host: "a:22", User: "aroot"},
		"b": {Alias: "b", Host: "b:22", User: "aroot"},
	}, map[string]*User{
		"C-7Hteo_D9vJXQ3UfzxbwnXaijM=": {Email: "root@foo", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot"},
	}, sftp)
	_, toDel := cfg.GetUserByEmail("root@foo")
	if toDel == nil {
		t.Errorf("No such user: %v", "root@foo")
	}
	cfg.DelUserFromHosts(toDel)
	testHosts := sftp.GetHosts()
	if strings.Contains(testHosts["a:22"].File, "ssh-rsa foo rootuser") {
		t.Errorf("User not deleted from host a:22")
	}
	if strings.Contains(testHosts["b:22"].File, "ssh-rsa foo rootuser") {
		t.Errorf("User not deleted from host b:22")
	}
	if !strings.Contains(testHosts["a:22"].File, "ssh-rsa bar1 user-a.com") {
		t.Errorf("Other user deleted from host a:22")
	}
	if !strings.Contains(testHosts["b:22"].File, "ssh-rsa bar2 user-b.com") {
		t.Errorf("Other user deleted from host b:22")
	}
}

func TestAddDeleteHost(t *testing.T) {
	cfg := testConfig("foo", map[string]*Host{
		"a": {Alias: "a", Host: "a:22", User: "aroot"},
		"b": {Alias: "b", Host: "b:22", User: "aroot"},
	}, map[string]*User{
		"asdfasdf": {Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot"},
	}, &SFTPConn{mock: true})
	_, err := cfg.AddHost("c", "c:22", "cuser", ".", "groupa")
	if err != nil {
		t.Errorf("Adding host did not work: %v", err)
	}
	_, err = cfg.AddHost("c", "c:22", "cuser", "nonexistent", "groupa")
	if err == nil {
		t.Errorf("Adding host did not try to read file")
	}
	if len(cfg.Hosts) != 3 {
		t.Errorf("Adding host did not work")
	}
	cfg.DeleteHost("c")
	if len(cfg.Hosts) != 2 {
		t.Errorf("Deleting host did not work")
	}
	if cfg.DeleteHost("c") {
		t.Errorf("Deleting host did work again")
	}
}

func TestAddDeleteUser(t *testing.T) {
	sftp := &SFTPConn{
		mock:     true,
		expected: "ssh-rsa foo rootuser\nssh-rsa bar1 user-a.com\nssh-rsa key user\n",
		testHosts: map[string]SFTPMockHost{
			"a:22": {Host: "a:22", User: "test", File: "ssh-rsa foo rootuser\nssh-rsa bar1 user-a.com\n"},
			"b:22": {Host: "b:22", User: "test", File: "ssh-rsa foo rootuser\nssh-rsa bar2 user-b.com\n"},
		},
		testError: false,
	}
	cfg := testConfig("foo", map[string]*Host{
		"a": {Alias: "a", Host: "a:22", User: "aroot", Groups: []string{"groupa"}, Users: []string{"foo@email"}},
		"b": {Alias: "b", Host: "b:22", User: "aroot"},
	}, map[string]*User{
		"asdfasdf":                     {Email: "foo@email", KeyType: "ssh-rsa", Key: "foo", Name: "rootuser", Groups: []string{"groupa"}},
		"djZ11qHY0KOijeymK7aKvYuvhvM=": {Email: "user-a.com@a", KeyType: "ssh-rsa", Key: "bar1", Name: "user-a.com", Groups: []string{"groupb"}},
		"GM0J1PU4m76_UN8SIJ3jrmPePq8=": {Email: "rootuser@b", KeyType: "ssh-rsa", Key: "foo2", Name: "aroot", Groups: []string{"groupc"}},
	}, sftp)
	u, err := cfg.PrepareUser("bar@email", "test/dummy.key", "groupa", "groupb")
	if err != nil {
		t.Errorf("error while preparing user: %v %v", err, cfg.Users)
	}
	u.UpdateGroups(cfg, []string{})
	err = cfg.AddUser(u)
	if err != nil {
		t.Errorf("error while registering user: %v %v", err, cfg.Users)
	}
	if len(cfg.Users) != 5 {
		t.Errorf("Adding user did not work %#v", len(cfg.Users))
	}
	cfg.Hosts["a"] = &Host{Config: cfg, Alias: "a", Host: "a:22", User: "aroot", Groups: []string{"groupa"}, Users: []string{"foo@email", "bar@email"}}
	u, err = cfg.PrepareUser("bar@email", "test/dummy.key", "groupa", "groupb")
	if err != nil {
		t.Errorf("error while preparing user: %v %v", err, cfg.Users)
	}
	u.UpdateGroups(cfg, []string{})
	err = cfg.AddUser(u)
	if err != nil {
		t.Errorf("error while registering user: %v %v", err, cfg.Users)
	}
	g := cfg.GetGroups()
	if grp, ok := g["groupa"]; ok {
		if !contains(grp.Users, "bar@email") {
			t.Errorf("Group users doesn't have bar@email: %v", grp.Users)
		}
	} else {
		t.Errorf("Group groupa doesn't exits: %v", g)
	}
	// set expected file content for removal
	sftp.testHosts["a:22"] = SFTPMockHost{Host: "b:22", User: "test", File: "ssh-rsa foo rootuser\nssh-rsa bar1 user-a.com\n"}
	key, u := cfg.GetUserByEmail("bar@email")
	if key == "" {
		t.Errorf("User not found")
	}
	u, err = cfg.PrepareUser("bar@email", "test/dummy.key")
	if err != nil {
		t.Errorf("error while registering user: %v %v", err, cfg.Users)
	}
	if len(cfg.Users) != 5 {
		t.Errorf("Adding user did not work %#v", len(cfg.Users))
	}
	u.UpdateGroups(cfg, []string{})
	err = cfg.AddUser(u)
	if err != nil {
		t.Errorf("error while registering user: %v %v", err, cfg.Users)
	}
	g = cfg.GetGroups()
	if grp, ok := g["groupa"]; ok {
		if contains(grp.Users, "bar@email") {
			t.Errorf("Group users still have bar@email: %v", grp.Users)
		}
	} else {
		t.Errorf("Group groupa doesn't exits: %v", g)
	}
	cfg.DeleteUser("bar@email")
	if len(cfg.Users) != 4 {
		t.Errorf("Deleting user did not work %v", len(cfg.Users))
	}
	if cfg.DeleteUser("bar@email") {
		t.Errorf("Deleting user did work again")
	}
}

func TestBrokenKey(t *testing.T) {
	sftp := &SFTPConn{mock: true, testHosts: map[string]SFTPMockHost{
		"a:22": {Host: "a:22", User: "test", File: "ssh-rsa foo rootuser\nssh-rsa bar1 user-a.com\n"},
		"b:22": {Host: "b:22", User: "test", File: "ssh-rsa foo rootuser\nssh-rsa bar2 user-b.com\n"},
	}}
	cfg := testConfig("foo", map[string]*Host{
		"a": {Alias: "a", Host: "a:22", User: "aroot"},
		"b": {Alias: "b", Host: "b:22", User: "aroot"},
	}, map[string]*User{
		"asdfasdf": {Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot"},
	}, sftp)
	_, err := cfg.PrepareUser("bar@email", "test/broken.key")
	if err == nil {
		t.Errorf("could prepare with broken key info")
	}
}

func TestReadError(t *testing.T) {
	sftp := &SFTPConn{mock: true, testHosts: map[string]SFTPMockHost{
		"a:22": {Host: "a:22", User: "test", File: "ssh-rsa foo rootuser\nssh-rsa bar1 user-a.com\n"},
		"b:22": {Host: "b:22", User: "test", File: "ssh-rsa foo rootuser\nssh-rsa bar2 user-b.com\n"},
	}}
	cfg := testConfig("foo", map[string]*Host{
		"a": {Alias: "a", Host: "a:22", User: "aroot", Groups: []string{"groupa"}},
		"b": {Alias: "b", Host: "b:22", User: "aroot"},
	}, map[string]*User{
		"asdfasdf": {Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot"},
	}, sftp)
	sftp.SetError(true)
	_, err := cfg.PrepareUser("bar@email", "test/nonexistent.key", "groupa")
	if err == nil {
		t.Errorf("could prepare with read error %v", err)
	}
}

func TestConfigUpdate(t *testing.T) {
	sftp := &SFTPConn{mock: true, testHosts: map[string]SFTPMockHost{
		"a:22": {Host: "a:22", User: "test", File: "ssh-rsa foo rootuser\nssh-rsa bar1 user-a.com\n"},
		"b:22": {Host: "b:22", User: "test", File: "ssh-rsa foo rootuser\nssh-rsa bar2 user-b.com\n"},
	}}
	cfg := testConfig("foo", map[string]*Host{
		"a": {Alias: "a", Host: "a:22", User: "aroot"},
		"b": {Alias: "b", Host: "b:22", User: "aroot"},
	}, map[string]*User{
		"asdfasdf": {Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot"},
	}, sftp)
	cfg.Update("a", "b")
	if len(cfg.Hosts["a"].Users) != 2 {
		t.Errorf("Update() did not work, host a users: %v", cfg.Hosts["a"].Users)
	}
	if len(cfg.Hosts["b"].Users) != 2 {
		t.Errorf("Update() did not work, host b users: %v %v", cfg.Hosts["b"].Users, cfg.Hosts["a"].Users)
	}
}

func TestGetGroups(t *testing.T) {
	cfg := testConfig("foo", map[string]*Host{
		"a": {Alias: "a", Host: "a:22", User: "aroot", Groups: []string{"a", "b"}},
		"b": {Alias: "b", Host: "b:22", User: "aroot", Groups: []string{"a", "c"}},
	}, map[string]*User{
		"asdfasdf": {Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot", Groups: []string{"a"}},
	}, &SFTPConn{mock: true})
	g := cfg.GetGroups()
	if len(g) != 3 {
		t.Errorf("GetGroups did not work, groups: %v", g)
	}
	if grp, ok := g["a"]; ok {
		if len(grp.Hosts) != 2 {
			t.Errorf("Group a hosts incorrect: %v", grp.Hosts)
		}
		if len(grp.Users) != 1 {
			t.Errorf("Group a users incorrect: %v", grp.Users)
		}
	}
}

func TestGetHosts(t *testing.T) {
	cfg := testConfig("foo", map[string]*Host{
		"a": {Alias: "a", Host: "a:22", User: "aroot", Groups: []string{"a", "b"}},
		"b": {Alias: "b", Host: "b:22", User: "aroot", Groups: []string{"a", "c"}},
	}, map[string]*User{
		"asdfasdf": {Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot", Groups: []string{"a"}},
	}, &SFTPConn{mock: true})
	hosts := cfg.getHosts("a")
	if len(hosts) != 2 {
		t.Errorf("gethosts error: %v", hosts)
	}
	hosts = cfg.getHosts("b")
	if len(hosts) != 1 {
		t.Errorf("gethosts error: %v", hosts)
	}
	hosts = cfg.getHosts("c")
	if len(hosts) != 1 {
		t.Errorf("gethosts error: %v", hosts)
	}
}
