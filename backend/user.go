package backend

import "log"

type User struct {
	KeyType string   `json:"type"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Key     string   `json:"key"`
	Groups  []string `json:"groups"`
}

func (u *User) UpdateGroups(C *config, oldgroups []string) error {
	added, removed := updates(oldgroups, u.Groups)
	log.Printf("added: %v removed: %v\n", added, removed)
	for _, group := range added {
		servers := C.getServers(group)
		for _, h := range servers {
			h.readUsers()
			if !h.HasUser(u.Email) {
				err := h.AddUser(u)
				if err != nil {
					log.Printf("Error adding %s to %s\n", u.Email, h.Alias)
					continue
				}
				log.Printf("Added %s to %s\n", u.Email, h.Alias)
			}
		}
	}

	for _, group := range removed {
		servers := C.getServers(group)
		for _, h := range servers {
			h.readUsers()
			// are there other groups that keep user on server
			if h.HasMatchingGroups(u) {
				continue
			}
			if h.HasUser(u.Email) {
				err := h.DelUser(u)
				if err != nil {
					log.Printf("Error removing %s from %s\n", u.Email, h.Alias)
					continue
				}
				log.Printf("Removed %s from %s\n", u.Email, h.Alias)
			}
		}
	}
	return nil
}

func (u *User) GetGroups() []string {
	return u.Groups
}

func (u *User) SetGroups(groups []string) {
	u.Groups = groups
}
