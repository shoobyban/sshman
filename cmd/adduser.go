package cmd

import (
	"fmt"

	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

// addUserCmd represents the addUser command
var addUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Add a user to the configuration",
	Long: `To add a user:
	sshman add user email sshkey.pub {group1 group2 ...}`,
	Run: func(_ *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println(`Usage: sshman add user <email> <sshkey.pub> [group1 group2 ...]`)
			return
		}
		cfg := backend.DefaultConfig()
		_, u := cfg.GetUserByEmail(args[0])
		if u != nil {
			fmt.Printf("User with email %s already exists. Overwrite? [y/n]: ", args[0])
			exitIfNo()
		}
		u, err := cfg.PrepareUser(args[0], args[1], args[2:]...)
		if err != nil {
			fmt.Printf("Error preparing user: %v\n", err)
			return
		}
		cfg.AddUser(u, "")
		cfg.AddUserToHosts(u)
	},
}

func init() {
	addCmd.AddCommand(addUserCmd)
}
