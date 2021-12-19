package cmd

import (
	"fmt"

	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

// renameUserCmd represents the removeUser command
var renameUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Rename a user (modify email)",
	Long:  `Modify a user email in the configuration, keeping user data and servers intact`,
	Run: func(_ *cobra.Command, args []string) {
		cfg := backend.ReadConfig()
		if len(args) < 2 {
			return
		}
		key, user := cfg.GetUserByEmail(args[0])
		if user != nil {
			user.Email = args[1]
			cfg.Users[key] = user
			cfg.Write()
			fmt.Printf("Renamed %s to %s\n", args[0], args[1])
		}
	},
}

func init() {
	renameCmd.AddCommand(renameUserCmd)
}
