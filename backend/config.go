package backend

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Group struct {
	Users   []string
	Servers []string
}

type config struct {
	Key        string               `json:"key"`
	Hosts      map[string]Hostentry `json:"hosts"`
	Users      map[string]User      `json:"users"`
	conn       SFTP                 `json:"-"`
	persistent bool                 `json:"-"`
}

// initConfig reads in config file and ENV variables if set.
func ReadConfig(conn SFTP) *config {
	var C config
	home, _ := os.UserHomeDir()
	b, err := ioutil.ReadFile(home + "/.ssh/.sshman")
	if err != nil {
		log.Printf("Error: unable to read .sshman, %v\n", err)
	}
	err = json.Unmarshal(b, &C)
	if err != nil {
		log.Printf("Error: unable to decode into struct, %v\n", err)
	}
	for alias, host := range C.Hosts {
		host.Alias = alias
		host.Config = &C
		C.Hosts[alias] = host
	}
	C.conn = conn
	C.persistent = true
	return &C
}

// GetUserByEmail get a user from config by email
func (c *config) GetUserByEmail(email string) (string, *User) {
	for key, user := range c.Users {
		if user.Email == email {
			return key, &user
		}
	}
	return "", nil
}

func (c *config) Write() {
	if !c.persistent {
		return // when testing (so not from ReadConfig)
	}
	b, _ := json.MarshalIndent(c, "", "  ")
	home, _ := os.UserHomeDir()
	ioutil.WriteFile(home+"/.ssh/.sshman", b, 0644)
}

func (c *config) getServers(group string) []Hostentry {
	var servers []Hostentry
	for _, host := range c.Hosts {
		if contains(host.Groups, group) {
			servers = append(servers, host)
		}
	}
	return servers
}

func (c *config) getUsers(group string) []User {
	var users []User
	for _, user := range c.Users {
		if contains(user.Groups, group) {
			users = append(users, user)
		}
	}
	return users
}

// AddUserToHosts adds user to all allowed hosts' authorized_keys files
func (c *config) AddUserToHosts(newuser *User) {
	for alias, host := range c.Hosts {
		if len(newuser.Groups) != 0 && !match(newuser.Groups, host.Groups) {
			continue
		}
		userlist := []string{}
		have := false
		sum, lines, err := host.read()
		if err != nil {
			log.Printf("Error: error reading authorized keys on %s: %v\n", alias, err)
			continue
		}
		for _, line := range lines {
			if len(line) == 0 {
				continue
			}
			parts := strings.Split(line, " ")
			if len(parts) != 3 {
				log.Printf("Error: Not good line: '%s'\n", line)
			}
			if parts[0] == newuser.KeyType &&
				parts[1] == newuser.Key &&
				parts[2] == newuser.Name {
				have = true
			}
			lsum := checksum(parts[1])
			if _, ok := c.Users[lsum]; !ok {
				c.Users[lsum] = User{
					KeyType: parts[0],
					Key:     parts[1],
					Name:    parts[2],
					Email:   parts[2] + "@" + alias,
				}
			}
			userlist = append(userlist, c.Users[lsum].Email)
		}

		if have {
			log.Printf("User %s is already on %s\n", newuser.Email, host.Alias)
			host.Checksum = sum
			host.Users = userlist
			c.Hosts[alias] = host
			continue
		} else {
			log.Printf("Adding %s to %s\n", newuser.Email, host.Alias)
			lines = deleteEmpty(append(lines, newuser.KeyType+" "+newuser.Key+" "+newuser.Name))

			err = c.conn.Write(strings.Join(lines, "\n") + "\n")
			if err != nil {
				log.Printf("Error: error writing %s: %v\n", alias, err)
			}
			userlist = append(userlist, newuser.Email)
		}
		host.Checksum = sum
		host.Users = userlist
		c.Hosts[alias] = host
	}
	c.Write()
}

