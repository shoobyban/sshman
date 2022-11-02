package backend

// Group is used to store group information (connecting hosts and users)
type Group struct {
	Name  string
	Size  int
	Hosts []*Host
	Users []*User
}

func (g Group) HasHost(alias string) bool {
	for _, host := range g.Hosts {
		if host.Alias == alias {
			return true
		}
	}
	return false
}
