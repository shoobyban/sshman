package backend

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

type Config struct {
	StorageFilePath string `mapstructure:"STORAGE"`
}

var config *Config

type storageFile struct {
	Key   string           `json:"key"`
	Hosts map[string]*Host `json:"hosts"`
	Users map[string]*User `json:"users"`
}

// Storage is the main storage (data storage) for the sshman backend
type Storage struct {
	l          sync.Mutex
	key        string
	hosts      map[string]*Host
	users      map[string]*User
	Groups     map[string]Group
	Conn       SFTP
	persistent bool
	Log        *ILog
}

// LabelGroup is used for returning group information to frontend
type LabelGroup struct {
	Label string   `json:"label"`
	Hosts []string `json:"hosts"`
	Users []string `json:"users"`
}

// Group is used to store group information (connecting hosts and users)
type Group struct {
	Name  string
	Size  int
	Hosts []*Host
	Users []*User
}

func SetConfig(c *Config) {
	config = c
}

// NewStorage creates a new storage with a logger
func NewStorage() *Storage {
	return &Storage{
		hosts: map[string]*Host{},
		users: map[string]*User{},
		Conn:  &SFTPConn{},
		Log:   NewLog(false),
	}
}

// newStorageWithLog creates a new storage with a given logger, used for frontend
func newStorageWithLog(log *ILog) *Storage {
	return &Storage{
		hosts: map[string]*Host{},
		users: map[string]*User{},
		Conn:  &SFTPConn{},
		Log:   log,
	}
}

// ReadStorage reads the storage file ~/.ssh/.sshman and returns a new Storage
func ReadStorage() *Storage {
	c := NewStorage()
	err := c.load(config.StorageFilePath)
	if err != nil {
		c.Log.Infof("No storage file " + config.StorageFilePath + ", creating one")
		return c
	}
	return c
}

// ReadStorageWithLog reads the storage file ~/.ssh/.sshman and returns a new Storage with a logger, used for web
func ReadStorageWithLog(log *ILog) *Storage {
	c := newStorageWithLog(log)
	c.load(config.StorageFilePath)
	return c
}

func (c *Storage) load(filename string) error {
	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	c.persistent = true // testing doesn't have this where we just create the config
	var cf storageFile
	err = json.Unmarshal(b, &cf)
	if err != nil {
		log.Fatalf("Error: unable to decode into struct, please correct or remove broken %s %v\n", config.StorageFilePath, err)
	}
	c.key = cf.Key
	c.hosts = cf.Hosts
	c.users = cf.Users
	for alias, host := range c.hosts {
		host.Alias, host.Config = alias, c
		c.hosts[alias] = host
	}
	c.updateGroups()
	return nil
}

func (c *Storage) updateGroups() {
	groups := map[string]Group{}
	for _, host := range c.Hosts() {
		for _, group := range host.Groups {
			if v, ok := groups[group]; ok {
				v.Hosts = append(v.Hosts, host)
				groups[group] = v
			} else {
				groups[group] = Group{Hosts: []*Host{host}}
			}
		}
	}
	for _, user := range c.Users() {
		for _, group := range user.Groups {
			if _, ok := groups[group]; ok {
				g := groups[group]
				g.Users = append(g.Users, user)
				groups[group] = g
			} else {
				groups[group] = Group{Users: []*User{user}}
			}
		}
	}
	c.Groups = groups
}

// GetUserByEmail get a user from config by email as we store them by key checksum
func (c *Storage) GetUserByEmail(email string) (string, *User) {
	for key, user := range c.Users() {
		if user.Email == email {
			return key, user
		}
	}
	//	c.Log.Errorf("No user with email %s found", email)
	return "", nil
}

// Write storage file into ~/.ssh/.sshman (if not testing)
func (c *Storage) Write() {
	if !c.persistent {
		return // when testing (so not from ReadConfig)
	}
	cf := storageFile{Key: c.key, Hosts: c.hosts, Users: c.users}
	b, _ := json.MarshalIndent(cf, "", "  ")
	os.WriteFile(config.StorageFilePath, b, 0600)
	c.Log.Infof("storage saved to " + config.StorageFilePath)
}

func (c *Storage) getHosts(group string) []*Host {
	var hosts []*Host
	for _, host := range c.Hosts() {
		if contains(host.GetGroups(), group) {
			hosts = append(hosts, host)
		}
	}
	return hosts
}

// GetUsers will return users that have the given group
func (c *Storage) GetUsers(group string) []*User {
	var users []*User
	for _, user := range c.Users() {
		if group == "" || contains(user.Groups, group) {
			users = append(users, user)
		}
	}
	return users
}

