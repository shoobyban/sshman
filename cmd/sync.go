package cmd

import (
	"fmt"
	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync users from hosts",
	Long:  "Fetches all users from authorized_keys on all hosts and updates the local configuration.",
	Run: func(_ *cobra.Command, _ []string) {
		cfg := backend.DefaultConfig()
		cfg.Update()
		fmt.Println("Sync complete")
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
