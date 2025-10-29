package cmd

import (
	"fmt"

	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Manage groups",
}

var groupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List groups",
	Run: func(_ *cobra.Command, _ []string) {
		cfg := backend.DefaultConfig()
		groups := cfg.GetGroups()
		for label, grp := range groups {
			fmt.Printf("%s hosts: %v\n%s users: %v\n", label, grp.Hosts, label, grp.Users)
		}
	},
}

func init() {
	rootCmd.AddCommand(groupCmd)
	groupCmd.AddCommand(groupListCmd)
}
