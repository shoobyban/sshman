package cmd

import (
	"fmt"

	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

// renameServerCmd represents the removeUser command
var renameServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Rename a server (modify host alias)",
	Long:  `Modify a host alias in the configuration, keeping host data intact`,
	Run: func(_ *cobra.Command, args []string) {
		cfg := backend.ReadConfig()
		if len(args) < 2 {
			return
		}
		if host, ok := cfg.Hosts[args[0]]; ok {
			cfg.Hosts[args[1]] = host
			delete(cfg.Hosts, args[0])
			cfg.Write()
			fmt.Printf("Renamed %s to %s\n", args[0], args[1])
		}
	},
}

func init() {
	renameCmd.AddCommand(renameServerCmd)
}
