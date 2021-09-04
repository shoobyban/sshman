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
	for _, group := range added {
		servers := C.getServers(group)
		for _, server := range servers {
			server.readUsers()
			if !server.hasUser(u.Email) {
				log.Printf("Adding %s to %s\n", u.Email, server.Alias)
				err := server.addUser(u)
				if err != nil {
					return err
				}
			}
		}
	}
	for _, group := range removed {
		servers := C.getServers(group)
		for _, server := range servers {
			server.readUsers()
			if server.hasUser(u.Email) {
				err := server.delUser(u)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
