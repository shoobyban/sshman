package backend

import "fmt"

type User struct {
	KeyType string   `json:"type"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Key     string   `json:"key"`
	Groups  []string `json:"groups"`
	File    string   `json:"keyfile,omitempty"`
	Config  *Storage `json:"-"`
}

func NewUser(email, keytype, key, name string) *User {
	return &User{
		Email:   email,
		KeyType: keytype,
		Key:     key,
		Name:    name,
	}
}

func (u *User) UpdateGroups(C *Storage, oldgroups []string) error {
	var errors *Errors
	added, removed := updates(oldgroups, u.Groups)
	if u.Config == nil {
		return fmt.Errorf("user has no config")
	}
	u.Config.Log.Infof("added: %v removed: %v", added, removed)
	for _, group := range added {
		hosts := C.getHosts(group)
		for _, h := range hosts {
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
				C.SetHost(h.Alias, h)
			}
		}
	}

	for _, group := range removed {
		hosts := C.getHosts(group)
		for _, h := range hosts {
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
				h.Config.Log.Infof("removed %s from %s %v\n", u.Email, h.Alias, h.Groups)
				C.SetHost(h.Alias, h)
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
