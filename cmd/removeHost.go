package cmd

import (
	"fmt"

	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

// removeHostCmd represents the removeHost command
var removeHostCmd = &cobra.Command{
	Use:   "host",
	Short: "Remove a host from config",
	Long:  `Remove a host from the configuration`,
	Run: func(_ *cobra.Command, args []string) {
		cfg := backend.ReadStorage()
		if len(args) < 1 {
			return
		}
		if cfg.DeleteHost(args[0]) {
			fmt.Printf("deleting %s from configuration\n", args[0])
		}

	},
}

func init() {
	removeCmd.AddCommand(removeHostCmd)
}
