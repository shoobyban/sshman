package cmd

import (
	"log"

	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add user",
	Long:  `Add already existing user by email`,
	Run: func(_ *cobra.Command, args []string) {
		conf := backend.ReadConfig(backend.NewSFTP())
		for _, email := range args {
			_, u := conf.GetUserByEmail(email)
			if u != nil {
				conf.AddUserToHosts(u)
			} else {
				log.Printf("No such user, register user first\n")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
