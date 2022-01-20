package backend

import "fmt"

type User struct {
	KeyType string   `json:"type"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Key     string   `json:"key"`
	Groups  []string `json:"groups"`
}

func (u *User) UpdateGroups(C *Storage, oldgroups []string) error {
	var errors *Errors
	added, removed := updates(oldgroups, u.Groups)
	fmt.Printf("added: %v removed: %v\n", added, removed)
	for _, group := range added {
		hosts := C.getHosts(group)
		for _, h := range hosts {
			h.ReadUsers()
			if !h.HasUser(u.Email) {
				err := h.AddUser(u)
				if err != nil {
					if errors == nil {
						errors = &Errors{}
					}
					errors.Add("Error adding %s to %s: %v", u.Email, h.Alias, err)
					continue
				}
				//fmt.Printf("Added %s to %s %v\n", u.Email, h.Alias, h.Groups)
				C.Hosts[h.Alias] = h
			}
		}
	}

	for _, group := range removed {
		hosts := C.getHosts(group)
		for _, h := range hosts {
			h.ReadUsers()
			// are there other groups that keep user on host
			if h.HasMatchingGroups(u) {
				continue
			}
			if h.HasUser(u.Email) {
				err := h.DelUser(u)
				if err != nil {
					if errors == nil {
						errors = &Errors{}
					}
					errors.Add("Error removing %s from %s", u.Email, h.Alias)
					continue
				}
				fmt.Printf("Removed %s from %s %v\n", u.Email, h.Alias, h.Groups)
				C.Hosts[h.Alias] = h
			}
		}
	}
	C.Write()
	return errors
}

func (u *User) GetGroups() []string {
	return u.Groups
}

func (u *User) SetGroups(groups []string) {
	u.Groups = groups
}