// AddUserToHosts adds user to all allowed hosts' authorized_keys files
func (c *Storage) AddUserToHosts(newuser *User) {
	for alias, host := range c.Hosts() {
		if match(host.GetGroups(), newuser.Groups) {
			c.Log.Infof("adding %s to %s", newuser.Email, alias)
			host.AddUser(newuser)
		}
	}
	c.Write()
}

// SetHost is a setter for hosts
func (c *Storage) SetHost(alias string, host *Host) {
	c.l.Lock()
	defer c.l.Unlock()
	host.Config = c
	c.hosts[alias] = host
	c.Write()
}

// Hosts is a getter for hosts
func (c *Storage) Hosts() map[string]*Host {
	c.l.Lock()
	defer c.l.Unlock()
	return c.hosts
}

// Users is a getter for users
func (c *Storage) Users() map[string]*User {
	c.l.Lock()
	defer c.l.Unlock()
	return c.users
}

// UserExists checks if a user exists ignoring the result
func (c *Storage) UserExists(lsum string) bool {
	_, ok := c.users[lsum]
	return ok
}

// DelUserFromHosts removes user's key from all hosts' authorized_keys files
func (c *Storage) DelUserFromHosts(deluser *User) error {
	if deluser == nil {
		return fmt.Errorf("User is nil")
	}
	for alias, host := range c.Hosts() {
		err := host.DelUser(deluser)
		if err != nil {
			c.Log.Errorf("Can't delete user %s from host %s %v", deluser.Email, host.Alias, err)
			continue
		}
		c.SetHost(alias, host)
	}
	c.Write()
	return nil
}

