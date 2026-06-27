package api

import (
	"sort"

	"github.com/shoobyban/sshman/backend"
)

func normalizeUser(user *backend.User) *backend.User {
	if user == nil {
		return nil
	}

	normalized := *user
	if normalized.Groups == nil {
		normalized.Groups = []string{}
	}
	if normalized.Hosts == nil {
		normalized.Hosts = []string{}
	}
	if normalized.Roles == nil {
		normalized.Roles = []string{}
	}

	return &normalized
}

func hostList(hosts map[string]*backend.Host) []*backend.Host {
	result := make([]*backend.Host, 0, len(hosts))
	for _, host := range hosts {
		result = append(result, host)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Alias < result[j].Alias
	})
	return result
}

func userList(users map[string]*backend.User) []*backend.User {
	result := make([]*backend.User, 0, len(users))
	for _, user := range users {
		result = append(result, normalizeUser(user))
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Email < result[j].Email
	})
	return result
}

func groupList(groups map[string]backend.LabelGroup) []backend.LabelGroup {
	result := make([]backend.LabelGroup, 0, len(groups))
	for _, group := range groups {
		result = append(result, group)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Label < result[j].Label
	})
	return result
}

func groupDetails(label string, group *backend.Group) *backend.LabelGroup {
	if group == nil {
		return nil
	}

	result := &backend.LabelGroup{
		Label: label,
		Hosts: make([]string, 0, len(group.Hosts)),
		Users: make([]string, 0, len(group.Users)),
	}

	for _, host := range group.Hosts {
		result.Hosts = append(result.Hosts, host.Alias)
	}
	for _, user := range group.Users {
		result.Users = append(result.Users, user.Email)
	}

	sort.Strings(result.Hosts)
	sort.Strings(result.Users)

	return result
}