package cmd

import (
	"fmt"

	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

// renameHostCmd represents the removeUser command
var renameHostCmd = &cobra.Command{
	Use:   "host",
	Short: "Rename a host (modify host alias)",
	Long:  `Modify a host alias in the configuration, keeping host data intact`,
	Run: func(_ *cobra.Command, args []string) {
		cfg := backend.ReadStorage()
		if len(args) < 2 {
			return
		}
		host := cfg.GetHost(args[0])
		if host != nil {
			cfg.SetHost(args[1], host)
			cfg.DeleteHost(args[0])
			cfg.Write()
			fmt.Printf("Renamed %s to %s\n", args[0], args[1])
		}
	},
}

func init() {
	renameCmd.AddCommand(renameHostCmd)
}