// PrepareHost will prepare a host entry from array of strings
// Args: alias, hostname, user, keyfile, groups
func (c *Storage) PrepareHost(args ...string) (*Host, error) {
	alias := args[0]
	if _, err := os.Stat(args[3]); os.IsNotExist(err) {
		wd, _ := os.Getwd()
		c.Log.Errorf("key file '%s' is not in %s", args[3], wd)
		return nil, fmt.Errorf("no such file '%s'", args[3])
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
	return host, nil
}

// AddHost adds a host to the storage
func (c *Storage) AddHost(host *Host, withUsers bool) error {
	c.l.Lock()
	defer c.l.Unlock()
	c.Log.Infof("Adding host %s", host.Alias)
	if _, ok := c.hosts[host.Alias]; ok {
		return fmt.Errorf("Host %s already exists", host.Alias)
	}
	c.hosts[host.Alias] = host
	c.Write()
	if withUsers {
		users, err := host.ReadUsers()
		if err != nil {
			return err
		}
		for _, user := range users {
			c.AddUser(user)
		}
	}
	return nil
}

// DeleteUserByID removes a user from the storage
func (c *Storage) DeleteUserByID(id string) bool {
	c.l.Lock()
	defer c.l.Unlock()
	var ok bool
	if _, ok = c.users[id]; ok {
		delete(c.users, id)
		c.Write()
		c.Log.Infof("Deleted user %s", id)
	}
	return ok
}

// DeleteUser removes a user from the storage
func (c *Storage) DeleteUser(email string) bool {
	c.l.Lock()
	defer c.l.Unlock()
	found := false
	for id, user := range c.users {
		if email == user.Email {
			delete(c.users, id)
			c.Write()
			c.Log.Infof("Deleted user %s", email)
			found = true
		}
	}
	if !found {
		c.Log.Errorf("No user with email %s found", email)
	}
	return found
}

// PrepareUser will prepare a user entry from array of strings
// New user: old groups, email, key file, new groups
func (c *Storage) PrepareUser(args ...string) (*User, error) {
	parts, err := readKeyFile(args[1])
	if err != nil {
		return nil, err
	}
	newuser := NewUser(args[0], parts[0], parts[1], parts[2])
	newuser.Groups = args[2:]
	newuser.Config = c
	return newuser, nil
}

// AddUser adds a user to the config
func (c *Storage) AddUser(newuser *User) error {
	c.l.Lock()
	defer c.l.Unlock()
	if newuser == nil {
		return fmt.Errorf("User is nil")
	}
	lsum := checksum(newuser.Key)
	if _, ok := c.users[lsum]; ok {
		return fmt.Errorf("user with key %s already exists", lsum)
	}
	c.users[lsum] = newuser
	c.Write()
	return nil
}

// UpdateUser finds and replaces user
func (c *Storage) UpdateUser(newuser *User) error {
	if newuser.Email == "" {
		return fmt.Errorf("no email provided")
	}
	oldKey, oldUser := c.GetUserByEmail(newuser.Email)
	if oldUser == nil {
		return fmt.Errorf("user %s not found", newuser.Email)
	}
	c.l.Lock()
	defer c.l.Unlock()
	delete(c.users, oldKey)
	c.Log.Infof("deleting user %s with key %s", oldUser.Email, oldKey)
	lsum := checksum(newuser.Key)
	c.users[lsum] = newuser
	c.Log.Infof("key %s with content %v", lsum, newuser)
	c.Write()
	return nil
}

// DeleteHost removes a host from the config
func (c *Storage) DeleteHost(alias string) bool {
	c.l.Lock()
	defer c.l.Unlock()
	if _, ok := c.hosts[alias]; ok {
		delete(c.hosts, alias)
		c.Write()
		return true
	}
	return false
}

// Update updates all hosts by given aliases
func (c *Storage) Update(aliases ...string) {
	c.l.Lock()
	hosts := c.hosts
	if len(aliases) > 0 {
		hosts = map[string]*Host{}
		for _, a := range aliases {
			if host, ok := c.hosts[a]; ok {
				hosts[a] = host
			}
		}
	}
	c.l.Unlock()
	for _, host := range hosts {
		users, err := host.ReadUsers()
		if err != nil {
			c.Log.Errorf("Can't read users from host %s: %v", host.Alias, err)
			continue
		}
		for _, user := range users {
			c.AddUser(user)
		}
	}
}

// Regenerate updates group information for given hosts
func (c *Storage) Regenerate(aliases ...string) {
	c.l.Lock()
	defer c.l.Unlock()
	if len(aliases) > 0 {
		for _, a := range aliases {
			if host, ok := c.hosts[a]; ok {
				host.UpdateGroups(c, []string{})
			}
		}
	}
}

// GetGroups returns all groups for frontend (web and cli)
func (c *Storage) GetGroups() map[string]LabelGroup {
	c.l.Lock()
	defer c.l.Unlock()
	groups := map[string]LabelGroup{}
	for alias, host := range c.hosts {
		for _, group := range host.Groups {
			if v, ok := groups[group]; ok {
				v.Hosts = append(v.Hosts, alias)
				groups[group] = v
			} else {
				groups[group] = LabelGroup{Label: group, Hosts: []string{alias}}
			}
		}
	}
	for _, user := range c.users {
		for _, group := range user.Groups {
			if _, ok := groups[group]; ok {
				g := groups[group]
				g.Users = append(g.Users, user.Email)
				groups[group] = g
			} else {
				groups[group] = LabelGroup{Label: group, Users: []string{user.Email}}
			}
		}
	}
	return groups
}

// AddUserByEmail adds a user to the hosts indicated by config groups
func (c *Storage) AddUserByEmail(email string) bool {
	_, u := c.GetUserByEmail(email)
	if u != nil {
		c.AddUserToHosts(u)
		return true
	}
	return false
}

// GetUser is getting user by id
func (c *Storage) GetUser(lsum string) *User {
	c.l.Lock()
	defer c.l.Unlock()
	return c.users[lsum]
}

// GetHost is getting host by alias
func (c *Storage) GetHost(alias string) *Host {
	c.l.Lock()
	defer c.l.Unlock()
	return c.hosts[alias]
}

// DeleteGroup removes a group from the config, removing all users and hosts group labels
func (c *Storage) DeleteGroup(label string) bool {
	c.l.Lock()
	defer c.l.Unlock()
	if _, ok := c.Groups[label]; ok {
		// loop through hosts and remove group from host
		for _, host := range c.hosts {
			host.Groups = remove(host.Groups, label)
		}
		// loop through users and remove group from user
		for _, user := range c.users {
			user.Groups = remove(user.Groups, label)
		}

		delete(c.Groups, label)
		c.Write()
		return true
	}
	return false
}

// UpdateGroup updates a group in the config
// removing all users and hosts group labels then adding them again
func (c *Storage) UpdateGroup(groupLabel string, users, servers []string) {
	c.l.Lock()
	// loop through hosts and add groupLabel to host.Groups
	for _, host := range c.hosts {
		if !contains(host.Groups, groupLabel) {
			host.Groups = append(host.Groups, groupLabel)
		}
	}
	// loop through users and add groupLabel to user.Groups
	for _, user := range c.users {
		if !contains(user.Groups, groupLabel) {
			user.Groups = append(user.Groups, groupLabel)
		}
	}
	// loop through c.Users and remove item if not in users
	for _, user := range c.users {
		if !contains(users, user.Email) {
			user.Groups = remove(user.Groups, groupLabel)
		}
	}
	// loop through c.Hosts and remove item if not in servers
	for _, host := range c.hosts {
		if !contains(servers, host.Alias) {
			host.Groups = remove(host.Groups, groupLabel)
		}
	}
	// update c.Groups.Users and c.Groups.Hosts
	c.l.Unlock()
	c.updateGroups()
	c.Write()
}
