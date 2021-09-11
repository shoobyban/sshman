package cmd

import (
	"fmt"

	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

// removeServerCmd represents the removeServer command
var removeServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Remove a server from config",
	Long:  `Remove a server from the configuration`,
	Run: func(_ *cobra.Command, args []string) {
		cfg := backend.ReadConfig()
		if len(args) < 1 {
			return
		}
		if cfg.UnregisterServer(args[0]) {
			fmt.Printf("deleting %s from configuration\n", args[0])
		}

	},
}

func init() {
	removeCmd.AddCommand(removeServerCmd)
}
