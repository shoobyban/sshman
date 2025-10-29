package cmd

import (
	"fmt"

	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

// listHostsCmd represents the list hosts command
var listHostsCmd = &cobra.Command{
	Use:   "hosts",
	Short: "List hosts",
	Long:  `Lists registered hosts`,
	Run: func(_ *cobra.Command, _ []string) {
		cfg := backend.DefaultConfig()
		for alias, host := range cfg.Hosts() {
			fmt.Printf("%-25s\t%-50s\t%v\n", alias, host.Host, host.GetGroups())
		}
	},
}

func init() {
	listCmd.AddCommand(listHostsCmd)
}
