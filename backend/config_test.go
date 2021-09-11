package backend

import (
	"strings"
	"testing"
)

func testConfig(key string, hosts map[string]Hostentry, users map[string]User, conn SFTP) *config {
	c := &config{
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
	cfg := testConfig("foo", map[string]Hostentry{}, map[string]User{
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
	sftp := &SFTPConn{mock: true, testServers: map[string]SFTPMockServer{
		"a:22": {Host: "a:22", User: "test", File: "ssh-rsa foo rootuser\nssh-rsa bar1 user-a.com\n"},
		"b:22": {Host: "b:22", User: "test", File: "ssh-rsa foo rootuser\nssh-rsa bar2 user-b.com\n"},
	}}
	cfg := testConfig("foo", map[string]Hostentry{
		"a": {Alias: "a", Host: "a:22", User: "aroot", Groups: []string{"a"}},
		"b": {Alias: "b", Host: "b:22", User: "aroot", Groups: []string{"a", "b"}},
	}, map[string]User{
		"asdfasdf": {Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot", Groups: []string{"a"}},
	}, sftp)
	cfg.AddUserByEmail("foo@email")
	testServers := sftp.GetServers()
	if !strings.Contains(testServers["a:22"].File, "ssh-rsa keydata aroot") {
		t.Errorf("User not upoaded to server a:22")
	}
	if !strings.Contains(testServers["b:22"].File, "ssh-rsa keydata aroot") {
		t.Errorf("User not upoaded to server b:22")
	}
	if cfg.AddUserByEmail("bar@email") {
		t.Errorf("Adding unknown user worked")
	}
}

func TestDelUser(t *testing.T) {
	sftp := &SFTPConn{mock: true, testServers: map[string]SFTPMockServer{
		"a:22": {Host: "a:22", User: "test", File: "ssh-rsa foo rootuser\nssh-rsa bar1 user-a.com\n\n"},
		"b:22": {Host: "b:22", User: "test", File: "ssh-rsa foo rootuser\nssh-rsa bar2 user-b.com\n"},
	}}
	cfg := testConfig("foo", map[string]Hostentry{
		"a": {Alias: "a", Host: "a:22", User: "aroot"},
		"b": {Alias: "b", Host: "b:22", User: "aroot"},
	}, map[string]User{
		"asdfasdf": {Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot"},
	}, sftp)
	cfg.DelUserFromHosts(&User{Email: "root@foo", KeyType: "ssh-rsa", Key: "foo", Name: "rootuser"})
	testServers := sftp.GetServers()
	if strings.Contains(testServers["a:22"].File, "ssh-rsa foo rootuser") {
		t.Errorf("User not deleted from server a:22")
	}
	if strings.Contains(testServers["b:22"].File, "ssh-rsa foo rootuser") {
		t.Errorf("User not deleted from server b:22")
	}
	if !strings.Contains(testServers["a:22"].File, "ssh-rsa bar1 user-a.com") {
		t.Errorf("Other user deleted from server a:22")
	}
	if !strings.Contains(testServers["b:22"].File, "ssh-rsa bar2 user-b.com") {
		t.Errorf("Other user deleted from server b:22")
	}
}

func TestRegisterUnregisterServer(t *testing.T) {
	cfg := testConfig("foo", map[string]Hostentry{
		"a": {Alias: "a", Host: "a:22", User: "aroot"},
		"b": {Alias: "b", Host: "b:22", User: "aroot"},
	}, map[string]User{
		"asdfasdf": {Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot"},
	}, &SFTPConn{mock: true})
	err := cfg.RegisterServer("c", "c:22", "cuser", ".", "groupa")
	if err != nil {
		t.Errorf("Registering server did not work: %v", err)
	}
	err = cfg.RegisterServer("c", "c:22", "cuser", "nonexistent", "groupa")
	if err == nil {
		t.Errorf("Registering server did not try to read file")
	}
	if len(cfg.Hosts) != 3 {
		t.Errorf("Registering server did not work")
	}
	cfg.UnregisterServer("c")
	if len(cfg.Hosts) != 2 {
		t.Errorf("Unregistering server did not work")
	}
	if cfg.UnregisterServer("c") {
		t.Errorf("Unregistering server did work again")
	}
}

func TestRegisterUnregisterUser(t *testing.T) {
	sftp := &SFTPConn{
		mock:     true,
		expected: "ssh-rsa foo rootuser\nssh-rsa bar1 user-a.com\nssh-rsa key user\n",
		testServers: map[string]SFTPMockServer{
			"a:22": {Host: "a:22", User: "test", File: "ssh-rsa foo rootuser\nssh-rsa bar1 user-a.com\n"},
			"b:22": {Host: "b:22", User: "test", File: "ssh-rsa foo rootuser\nssh-rsa bar2 user-b.com\n"},
		},
	}
	cfg := testConfig("foo", map[string]Hostentry{
		"a": {Alias: "a", Host: "a:22", User: "aroot", Groups: []string{"groupa"}, Users: []string{"foo@email"}},
		"b": {Alias: "b", Host: "b:22", User: "aroot"},
	}, map[string]User{
		"asdfasdf":                     {Email: "foo@email", KeyType: "ssh-rsa", Key: "foo", Name: "rootuser", Groups: []string{"groupa"}},
		"djZ11qHY0KOijeymK7aKvYuvhvM=": {Email: "user-a.com@a", KeyType: "ssh-rsa", Key: "bar1", Name: "user-a.com", Groups: []string{"groupb"}},
		"GM0J1PU4m76_UN8SIJ3jrmPePq8=": {Email: "rootuser@b", KeyType: "ssh-rsa", Key: "foo2", Name: "aroot", Groups: []string{"groupc"}},
	}, sftp)
	err := cfg.RegisterUser([]string{}, "bar@email", "test/dummy.key", "groupa", "groupb")
	if err != nil {
		t.Errorf("error while registering user: %v %v", err, cfg.Users)
	}
	if len(cfg.Users) != 5 {
		t.Errorf("Registering user did not work %#v", len(cfg.Users))
	}
	cfg.Hosts["a"] = Hostentry{Config: cfg, Alias: "a", Host: "a:22", User: "aroot", Groups: []string{"groupa"}, Users: []string{"foo@email", "bar@email"}}
	err = cfg.RegisterUser([]string{}, "bar@email", "test/dummy.key", "groupa", "groupb")
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
	sftp.testServers["a:22"] = SFTPMockServer{Host: "b:22", User: "test", File: "ssh-rsa foo rootuser\nssh-rsa bar1 user-a.com\n"}
	_, u := cfg.GetUserByEmail("bar@email")
	err = cfg.RegisterUser(u.Groups, "bar@email", "test/dummy.key")
	if err != nil {
		t.Errorf("error while registering user: %v %v", err, cfg.Users)
	}
	if len(cfg.Users) != 5 {
		t.Errorf("Registering user did not work %#v", len(cfg.Users))
	}
	g = cfg.GetGroups()
	if grp, ok := g["groupa"]; ok {
		if contains(grp.Users, "bar@email") {
			t.Errorf("Group users still have bar@email: %v", grp.Users)
		}
	} else {
		t.Errorf("Group groupa doesn't exits: %v", g)
	}
	cfg.UnregisterUser("bar@email")
	if len(cfg.Users) != 4 {
		t.Errorf("Unregistering user did not work %v", len(cfg.Users))
	}
	if cfg.UnregisterUser("bar@email") {
		t.Errorf("Unregistering user did work again")
	}
}

func TestBrokenKey(t *testing.T) {
	sftp := &SFTPConn{mock: true, testServers: map[string]SFTPMockServer{
		"a:22": {Host: "a:22", User: "test", File: "ssh-rsa foo rootuser\nssh-rsa bar1 user-a.com\n"},
		"b:22": {Host: "b:22", User: "test", File: "ssh-rsa foo rootuser\nssh-rsa bar2 user-b.com\n"},
	}}
	cfg := testConfig("foo", map[string]Hostentry{
		"a": {Alias: "a", Host: "a:22", User: "aroot"},
		"b": {Alias: "b", Host: "b:22", User: "aroot"},
	}, map[string]User{
		"asdfasdf": {Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot"},
	}, sftp)
	err := cfg.RegisterUser([]string{}, "bar@email", "test/broken.key")
	if err == nil {
		t.Errorf("could register with broken key info")
	}

}

func TestConfigUpdate(t *testing.T) {
	sftp := &SFTPConn{mock: true, testServers: map[string]SFTPMockServer{
		"a:22": {Host: "a:22", User: "test", File: "ssh-rsa foo rootuser\nssh-rsa bar1 user-a.com\n"},
		"b:22": {Host: "b:22", User: "test", File: "ssh-rsa foo rootuser\nssh-rsa bar2 user-b.com\n"},
	}}
	cfg := testConfig("foo", map[string]Hostentry{
		"a": {Alias: "a", Host: "a:22", User: "aroot"},
		"b": {Alias: "b", Host: "b:22", User: "aroot"},
	}, map[string]User{
		"asdfasdf": {Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot"},
	}, sftp)
	cfg.Update("a", "b")
	if len(cfg.Hosts["a"].Users) != 2 {
		t.Errorf("Update() did not work, host a users: %v", cfg.Hosts["a"].Users)
	}
	if len(cfg.Hosts["b"].Users) != 2 {
		t.Errorf("Update() did not work, host b users: %v", cfg.Hosts["b"].Users)
	}
}

func TestGetGroups(t *testing.T) {
	cfg := testConfig("foo", map[string]Hostentry{
		"a": {Alias: "a", Host: "a:22", User: "aroot", Groups: []string{"a", "b"}},
		"b": {Alias: "b", Host: "b:22", User: "aroot", Groups: []string{"a", "c"}},
	}, map[string]User{
		"asdfasdf": {Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot", Groups: []string{"a"}},
	}, &SFTPConn{mock: true})
	g := cfg.GetGroups()
	if len(g) != 3 {
		t.Errorf("GetGroups did not work, groups: %v", g)
	}
	if grp, ok := g["a"]; ok {
		if len(grp.Servers) != 2 {
			t.Errorf("Group a servers incorrect: %v", grp.Servers)
		}
		if len(grp.Users) != 1 {
			t.Errorf("Group a users incorrect: %v", grp.Users)
		}
	}
}

func TestGetServers(t *testing.T) {
	cfg := testConfig("foo", map[string]Hostentry{
		"a": {Alias: "a", Host: "a:22", User: "aroot", Groups: []string{"a", "b"}},
		"b": {Alias: "b", Host: "b:22", User: "aroot", Groups: []string{"a", "c"}},
	}, map[string]User{
		"asdfasdf": {Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot", Groups: []string{"a"}},
	}, &SFTPConn{mock: true})
	servers := cfg.getServers("a")
	if len(servers) != 2 {
		t.Errorf("getservers error: %v", servers)
	}
	servers = cfg.getServers("b")
	if len(servers) != 1 {
		t.Errorf("getservers error: %v", servers)
	}
	servers = cfg.getServers("c")
	if len(servers) != 1 {
		t.Errorf("getservers error: %v", servers)
	}
}
