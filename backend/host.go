package backend

import (
	"fmt"
	"strings"
)

// Host holds and manages a host entry in the config
type Host struct {
	Host     string   `json:"host"`
	User     string   `json:"user"`
	Key      string   `json:"key"`
	Checksum string   `json:"checksum"`
	Alias    string   `json:"alias"`
	Config   *Storage `json:"-"`
	Users    []string `json:"users"`
	Groups   []string `json:"groups"`
}

// ReadUsers reads any new users from the host's authorized_keys file
func (h *Host) ReadUsers() (map[string]*User, error) {
	sum, lines, err := h.read()
	if err != nil {
		return nil, err
	}
	newUsers := map[string]*User{}
	var userlist []string
	for _, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) == 3 {
			lsum := checksum(parts[1])
			email := parts[2] + "@" + h.Alias
			if !h.Config.UserExists(lsum) {
				user := &User{
					KeyType: parts[0],
					Key:     parts[1],
					Name:    parts[2],
					Email:   email,
					Config:  h.Config,
				}
				newUsers[lsum] = user
			}
			userlist = append(userlist, email)
		}
	}
	h.Checksum = sum
	h.Users = userlist
	h.Config.Write()
	return newUsers, nil
}

func (h *Host) read() (string, []string, error) {
	if h.Config == nil {
		return "", nil, fmt.Errorf("host is nil")
	}
	err := h.Config.Conn.Connect(h.Key, h.Host, h.User)
	if err != nil {
		return "", nil, fmt.Errorf("error connecting %s: %v", h.Alias, err)
	}
	defer h.Config.Conn.Close()
	b, err := h.Config.Conn.Read()
	if err != nil {
		return "", nil, fmt.Errorf("error reading authorized keys on %s: %v", h.Alias, err)
	}
	sum := checksum(string(b))
	lines := deleteEmpty(strings.Split(string(b), "\n"))
	return sum, lines, nil
}

func (h *Host) write(lines []string) error {
	if len(lines) == 0 {
		return fmt.Errorf("no keys in new file for host '%s', host would be inaccessible", h.Alias)
	}
	err := h.Config.Conn.Connect(h.Key, h.Host, h.User)
	if err != nil {
		return fmt.Errorf("error connecting %s: %v", h.Alias, err)
	}
	defer h.Config.Conn.Close()
	return h.Config.Conn.Write(strings.Join(lines, "\n") + "\n")
}

// GetUsers is a getter for host's Users
func (h *Host) GetUsers() []string {
	return h.Users
}

// GetGroups is a getter for host's Groups
func (h *Host) GetGroups() []string {
	return h.Groups
}

// SetGroups is a setter for host's Groups (overwrite all groups at once)
func (h *Host) SetGroups(groups []string) {
	h.Groups = groups
}

// HasMatchingGroups checks if the host has any matching groups with the user
func (h *Host) HasMatchingGroups(user *User) bool {
	return match(h.GetGroups(), user.GetGroups())
}

// HasUser checks if the host has a user with the given email
func (h *Host) HasUser(email string) bool {
	for _, e := range h.Users {
		if e == email {
			return true
		}
	}
	return false
}

// AddUser adds a user to the host's authorized_keys file
func (h *Host) AddUser(u *User) error {
	h.Users = append(h.Users, u.Email)
	var lines []string
	for _, email := range h.Users {
		// double-check old user entries
		_, userentry := h.Config.GetUserByEmail(email)
		if userentry == nil {
			// Shall we add a warning here?
			continue
		}
		lines = append(lines, userentry.KeyType+" "+userentry.Key+" "+userentry.Name)
	}
	return h.write(lines)
}

// DelUser removes a user from the host's authorized_keys file
func (h *Host) DelUser(u *User) error {
	if u == nil {
		return fmt.Errorf("user is nil")
	}
	sum, lines, err := h.read()
	if err != nil {
		return err
	}
	userlist := []string{}
	found := false
	newlines := []string{}
	for _, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) == 3 {
			lsum := checksum(parts[1])
			h.Config.Log.Infof("deleting user with checksum: %s, key: %s, userkey: %s", lsum, parts[1], u.Key)
			if parts[1] == u.Key {
				found = true
				continue
			}
			user := h.Config.GetUser(lsum)
			if user != nil {
				newlines = append(newlines, line)
				userlist = append(userlist, user.Email)
			}
		}
	}

	if found {
		h.Config.Log.Infof("found user, deleting %v", u)
		newlines = deleteEmpty(newlines)
		err = h.write(newlines)
		if err != nil {
			return fmt.Errorf("error writing %s: %v", h.Alias, err)
		}
	}
	h.Checksum = sum
	h.Users = userlist
	return nil
}

// UpdateGroups updates the host's groups based on old groups
func (h *Host) UpdateGroups(c *Storage, oldgroups []string) bool {
	success := true
	added, removed := updates(oldgroups, h.Groups)
	h.Config.Log.Infof("added: %v removed: %v", added, removed)
	for _, group := range added {
		users := c.GetUsers(group)
		for _, u := range users {
			if !h.HasUser(u.Email) {
				h.Config.Log.Infof("Adding %s (group %s) to %s", u.Email, group, h.Alias)
				err := h.AddUser(u)
				if err != nil {
					h.Config.Log.Errorf("error adding %s to %s", u.Email, h.Alias)
					success = false
					continue
				}
				h.Config.Log.Infof("Added %s to %s (host groups %v)", u.Email, h.Alias, h.Groups)
			}
		}
	}

	for _, group := range removed {
		users := c.GetUsers(group)
		for _, u := range users {
			// are there other groups that keep user on host
			if h.HasMatchingGroups(u) {
				continue
			}
			if h.HasUser(u.Email) {
				err := h.DelUser(u)
				if err != nil {
					c.Log.Errorf("error removing %s from %s", u.Email, h.Alias)
					success = false
					continue
				}
				c.Log.Infof("removed %s from %s", u.Email, h.Alias)
			}
		}
	}
	c.SetHost(h.Alias, h)
	c.Write()
	return success
}
