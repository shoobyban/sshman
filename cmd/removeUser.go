package cmd

import (
	"github.com/apioapp/slog"
	"github.com/spf13/cobra"
)

// removeUserCmd represents the removeUser command
var removeUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Remove user from config",
	Long:  `Remove a user by email from the configuration`,
	Run: func(_ *cobra.Command, args []string) {
		cfg := readConfig()
		if len(args) < 1 {
			return
		}
		for id, user := range cfg.Users {
			if args[0] == user.Email {
				slog.Infof("deleting %s from configuration", user.Email)
				delete(cfg.Users, id)
			}
		}
		writeConfig(cfg)
	},
}

func init() {
	removeCmd.AddCommand(removeUserCmd)
}
