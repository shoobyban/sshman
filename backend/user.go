package backend

import "fmt"

// User holds the user information in the config
type User struct {
	KeyType string   `json:"type"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Key     string   `json:"key"`
	Groups  []string `json:"groups"`
	Hosts   []string `json:"hosts"`
	File    string   `json:"keyfile,omitempty"`
	Config  *Storage `json:"-"`
}

// NewUser creates a new user
func NewUser(email, keytype, key, name string) *User {
	return &User{
		Email:   email,
		KeyType: keytype,
		Key:     key,
		Name:    name,
	}
}

// UpdateGroups updates the user's groups based on old groups
func (u *User) UpdateGroups(C *Storage, oldgroups []string) error {
	added, removed := splitUpdates(oldgroups, u.Groups)
	if u.Config == nil {
		return fmt.Errorf("user has no config")
	}
	u.Config.Log.Infof("added: %v removed: %v", added, removed)

	errors := processUserAdded(added, C, u)

	// are there other groups that keep user on host
	errors = processUserRemoved(removed, C, u, errors)
	C.Write()
	return errors
}

func processUserRemoved(removed []string, C *Storage, u *User, errors *Errors) *Errors {
	for _, group := range removed {
		hosts := C.getHosts(group)
		for _, h := range hosts {

			if h.HasMatchingGroups(u) {
				continue
			}
			if h.HasUser(u.Email) {
				err := h.RemoveUser(u)
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
	return errors
}

func processUserAdded(added []string, C *Storage, u *User) *Errors {
	var errors *Errors
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
				h.Config.Log.Infof("added %s to %s %v\n", u.Email, h.Alias, h.Groups)
				C.SetHost(h.Alias, h)
			}
		}
	}
	return errors
}

// GetGroups is a getter for Groups
func (u *User) GetGroups() []string {
	return u.Groups
}

// SetGroups is a setter for Groups
func (u *User) SetGroups(groups []string) {
	u.Groups = groups
}
