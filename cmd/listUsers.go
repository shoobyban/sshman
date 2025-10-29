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
		fmt.Println("DEPRECATED: 'sshman list users' is deprecated. Use 'sshman user list' instead.")
		cfg := backend.DefaultConfig()
		for _, user := range cfg.Users() {
			fmt.Printf("%-25s\t%v\n", user.Email, user.GetGroups())
		}
	},
}

func init() {
	listCmd.AddCommand(listUsersCmd)
}
