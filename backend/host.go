package backend

import (
	"fmt"
	"strings"
	"time"
)

const AUTHFILE = ".ssh/authorized_keys"

// Host holds and manages a host entry in the config
type Host struct {
	Host        string              `json:"host"`
	User        string              `json:"user"`
	Key         string              `json:"key"`
	Protected   map[string]struct{} `json:"protected"`
	Alias       string              `json:"alias"`
	Config      *Data               `json:"-"`
	Users       []*User             `json:"userlist"`
	Groups      []string            `json:"groups"`
	LastUpdated time.Time           `json:"last_updated"`
	Checksum    string              `json:"checksum"`
	Modified    bool                `json:"modified"`
}

// ReadUsers reads ssh entries from the host's authorized_keys file,
// returning user list, checksum and error
func (h *Host) ReadUsers() ([]*User, string, error) {
	lines, err := h.read()
	if err != nil {
		return nil, "", err
	}
	sum := checksum(strings.Join(lines, "\n"))
	var userlist []*User
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			lsum := checksum(parts[1])
			email := ""
			if len(parts) > 2 {
				email = parts[2] + "@" + h.Alias
			}
			user := h.Config.GetUserByKey(lsum)
			if user == nil {
				name := ""
				if len(parts) > 2 {
					name = parts[2]
				}
				user = NewUser(email, parts[0], parts[1], name)
				user.Config = h.Config
			}
			userlist = append(userlist, user)
		}
	}
	return userlist, sum, nil
}

// UpdateUsersList sets the host's user list
func (h *Host) UpdateUsersList(userlist []*User) error {
	if h.Modified {
		return fmt.Errorf("host has been modified since last update, please refresh host first")
	}
	if len(userlist) > 0 && len(h.Protected) == 0 {
		// right after adding to host, all users should be protected (then we can untick them as we please)
		for _, u := range userlist {
			h.Protected = map[string]struct{}{u.Key: {}}
		}
	}
	uList := []*User{}
	for _, user := range userlist {
		if !h.Config.FromGroup(h, user.Email) {
			uList = append(uList, user)
		}
	}
	h.Users = uList
	h.LastUpdated = time.Now()
	h.Config.Write()
	return nil
}

func (h *Host) connectAndDo(action func() error) error {
	if h.Config == nil {
		return fmt.Errorf("host is nil")
	}
	err := h.Config.Conn.Connect(h.Key, h.Host, h.User)
	if err != nil {
		return fmt.Errorf("error connecting %s: %v", h.Alias, err)
	}
	defer h.Config.Conn.Close()
	return action()
}

// connects to host via ssh, downloads authorized_keys file content, updates Checksum field
func (h *Host) read() ([]string, error) {
	var lines []string
	err := h.connectAndDo(func() error {
		b, err := h.Config.Conn.Read(AUTHFILE)
		if err != nil {
			return fmt.Errorf("error reading authorized keys on %s: %v", h.Alias, err)
		}
		lines = deleteEmpty(strings.Split(string(b), "\n"))
		h.Checksum = checksum(strings.Join(lines, "\n"))
		return nil
	})
	return lines, err
}

// connects to host via ssh, uploads new authorized_keys file, updates Modified, LastUpdated and Checksum field
func (h *Host) write(lines []string) error {
	return h.connectAndDo(func() error {
		ls, err := h.read()
		if err != nil {
			return fmt.Errorf("error connecting %s: %v", h.Alias, err)
		}
		if checksum(strings.Join(ls, "\n")) != h.Checksum {
			return fmt.Errorf("host's authorized_keys file was modified since last update, please refresh hosts first")
		}
		if len(lines) == 0 {
			return fmt.Errorf("no keys in new file for host '%s', host would be inaccessible", h.Alias)
		}
		err = h.Config.Conn.Write(AUTHFILE, strings.Join(lines, "\n")+"\n")
		if err != nil {
			return err
		}
		h.LastUpdated = time.Now()
		h.Modified = false
		h.Checksum = checksum(strings.Join(lines, "\n"))
		return nil
	})
}

