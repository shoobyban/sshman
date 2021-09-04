package cmd

import (
	"fmt"

	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

// listServersCmd represents the list servers command
var listServersCmd = &cobra.Command{
	Use:   "servers",
	Short: "List servers",
	Long:  `Lists registered servers`,
	Run: func(_ *cobra.Command, _ []string) {
		conf := backend.ReadConfig(backend.NewSFTP())
		for alias, host := range conf.Hosts {
			fmt.Printf("%-25s\t%-50s\t%v\n", alias, host.Host, host.Groups)
		}
	},
}

func init() {
	listCmd.AddCommand(listServersCmd)
}
