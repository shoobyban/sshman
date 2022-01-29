package cmd

import (
	"fmt"

	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

// listUsersCmd represents the listUsers command
var listUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "List users",
	Long:  `Lists registered users`,
	Run: func(_ *cobra.Command, _ []string) {
		conf := backend.ReadConfig()
		for _, user := range conf.Users() {
			fmt.Printf("%-25s\t%v\n", user.Email, user.GetGroups())
		}
	},
}

func init() {
	listCmd.AddCommand(listUsersCmd)
}