// GetUsers is a getter for host's Users
func (h *Host) GetUsers() []*User {
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
	for _, u := range h.Users {
		if u.Email == email {
			return true
		}
	}
	return false
}

// AddUser adds a user to the host's authorized_keys file
func (h *Host) AddUser(u *User) error {
	if h.LastUpdated.IsZero() {
		return fmt.Errorf("host was never read, can't update, please refresh hosts first")
	}
	h.Users = append(h.Users, u)
	if h.Config == nil || h.Config.Conn == nil {
		// No connection configured (likely in unit tests); don't attempt upload
		return nil
	}
	return h.Upload()
}

// Upload new authorized_keys file content
func (h *Host) Upload() error {
	var lines []string
	for _, user := range h.Users {
		// double-check old user entries
		lines = append(lines, user.KeyType+" "+user.Key+" "+user.Name)
	}
	lines = deleteEmpty(lines)
	h.Config.Log().Infof("updating %s", h.Alias)
	// return nil
	return h.write(lines)
}

// RemoveUser removes a user from the host's authorized_keys file
func (h *Host) RemoveUser(u *User) error {
	if h.LastUpdated.IsZero() {
		return fmt.Errorf("host was never read, can't update, please refresh hosts first")
	}
	if u == nil {
		return fmt.Errorf("user is nil")
	}
	userlist := []*User{}
	for _, user := range h.Users {
		if user.Key == u.Key {
			if _, protected := h.Protected[user.Key]; protected {
				return fmt.Errorf("user is protected, please remove protection first")
			}
			h.Config.Log().Infof("removing %s from %s", u.Email, h.Alias)
			h.Modified = true
			continue
		}
		userlist = append(userlist, user)
	}
	h.Users = userlist
	if h.Modified {
		if h.Config == nil || h.Config.Conn == nil {
			// No connection configured (likely in unit tests); don't attempt upload
			return nil
		}
		return h.Upload()
	}
	return nil
}

// DueGroup answers the question if user is on host due to group membership
func (h *Host) DueGroup(u *User) bool {
	diff := Difference(u.Groups, h.Groups)
	if len(diff[0]) < len(u.Groups) {
		return true
	}
	if len(diff[1]) < len(h.Groups) {
		return true
	}
	return false
}

// UpdateGroups updates the host's groups based on old groups
func (h *Host) UpdateGroups(cfg Config, oldgroups []string) bool {

	added, removed := splitUpdates(oldgroups, h.Groups)
	h.Config.Log().Infof("added: %v removed: %v", added, removed)

	success := processHostAdded(added, cfg, h)

	// are there other groups that keep user on host
	success = processHostRemoved(removed, cfg, h, success)
	cfg.SetHost(h.Alias, h)
	cfg.Write()
	return success
}

func processHostRemoved(removed []string, cfg Config, h *Host, success bool) bool {
	for _, group := range removed {
		users := cfg.GetUsers(group)
		for _, u := range users {

			if h.HasMatchingGroups(u) {
				continue
			}
			if h.HasUser(u.Email) {
				err := h.RemoveUser(u)
				if err != nil {
					cfg.Log().Errorf("error removing %s from %s", u.Email, h.Alias)
					success = false
					continue
				}
				cfg.Log().Infof("removed %s from %s", u.Email, h.Alias)
			}
		}
	}
	return success
}

func processHostAdded(added []string, cfg Config, h *Host) bool {
	success := true
	for _, group := range added {
		users := cfg.GetUsers(group)
		for _, u := range users {
			if !h.HasUser(u.Email) {
				h.Config.Log().Infof("Adding %s (group %s) to %s", u.Email, group, h.Alias)
				err := h.AddUser(u)
				if err != nil {
					h.Config.Log().Errorf("error adding %s to %s: %v", u.Email, h.Alias, err)
					success = false
					continue
				}
				h.Config.Log().Infof("Added %s to %s (host groups %v)", u.Email, h.Alias, h.Groups)
			}
		}
	}
	return success
}
