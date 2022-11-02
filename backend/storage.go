package backend

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type InternalConfig struct {
	StorageFilePath string `mapstructure:"STORAGE"`
}

var config = &InternalConfig{StorageFilePath: "teststorage"}

func SetConfig(c *InternalConfig) {
	fmt.Printf("Setting config to %v\n", c)
	config = c
}

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
	log        *ILog
	Stop       bool
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

// NewStorage creates a new storage with a logger
func NewStorage(persistent bool) *Storage {
	return &Storage{
		hosts:      map[string]*Host{},
		users:      map[string]*User{},
		Conn:       &SFTPConn{},
		log:        NewLog(false),
		persistent: persistent,
	}
}

func NewTestStorage() *Storage {
	return &Storage{
		hosts:      map[string]*Host{},
		users:      map[string]*User{},
		Conn:       &SFTPConn{mock: true},
		log:        NewLog(false),
		persistent: false,
	}
}

// newStorageWithLog creates a new storage with a given logger, used for frontend
func newStorageWithLog(log *ILog) *Storage {
	return &Storage{
		hosts:      map[string]*Host{},
		users:      map[string]*User{},
		Conn:       &SFTPConn{},
		log:        log,
		persistent: true,
	}
}

// ReadStorage reads the storage file ~/.ssh/.sshman and returns a new Storage
func ReadStorage() *Storage {
	c := NewStorage(true)
	err := c.load(config.StorageFilePath)
	if err != nil {
		c.log.Infof("No storage file " + config.StorageFilePath + ", creating one")
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
	if len(b) == 0 {
		return nil
	}
	c.persistent = true // testing doesn't have this where we just create the config
	var cf storageFile
	err = json.Unmarshal(b, &cf)
	if err != nil {
		log.Printf("Error: unable to decode into struct, please correct or remove broken %s %v\n", config.StorageFilePath, err)
		return err
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
			if v, ok := groups[group]; ok {
				v.Users = append(v.Users, user)
				groups[group] = v
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
	//	c.log.Errorf("No user with email %s found", email)
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
	// c.log.Infof("storage saved to " + config.StorageFilePath)
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
			c.log.Infof("adding %s to %s", newuser.Email, alias)
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

// GetUserByKey returns user if exists by ssh key
func (c *Storage) GetUserByKey(lsum string) *User {
	if v, ok := c.users[lsum]; ok {
		return v
	}
	return nil
}

// RemoveUserFromHosts removes user's key from all hosts' authorized_keys files
func (c *Storage) RemoveUserFromHosts(deluser *User) error {
	if deluser == nil {
		return fmt.Errorf("User is nil")
	}
	for alias, host := range c.Hosts() {
		err := host.RemoveUser(deluser)
		if err != nil {
			c.log.Errorf("Can't delete user %s from host %s %v", deluser.Email, host.Alias, err)
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
		c.log.Errorf("key file '%s' is not in %s", args[3], wd)
		return nil, fmt.Errorf("no such file '%s'", args[3])
	}
	groups := args[4:]
	host := &Host{
		Host:   args[1],
		User:   args[2],
		Key:    c.key,
		Users:  []*User{},
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
	c.log.Infof("Adding host %s", host.Alias)
	if _, ok := c.hosts[host.Alias]; ok {
		return fmt.Errorf("Host %s already exists", host.Alias)
	}
	c.hosts[host.Alias] = host
	c.Write()
	if withUsers {
		c.UpdateHost(host)
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
		c.log.Infof("Deleted user %s", id)
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
			c.log.Infof("Deleted user %s", email)
			found = true
		}
	}
	if !found {
		c.log.Errorf("No user with email %s found", email)
	}
	return found
}

// PrepareUser will prepare a user entry from array of strings
// New user: old groups, email, key file, new groups
func (c *Storage) PrepareUser(email, filename string, groups ...string) (*User, error) {
	parts, err := readKeyFile(filename)
	if err != nil {
		return nil, err
	}
	newuser := NewUser(email, parts[0], parts[1], parts[2])
	newuser.Groups = groups
	newuser.Config = c
	return newuser, nil
}

// AddUser adds a user to the config
func (c *Storage) AddUser(newuser *User, host string) error {
	c.l.Lock()
	defer c.l.Unlock()
	if newuser == nil {
		return fmt.Errorf("User is nil")
	}
	lsum := checksum(newuser.Key)
	if u, ok := c.users[lsum]; ok {
		if host != "" {
			c.log.Infof("User %s exists (%v)", u.Email, u.Hosts)
			diff := Difference(u.Hosts, []string{host})
			if len(diff[1]) > 0 {
				hh := c.hosts[host]
				if !hh.DueGroup(u) {
					u.Hosts = append(u.Hosts, host)
					c.users[lsum] = u
				}
			}
			return nil
		}
		return fmt.Errorf("user with key %s already exists", lsum)
	}
	if host != "" {
		newuser.Hosts = []string{host}
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
	//	c.log.Infof("deleting user %s with key %s", oldUser.Email, oldKey)
	lsum := checksum(newuser.Key)
	c.users[lsum] = newuser
	//	c.log.Infof("key %s with content %v", lsum, newuser)
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
	c.Stop = false
	c.l.Unlock()
	for _, host := range hosts {
		c.l.Lock()
		if c.Stop {
			c.l.Unlock()
			c.log.Infof("Received stop signal, stopping sync")
			return
		}
		c.l.Unlock()
		c.log.Infof("Updating host %s...", host.Alias)
		// check Stop channel
		c.UpdateHost(host)
	}
	c.Write()
}

// UpdateHost reads users from host and adds them to users list
func (c *Storage) UpdateHost(host *Host) error {
	users, _, err := host.ReadUsers()
	if err != nil {
		c.log.Errorf("Can't read users from host %s: %v", host.Alias, err)
		return err
	}
	for _, user := range users {
		c.AddUser(user, host.Alias)
		host.Users = append(host.Users, user)
	}
	host.LastUpdated = time.Now()
	return nil
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
			if v, ok := groups[group]; ok {
				v.Users = append(v.Users, user.Email)
				groups[group] = v
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
	if _, ok := c.Groups[label]; !ok {
		c.log.Errorf("group %s not found", label)
		return false
	}
	// loop through hosts and remove group from host
	for _, host := range c.hosts {
		host.Groups = remove(host.Groups, label)
	}
	// loop through users and remove group from user
	for _, user := range c.users {
		user.Groups = remove(user.Groups, label)
	}

	delete(c.Groups, label)
	c.log.Infof("deleted group %s", label)
	c.Write()
	return true
}

// UpdateGroup updates a group in the config
// removing all users and hosts group labels then adding them again
func (c *Storage) UpdateGroup(groupLabel string, users, hosts []string) {
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
		if !contains(hosts, host.Alias) {
			host.Groups = remove(host.Groups, groupLabel)
		}
	}
	// update c.Groups.Users and c.Groups.Hosts
	c.l.Unlock()
	c.updateGroups()
	c.Write()
}

// FromGroup returns true if user (by email) is in same group as host
func (c *Storage) FromGroup(host *Host, email string) bool {
	_, user := c.GetUserByEmail(email)
	if user == nil {
		return false
	}
	for _, g := range host.Groups {
		if contains(user.Groups, g) {
			return true
		}
	}
	return false
}

// StopSync stops Update() loop
func (c *Storage) StopUpdate() {
	c.log.Infof("Stopping sync...")
	c.l.Lock()
	defer c.l.Unlock()
	c.Stop = true
}

func (c *Storage) WatchFile(notify func()) {
	initWG := sync.WaitGroup{}
	initWG.Add(1)
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()
		eventsWG := sync.WaitGroup{}
		eventsWG.Add(1)
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						eventsWG.Done()
						return
					}
					const writeOrCreateMask = fsnotify.Write | fsnotify.Create
					if filepath.Clean(event.Name) == config.StorageFilePath &&
						event.Op&writeOrCreateMask != 0 {
						err := c.load(config.StorageFilePath)
						if err != nil {
							log.Printf("error reading config file: %v\n", err)
						}
						notify()
					} else if filepath.Clean(event.Name) == config.StorageFilePath &&
						event.Op&fsnotify.Remove != 0 {
						eventsWG.Done()
						return
					}

				case err, ok := <-watcher.Errors:
					if ok {
						log.Printf("watcher error: %v\n", err)
					}
					eventsWG.Done()
					return
				}
			}
		}()
		configDir := filepath.Dir(config.StorageFilePath)
		watcher.Add(configDir)
		initWG.Done()
		eventsWG.Wait()
	}()
	initWG.Wait()
}

func (c *Storage) Log() *ILog {
	return c.log
}
