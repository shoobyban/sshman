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
		conf := backend.ReadConfig()
		for _, email := range args {
			_, u := conf.GetUserByEmail(email)
			if u != nil {
				conf.DelUserFromHosts(u)
			} else {
				fmt.Printf("No such user\n")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(delCmd)
}
