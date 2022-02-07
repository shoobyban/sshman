package cmd

import (
	"fmt"

	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

// listGroupsCmd represents the listGroups command
var listGroupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List all groups",
	Long:  `List all groups from users and hosts`,
	Run: func(_ *cobra.Command, _ []string) {
		conf := backend.ReadStorage()
		groups := conf.GetGroups()
		for label, grp := range groups {
			fmt.Printf("%s hosts: %v\n%s users: %v\n", label, grp.Hosts, label, grp.Users)
		}
	},
}

func init() {
	listCmd.AddCommand(listGroupsCmd)
}
