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
		if len(line) == 0 {
			continue
		}
		parts := strings.Split(line, " ")
		if len(parts) != 3 {
			log.Printf("Error: Not good line: '%s'\n", line)
		}
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
		log.Printf("Error: error connecting %s: %v\n", h.Alias, err)
		return "", nil, err
	}
	defer h.Config.conn.Close()
	b, err := h.Config.conn.Read()
	if err != nil {
		log.Printf("Error: error reading authorized keys on %s: %v\n", h.Alias, err)
		return "", nil, err
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
		log.Printf("Error: error connecting %s: %v\n", h.Alias, err)
		return err
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
	return h.write(lines)
}

func (h *Hostentry) delUser(u *User) error {
	sum, lines, err := h.read()
	if err != nil {
		log.Printf("Error reading from host %s %v", h.Alias, err)
		return err
	}
	userlist := []string{}
	found := false
	newlines := []string{}
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		parts := strings.Split(line, " ")
		if len(parts) != 3 {
			log.Printf("Error: Not good line: '%s'\n", line)
		}
		lsum := checksum(parts[1])
		if _, ok := h.Config.Users[lsum]; !ok {
			delete(h.Config.Users, lsum)
		}
		if parts[1] == u.Key {
			found = true
			continue
		}
		newlines = append(newlines, line)
		userlist = append(userlist, h.Config.Users[lsum].Email)
	}

	if found {
		newlines = deleteEmpty(newlines)
		err = h.write(newlines)
		if err != nil {
			log.Printf("Error: error writing %s: %v\n", h.Alias, err)
			return err
		}
		log.Printf("Removed %s from %s\n", u.Email, h.Alias)
	}
	h.Checksum = sum
	h.Users = userlist
	return nil
}
