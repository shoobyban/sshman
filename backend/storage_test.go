package backend

import (
	"strings"
	"testing"
)

func testConfig(key string, hosts map[string]*Host, users []*User, conn SFTP) *Storage {
	c := &Storage{
		hosts:  map[string]*Host{},
		users:  map[string]*User{},
		Groups: map[string]Group{},
		Conn:   conn,
		Log:    NewLog(false),
	}
	for a, h := range hosts {
		h.Config = c
		h.Alias = a
		c.SetHost(a, h)
	}
	for _, u := range users {
		c.Log.Infof("Adding user %p %v", u, u.Email)
		u.Config = c
		c.AddUser(u)
	}
	emails := []string{}
	for _, u := range c.Users() {
		emails = append(emails, u.Email)
	}
	c.Log.Infof("Test Users: %v", emails)
	c.updateGroups()
	return c
}

func TestGetUserByEmail(t *testing.T) {
	cfg := testConfig("foo", map[string]*Host{}, []*User{{Email: "foo@email", Name: "foo"}}, &SFTPConn{mock: true})
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
	}, []*User{
		{Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot", Groups: []string{"a"}},
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
	}, []*User{
		{Email: "root@foo", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot"},
	}, sftp)
	_, toDel := cfg.GetUserByEmail("root@foo")
	if toDel == nil {
		t.Errorf("No such user: %v", "root@foo")
	}
	err := cfg.DelUserFromHosts(toDel)
	if err != nil {
		t.Errorf("DelUserFromHosts failed: %v", err)
	}
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

func TestAddAndDeleteHost(t *testing.T) {
	cfg := testConfig("foo", map[string]*Host{
		"a": {Alias: "a", Host: "a:22", User: "aroot"},
		"b": {Alias: "b", Host: "b:22", User: "aroot"},
	}, []*User{
		{Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot"},
	}, &SFTPConn{mock: true})

	// test reading nonexisting host
	_, err := cfg.PrepareHost("c", "c:22", "cuser", "nonexistent", "groupa")
	if err == nil {
		t.Errorf("Preparing host did not try to read file")
		return
	}

	h, err := cfg.PrepareHost("c", "c:22", "cuser", ".", "groupa")
	if err != nil {
		t.Errorf("Adding host did not work: %v", err)
	}
	err = cfg.AddHost(h, false)
	if err != nil {
		t.Errorf("Adding host did not work: %v", err)
	}
	if len(cfg.Hosts()) != 3 {
		t.Errorf("Adding host did not work")
	}
	cfg.DeleteHost("c")
	if len(cfg.Hosts()) != 2 {
		t.Errorf("Deleting host did not work")
	}
	if cfg.DeleteHost("c") {
		t.Errorf("Deleting host did work again")
	}
}

func TestAddAndDeleteUser(t *testing.T) {
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
	}, []*User{
		{Email: "foo@email", KeyType: "ssh-rsa", Key: "foo", Name: "rootuser", Groups: []string{"groupa"}},
		{Email: "user-a.com@a", KeyType: "ssh-rsa", Key: "bar1", Name: "user-a.com", Groups: []string{"groupb"}},
		{Email: "rootuser@b", KeyType: "ssh-rsa", Key: "foo2", Name: "aroot", Groups: []string{"groupc"}},
	}, sftp)
	u, err := cfg.PrepareUser("bar@email", "fixtures/dummy.key", "groupa", "groupb")
	if err != nil {
		t.Errorf("error while preparing user: %v %v", err, cfg.Users())
	}
	u.UpdateGroups(cfg, []string{})
	err = cfg.AddUser(u)
	if err != nil {
		t.Errorf("error while registering user: %v %v", err, cfg.Users())
	}
	if len(cfg.Users()) != 4 {
		t.Errorf("Adding user did not work %v", JSON(cfg.Users()))
	}
	// set expected file content for removal
	sftp.testHosts["a:22"] = SFTPMockHost{Host: "b:22", User: "test", File: "ssh-rsa foo rootuser\nssh-rsa bar1 user-a.com\n"}
	cfg.DeleteUser("bar@email")
	if len(cfg.Users()) != 3 {
		t.Errorf("Deleting user did not work %v", len(cfg.Users()))
	}
	if cfg.DeleteUser("bar@email") {
		t.Errorf("Deleting user did work again")
	}
}

