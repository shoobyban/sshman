package cmd

type hostentry struct {
	Host     string   `json:"host"`
	User     string   `json:"user"`
	Key      string   `json:"key"`
	Checksum string   `json:"checksum"`
	Users    []string `json:"users"`
	Groups   []string `json:"groups"`
}

type user struct {
	KeyType string   `json:"type"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Key     string   `json:"key"`
	Groups  []string `json:"groups"`
}

type config struct {
	Key   string               `json:"key"`
	Hosts map[string]hostentry `json:"hosts"`
	Users map[string]user      `json:"users"`
}
