package backend

import (
	"fmt"
	"log"
	"strings"
)

type Hostentry struct {
	Host     string   `json:"host"`
	User     string   `json:"user"`
	Key      string   `json:"key"`
	Checksum string   `json:"checksum"`
	Users    []string `json:"users"`
	Groups   []string `json:"groups"`
	Alias    string   `json:"-"`
	Config   *config  `json:"-"`
}

func (h *Hostentry) readUsers() error {
	sum, lines, err := h.read()
	if err != nil {
		return err
	}
	var userlist []string
	for _, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) == 3 {
			lsum := checksum(parts[1])
			if _, ok := h.Config.Users[lsum]; !ok {
				h.Config.Users[lsum] = User{
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
	h.Config.Hosts[h.Alias] = *h
	h.Config.Write()
	return nil
}

func (h *Hostentry) read() (string, []string, error) {
	key := h.Key
	if key == "" {
		key = h.Config.Key
	}
	err := h.Config.conn.Connect(key, h.Host, h.User)
	if err != nil {
		return "", nil, fmt.Errorf("error connecting %s: %v", h.Alias, err)
	}
	defer h.Config.conn.Close()
	b, err := h.Config.conn.Read()
	if err != nil {
		return "", nil, fmt.Errorf("error reading authorized keys on %s: %v", h.Alias, err)
	}
	sum := checksum(string(b))
	lines := deleteEmpty(strings.Split(string(b), "\n"))
	return sum, lines, nil
}

func (h *Hostentry) write(lines []string) error {
	if len(lines) == 0 {
		return fmt.Errorf("no keys in new file for server '%s', server would be inaccessible", h.Alias)
	}
	key := h.Key
	if key == "" {
		key = h.Config.Key
	}
	err := h.Config.conn.Connect(key, h.Host, h.User)
	if err != nil {
		return fmt.Errorf("error connecting %s: %v", h.Alias, err)
	}
	defer h.Config.conn.Close()
	return h.Config.conn.Write(strings.Join(lines, "\n") + "\n")
}

func (h *Hostentry) GetUsers() []string {
	return h.Users
}

func (h *Hostentry) hasUser(email string) bool {
	for _, e := range h.Users {
		if e == email {
			return true
		}
	}
	return false
}

func (h *Hostentry) addUser(u *User) error {
	h.Users = append(h.Users, u.Email)
	var lines []string
	for _, email := range h.Users {
		_, userentry := h.Config.GetUserByEmail(email)
		if userentry.Email == email {
			lines = append(lines, userentry.KeyType+" "+userentry.Key+" "+userentry.Name)
		}
	}
	h.Users = append(h.Users, u.Email)
	return h.write(lines)
}

func (h *Hostentry) delUser(u *User) error {
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
			if parts[1] == u.Key {
				found = true
				continue
			}
			newlines = append(newlines, line)
			userlist = append(userlist, h.Config.Users[lsum].Email)
		}
	}

	if found {
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

func (h *Hostentry) UpdateGroups(c *config, oldgroups []string) error {
	added, removed := updates(oldgroups, h.Groups)
	for _, group := range added {
		users := c.getUsers(group)
		for _, u := range users {
			if !h.hasUser(u.Email) {
				err := h.addUser(u)
				if err != nil {
					log.Printf("Error adding %s to %s\n", u.Email, h.Alias)
					continue
				}
				log.Printf("Added %s to %s\n", u.Email, h.Alias)
			}
		}
	}

	for _, group := range removed {
		users := c.getUsers(group)
		for _, u := range users {
			if match(u.Groups, h.Groups) {
				continue
			}
			if h.hasUser(u.Email) {
				err := h.delUser(u)
				if err != nil {
					log.Printf("Error removing %s from %s\n", u.Email, h.Alias)
					continue
				}
				log.Printf("Removed %s from %s\n", u.Email, h.Alias)
			}
		}
	}
	c.Hosts[h.Alias] = *h
	c.Write()
	return nil
}