func TestModifyUserGroups(t *testing.T) {
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
	}, []*User{
		{Email: "foo@email", KeyType: "ssh-rsa", Key: "foo", Name: "rootuser", Groups: []string{"groupa"}},
		{Email: "user-a.com@a", KeyType: "ssh-rsa", Key: "bar1", Name: "user-a.com", Groups: []string{"groupb"}},
		{Email: "rootuser@b", KeyType: "ssh-rsa", Key: "foo2", Name: "aroot", Groups: []string{"groupc"}},
	}, sftp)
	u, err := cfg.PrepareUser("bar@email", "fixtures/dummy.key", "groupa", "groupb")
	if err != nil {
		t.Errorf("error while preparing user: %v %v", err, cfg.Users())
	}
	u.UpdateGroups(cfg, []string{})
	err = cfg.AddUser(u)
	if err != nil {
		t.Errorf("error while registering user: %v %v", err, cfg.Users())
	}
	if len(cfg.Users()) != 4 {
		t.Errorf("Adding user did not work #1 %v", JSON(cfg.Users()))
	}
	g := cfg.GetGroups()
	if grp, ok := g["groupa"]; ok {
		if !contains(grp.Users, "bar@email") {
			t.Errorf("Group users doesn't have bar@email: %v", JSON(grp.Users))
		}
	} else {
		t.Errorf("Group groupa doesn't exits: %v", g)
	}
	cfg.DeleteUser("bar@email")
	g = cfg.GetGroups()
	if grp, ok := g["groupa"]; ok {
		if contains(grp.Users, "bar@email") {
			t.Errorf("Group users still have bar@email: %v", grp.Users)
		}
	} else {
		t.Errorf("Group groupa doesn't exits: %v", g)
	}
	if len(cfg.Users()) != 3 {
		t.Errorf("Deleting user did not work #2 %v", len(cfg.Users()))
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
	}, []*User{
		{Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot"},
	}, sftp)
	_, err := cfg.PrepareUser("bar@email", "fixtures/broken.key")
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
	}, []*User{
		{Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot"},
	}, sftp)
	sftp.SetError(true)
	_, err := cfg.PrepareUser("bar@email", "fixtures/nonexistent.key", "groupa")
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
	}, []*User{
		{Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot"},
	}, sftp)
	cfg.Update("a", "b")
	h := cfg.GetHost("a")
	if h == nil {
		t.Errorf("host a not found")
	}
	if len(h.Users) != 2 {
		t.Errorf("Update() did not work, host a users: %v", h.Users)
	}
	h = cfg.GetHost("b")
	if h == nil {
		t.Errorf("host b not found")
	}
	if len(h.Users) != 2 {
		t.Errorf("Update() did not work, host b users: %v", h.Users)
	}
}

func TestGetGroups(t *testing.T) {
	cfg := testConfig("foo", map[string]*Host{
		"a": {Alias: "a", Host: "a:22", User: "aroot", Groups: []string{"a", "b"}},
		"b": {Alias: "b", Host: "b:22", User: "aroot", Groups: []string{"a", "c"}},
	}, []*User{
		{Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot", Groups: []string{"a"}},
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
	}, []*User{
		{Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot", Groups: []string{"a"}},
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

func TestUpdateGroup(t *testing.T) {
	cfg := testConfig("foo", map[string]*Host{
		"a": {Alias: "a", Host: "a:22", User: "aroot", Groups: []string{"a", "b"}},
		"b": {Alias: "b", Host: "a:22", User: "aroot", Groups: []string{"a", "c"}},
	},
		[]*User{
			{Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot", Groups: []string{"a"}},
		}, &SFTPConn{mock: true})
	cfg.UpdateGroup("a", []string{"b"}, []string{"c"})
	h := cfg.GetHost("a")
	if h == nil {
		t.Errorf("host a not found")
		return
	}
	if len(h.Groups) != 1 {
		t.Errorf("UpdateGroup did not work, host a groups: %v", h.Groups)
	}
	h = cfg.GetHost("b")
	if h == nil {
		t.Errorf("host b not found")
	}
	if len(h.Groups) != 1 {
		t.Errorf("UpdateGroup did not work, host b groups: %v", h.Groups)
	}
}

func TestDeleteGroup(t *testing.T) {
	cfg := testConfig("foo", map[string]*Host{
		"a": {Alias: "a", Host: "a:22", User: "aroot", Groups: []string{"a", "b"}},
		"b": {Alias: "b", Host: "a:22", User: "aroot", Groups: []string{"a", "c"}},
	},
		[]*User{
			{Email: "foo@email", KeyType: "ssh-rsa", Key: "keydata", Name: "aroot", Groups: []string{"a"}},
		}, &SFTPConn{mock: true})
	cfg.DeleteGroup("a")
	h := cfg.GetHost("a")
	if h == nil {
		t.Errorf("host a not found")
	}
	if len(h.Groups) != 1 {
		t.Errorf("DeleteGroup did not work, host a groups: %v", h.Groups)
	}
	h = cfg.GetHost("b")
	if h == nil {
		t.Errorf("host b not found")
	}
	if len(h.Groups) != 1 {
		t.Errorf("DeleteGroup did not work, host b groups: %v", h.Groups)
	}
}
