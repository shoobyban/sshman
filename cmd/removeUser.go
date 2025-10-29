package cmd

import (
	"fmt"

	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

// delCmd represents the del command
var delCmd = &cobra.Command{
	Use:   "del",
	Short: "Delete user by email",
	Long:  `Check all hosts and delete user with given email`,
	Run: func(_ *cobra.Command, args []string) {
		cfg := backend.DefaultConfig()
		for _, email := range args {
			_, u := cfg.GetUserByEmail(email)
			if u != nil {
				cfg.RemoveUserFromHosts(u)
			} else {
				fmt.Printf("No such user\n")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(delCmd)
}
