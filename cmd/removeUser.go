package cmd

import (
	"fmt"

	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

// removeUserCmd represents the removeUser command
var removeUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Remove user from config",
	Long:  `Remove a user by email from the configuration`,
	Run: func(_ *cobra.Command, args []string) {
		cfg := backend.ReadConfig()
		if len(args) < 1 {
			return
		}
		if cfg.UnregisterUser(args[0]) {
			fmt.Printf("Deleted %s from configuration\n", args[0])
		}
	},
}

func init() {
	removeCmd.AddCommand(removeUserCmd)
}
