package backend

import (
	"fmt"
	"log"
	"strings"
)

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

func (h *Host) ReadUsers() error {
	sum, lines, err := h.read()
	if err != nil {
		return err
	}
	var userlist []string
	for _, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) == 3 {
			lsum := checksum(parts[1])
			if _, exists := h.Config.Users[lsum]; !exists {
				h.Config.Users[lsum] = &User{
					KeyType: parts[0],
					Key:     parts[1],
					Name:    parts[2],
					Email:   parts[2] + "@" + h.Alias,
				}
			}
			userlist = append(userlist, h.Config.Users[lsum].Email)
		}
	}
	h.Checksum = sum
	h.Users = userlist
	h.Config.Write()
	return nil
}

func (h *Host) read() (string, []string, error) {
	key := h.Key
	if key == "" {
		key = h.Config.Key
	}
	err := h.Config.Conn.Connect(key, h.Host, h.User)
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
	key := h.Key
	if key == "" {
		key = h.Config.Key
	}
	err := h.Config.Conn.Connect(key, h.Host, h.User)
	if err != nil {
		return fmt.Errorf("error connecting %s: %v", h.Alias, err)
	}
	defer h.Config.Conn.Close()
	return h.Config.Conn.Write(strings.Join(lines, "\n") + "\n")
}

func (h *Host) GetUsers() []string {
	return h.Users
}

func (h *Host) GetGroups() []string {
	return h.Groups
}

func (h *Host) SetGroups(groups []string) {
	h.Groups = groups
}

func (h *Host) HasMatchingGroups(user *User) bool {
	return match(h.GetGroups(), user.GetGroups())
}

func (h *Host) HasUser(email string) bool {
	for _, e := range h.Users {
		if e == email {
			return true
		}
	}
	return false
}

func (h *Host) AddUser(u *User) error {
	h.Users = append(h.Users, u.Email)
	var lines []string
	for _, email := range h.Users {
		_, userentry := h.Config.GetUserByEmail(email)
		if userentry == nil {
			// Shall we add a warning here?
			continue
		}
		if userentry.Email == email {
			lines = append(lines, userentry.KeyType+" "+userentry.Key+" "+userentry.Name)
		}
	}
	return h.write(lines)
}

func (h *Host) DelUser(u *User) error {
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
			log.Printf("Del Checksum: %s, key: %s, userkey: %s", lsum, parts[1], u.Key)
			if parts[1] == u.Key {
				found = true
				continue
			}
			if user, ok := h.Config.Users[lsum]; ok {
				newlines = append(newlines, line)
				userlist = append(userlist, user.Email)
			}
		}
	}

	if found {
		log.Printf("found, deleting %v", u)
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

func (h *Host) UpdateGroups(c *Storage, oldgroups []string) bool {
	success := true
	added, removed := updates(oldgroups, h.Groups)
	fmt.Printf("added: %v removed: %v\n", added, removed)
	for _, group := range added {
		users := c.GetUsers(group)
		for _, u := range users {
			fmt.Printf("User %s from group %s\n", u.Email, group)
			if !h.HasUser(u.Email) {
				err := h.AddUser(u)
				if err != nil {
					fmt.Printf("Error adding %s to %s\n", u.Email, h.Alias)
					success = false
					continue
				}
				fmt.Printf("Added %s to %s %v\n", u.Email, h.Alias, h.Groups)
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
					fmt.Printf("Error removing %s from %s\n", u.Email, h.Alias)
					success = false
					continue
				}
				fmt.Printf("Removed %s from %s\n", u.Email, h.Alias)
			}
		}
	}
	c.Hosts[h.Alias] = h
	c.Write()
	return success
}