// DelUserFromHosts removes user's key from all hosts' authorized_keys files
func (c *config) DelUserFromHosts(deluser *User) {
	for alias, host := range c.Hosts {
		err := host.delUser(deluser)
		if err != nil {
			log.Printf("Can't delete user %s from host %s %v\n", deluser.Email, host.Alias, err)
			continue
		}
		c.Hosts[alias] = host
	}
	c.Write()
}

// RegisterServer adds a server to the configuration
func (c *config) RegisterServer(args ...string) error {
	alias := args[0]
	if _, err := os.Stat(args[3]); os.IsNotExist(err) {
		log.Fatalf("no such file '%s'\n", args[3])
		return err
	}
	groups := args[4:]
	server := Hostentry{
		Host:   args[1],
		User:   args[2],
		Key:    args[3],
		Users:  []string{},
		Groups: groups,
		Alias:  alias,
		Config: c,
	}
	c.Hosts[alias] = server
	log.Printf("Registering %s to server %s with %s user\n", alias, args[1], args[1])
	c.Write()
	server.readUsers()
	return nil
}

// UnregisterUser removes a user from the configuration
func (c *config) UnregisterUser(email string) bool {
	for id, user := range c.Users {
		if email == user.Email {
			delete(c.Users, id)
			c.Write()
			return true
		}
	}
	return false
}

// RegisterUser adds a user to the config
func (c *config) RegisterUser(oldgroups []string, args ...string) error {
	b, err := ioutil.ReadFile(args[1])
	if err != nil {
		log.Printf("Error: error reading public key file: '%s' %v\n", args[1], err)
		return err
	}
	parts := strings.Split(strings.TrimSuffix(string(b), "\n"), " ")
	if len(parts) != 3 {
		log.Printf("Error: not a proper public key file\n")
	}
	lsum := checksum(parts[1])
	groups := args[2:]
	newuser := User{
		KeyType: parts[0],
		Key:     parts[1],
		Name:    parts[2],
		Email:   args[0],
		Groups:  groups,
	}
	c.Users[lsum] = newuser
	log.Printf("Registering %s %s %s %s\n", parts[0], parts[2], args[0], lsum)
	c.Write()
	return newuser.UpdateGroups(c, oldgroups)
}

// UnregisterServer removes a server from the config
func (c *config) UnregisterServer(alias string) bool {
	if _, ok := c.Hosts[alias]; ok {
		delete(c.Hosts, alias)
		c.Write()
		return true
	}
	return false
}

func (c *config) Update(aliases ...string) {
	hosts := c.Hosts
	if len(aliases) > 0 {
		hosts = map[string]Hostentry{}
		for _, a := range aliases {
			if host, ok := c.Hosts[a]; ok {
				hosts[a] = host
			}
		}
	}
	for alias, host := range hosts {
		log.Printf("Reading %s\n", host.Alias)
		sum, lines, err := host.read()
		if err != nil {
			log.Printf("Error: error reading authorized keys on %s: %v\n", alias, err)
			continue
		}
		userlist := []string{}
		for _, line := range lines {
			if len(line) == 0 {
				continue
			}
			parts := strings.Split(line, " ")
			if len(parts) != 3 {
				log.Printf("Error: Not good line: '%s'\n", line)
			}
			lsum := checksum(parts[1])
			if _, ok := c.Users[lsum]; !ok {
				c.Users[lsum] = User{
					KeyType: parts[0],
					Key:     parts[1],
					Name:    parts[2],
					Email:   parts[2] + "@" + alias,
				}
			}
			userlist = append(userlist, c.Users[lsum].Email)
		}
		host.Checksum = sum
		host.Users = userlist
		c.Hosts[alias] = host
	}
	c.Write()
}

func (c *config) GetGroups() map[string]Group {
	groups := map[string]Group{}
	for alias, host := range c.Hosts {
		for _, group := range host.Groups {
			if v, ok := groups[group]; ok {
				v.Servers = append(v.Servers, alias)
				groups[group] = v
			} else {
				groups[group] = Group{Servers: []string{alias}}
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
