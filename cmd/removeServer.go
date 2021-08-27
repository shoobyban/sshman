package cmd

import (
	"github.com/apioapp/slog"
	"github.com/spf13/cobra"
)

// removeServerCmd represents the removeServer command
var removeServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Remove a server from config",
	Long:  `Remove a server from the configuration`,
	Run: func(_ *cobra.Command, args []string) {
		cfg := readConfig()
		if len(args) < 1 {
			return
		}
		for alias := range cfg.Hosts {
			if args[0] == alias {
				slog.Infof("deleting %s from configuration", alias)
				delete(cfg.Hosts, alias)
			}
		}
		writeConfig(cfg)
	},
}

func init() {
	removeCmd.AddCommand(removeServerCmd)
}
