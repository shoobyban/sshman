package cmd

import (
	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add user",
	Long:  `Add already existing user by email`,
	Run: func(_ *cobra.Command, args []string) {
		conf := backend.ReadConfig()
		for _, email := range args {
			conf.AddUserByEmail(email)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
