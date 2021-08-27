package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type Group struct {
	Users   []string
	Servers []string
}

// listGroupsCmd represents the listGroups command
var listGroupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List all groups",
	Long:  `List all groups from users and servers`,
	Run: func(_ *cobra.Command, _ []string) {
		conf := readConfig()
		groups := map[string]Group{}
		for alias, host := range conf.Hosts {
			for _, group := range host.Groups {
				if v, ok := groups[group]; ok {
					v.Servers = append(v.Servers, alias)
					groups[group] = v
				} else {
					groups[group] = Group{Servers: []string{alias}}
				}
			}
		}
		for _, user := range conf.Users {
			for _, group := range user.Groups {
				if _, ok := groups[group]; ok {
					g := groups[group]
					g.Users = append(g.Users, user.Email)
					groups[group] = g
				} else {
					groups[group] = Group{Users: []string{user.Email}}
				}
			}
		}
		for label, grp := range groups {
			fmt.Printf("%s servers: %v\n%s users: %v\n", label, grp.Servers, label, grp.Users)
		}

	},
}

func init() {
	listCmd.AddCommand(listGroupsCmd)
}
