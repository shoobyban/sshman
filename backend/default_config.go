package backend

func DefaultConfig() *Data {
	s := &JsonStorage{
		Path: ".ssh/.sshman",
	}
	data := ReadData(s)

	// Initialize default roles and permissions
	if _, ok := data.roles["admin"]; !ok {
		data.roles["admin"] = &Role{Name: "admin", Permissions: []string{"add_user", "remove_user", "add_host", "remove_host", "assign_role"}}
	}
	if _, ok := data.roles["user"]; !ok {
		data.roles["user"] = &Role{Name: "user", Permissions: []string{"view_hosts", "view_users"}}
	}

	return data
}
