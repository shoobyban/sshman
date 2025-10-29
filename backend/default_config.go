package backend

func DefaultConfig() *Data {
	s := &JsonStorage{
		Path: ".ssh/.sshman",
	}
	data := NewData(s)

	// Initialize default roles and permissions
	data.roles["admin"] = &Role{Name: "admin", Permissions: []string{"add_user", "remove_user", "add_host", "remove_host", "assign_role"}}
	data.roles["user"] = &Role{Name: "user", Permissions: []string{"view_hosts", "view_users"}}

	return data
}
