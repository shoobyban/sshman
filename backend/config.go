package backend

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type Group struct {
	Users []string
	Hosts []string
}

type Storage struct {
	Key        string           `json:"key"`
	Hosts      map[string]*Host `json:"hosts"`
	Users      map[string]*User `json:"users"`
	conn       SFTP             `json:"-"`
	persistent bool             `json:"-"`
	home       string           `json:"-"`
}

func ReadConfig() *Storage {
	C := Storage{Hosts: map[string]*Host{}, Users: map[string]*User{}, conn: &SFTPConn{}}
	C.home, _ = os.UserHomeDir()
	b, err := os.ReadFile(C.home + "/.ssh/.sshman")
	if err != nil {
		C.persistent = true // testing doesn't have this where we just create the config
		fmt.Println("No configuration file ~/.ssh/.sshman, creating one")
		return &C
	}
	err = json.Unmarshal(b, &C)
	if err != nil {
		log.Fatalf("Error: unable to decode into struct, please correct or remove broken ~/.ssh/.sshman %v\n", err)
	}
	C.persistent = true // testing doesn't have this where we just create the config
	for alias, host := range C.Hosts {
		host.Alias, host.Config = alias, &C
		C.Hosts[alias] = host
	}
	return &C
}

// GetUserByEmail get a user from config by email as we store them by key checksum
func (c *Storage) GetUserByEmail(email string) (string, *User) {
	for key, user := range c.Users {
		if user.Email == email {
			return key, user
		}
	}
	return "", nil
}

// Write configuration file into ~/.ssh/.sshman (if not testing)
func (c *Storage) Write() {
	if !c.persistent {
		return // when testing (so not from ReadConfig)
	}
	b, _ := json.MarshalIndent(c, "", "  ")
	os.WriteFile(c.home+"/.ssh/.sshman", b, 0644)
}

func (c *Storage) getHosts(group string) []*Host {
	var hosts []*Host
	for _, host := range c.Hosts {
		if contains(host.GetGroups(), group) {
			hosts = append(hosts, host)
		}
	}
	return hosts
}

// getUsers will return users that have the given group
func (c *Storage) GetUsers(group string) []*User {
	var users []*User
	for _, user := range c.Users {
		if contains(user.Groups, group) {
			fmt.Printf("Checking for %s, User %s has %v\n", group, user.Email, user.Groups)
			users = append(users, user)
		}
	}
	return users
}

// AddUserToHosts adds user to all allowed hosts' authorized_keys files
func (c *Storage) AddUserToHosts(newuser *User) {
	for alias, host := range c.Hosts {
		if match(host.GetGroups(), newuser.Groups) {
			fmt.Printf("Adding %s to %s\n", newuser.Email, alias)
			host.AddUser(newuser)
		}
	}
	c.Write()
}

// DelUserFromHosts removes user's key from all hosts' authorized_keys files
func (c *Storage) DelUserFromHosts(deluser *User) {
	for alias, host := range c.Hosts {
		host.ReadUsers()
		err := host.DelUser(deluser)
		if err != nil {
			fmt.Printf("Can't delete user %s from host %s %v\n", deluser.Email, host.Alias, err)
			continue
		}
		c.Hosts[alias] = host
	}
	c.Write()
}

// AddHost adds a host to the configuration
func (c *Storage) AddHost(args ...string) error {
	alias := args[0]
	if _, err := os.Stat(args[3]); os.IsNotExist(err) {
		return fmt.Errorf("no such file '%s'", args[3])
	}
	groups := args[4:]
	host := &Host{
		Host:   args[1],
		User:   args[2],
		Key:    args[3],
		Users:  []string{},
		Groups: groups,
		Alias:  alias,
		Config: c,
	}
	c.Hosts[alias] = host
	fmt.Printf("Adding %s to host %s with %s user\n", alias, args[1], args[1])
	c.Write()
	host.ReadUsers()
	return nil
}

// DeleteUser removes a user from the configuration
func (c *Storage) DeleteUser(email string) bool {
	for id, user := range c.Users {
		if email == user.Email {
			delete(c.Users, id)
			c.Write()
			return true
		}
	}
	return false
}

// New user: old groups, email, key file, new groups
func (c *Storage) PrepareUser(args ...string) (*User, error) {
	b, err := os.ReadFile(args[1])
	if err != nil {
		return nil, fmt.Errorf("error: error reading public key file: '%s' %v", args[1], err)
	}
	parts := strings.Split(strings.TrimSuffix(string(b), "\n"), " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("error: not a proper public key file")
	}
	groups := args[2:]
	newuser := &User{
		KeyType: parts[0],
		Key:     parts[1],
		Name:    parts[2],
		Email:   args[0],
		Groups:  groups,
	}
	return newuser, nil
}

// AddUser adds a user to the config
func (c *Storage) AddUser(newuser *User) error {
	lsum := checksum(newuser.Key)
	c.Users[lsum] = newuser
	c.Write()
	return nil
}

// DeleteHost removes a host from the config
func (c *Storage) DeleteHost(alias string) bool {
	if _, ok := c.Hosts[alias]; ok {
		delete(c.Hosts, alias)
		c.Write()
		return true
	}
	return false
}

func (c *Storage) Update(aliases ...string) {
	hosts := c.Hosts
	if len(aliases) > 0 {
		hosts = map[string]*Host{}
		for _, a := range aliases {
			if host, ok := c.Hosts[a]; ok {
				hosts[a] = host
			}
		}
	}
	for _, host := range hosts {
		host.ReadUsers()
	}
}

func (c *Storage) Regenerate(aliases ...string) {
	if len(aliases) > 0 {
		for _, a := range aliases {
			if host, ok := c.Hosts[a]; ok {
				host.UpdateGroups(c, []string{})
			}
		}
	}
}

func (c *Storage) GetGroups() map[string]Group {
	groups := map[string]Group{}
	for alias, host := range c.Hosts {
		for _, group := range host.Groups {
			if v, ok := groups[group]; ok {
				v.Hosts = append(v.Hosts, alias)
				groups[group] = v
			} else {
				groups[group] = Group{Hosts: []string{alias}}
			}
		}
	}
	for _, user := range c.Users {
		for _, group := range user.Groups {
			if _, ok := groups[group]; ok {
				g := groups[group]
				g.Users = append(g.Users, user.Email)
				groups[group] = g
			} else {
				groups[group] = Group{Users: []string{user.Email}}
			}
		}
	}
	return groups
}

func (c *Storage) AddUserByEmail(email string) bool {
	_, u := c.GetUserByEmail(email)
	if u != nil {
		c.AddUserToHosts(u)
		return true
	}
	return false
}

func (c *Storage) GetUser(key string) *User {
	return c.Users[key]
}
