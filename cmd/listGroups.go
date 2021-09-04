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
	Long:  `List all groups from users and servers`,
	Run: func(_ *cobra.Command, _ []string) {
		conf := backend.ReadConfig(backend.NewSFTP())
		groups := conf.GetGroups()
		for label, grp := range groups {
			fmt.Printf("%s servers: %v\n%s users: %v\n", label, grp.Servers, label, grp.Users)
		}
	},
}

func init() {
	listCmd.AddCommand(listGroupsCmd)
}
