package backend

type User struct {
	KeyType string   `json:"type"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Key     string   `json:"key"`
	Groups  []string `json:"groups"`
}

func (u *User) updateGroups(C *config, oldgroups, newgroups []string) error {
	changes := updates(oldgroups, newgroups)
	for _, group := range changes {
		servers := C.getServers(group)
		for _, server := range servers {
			server.readUsers()
			if !server.hasUser(u.Email) {
				err := server.addUser(u)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
