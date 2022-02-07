package cmd

import (
	"fmt"

	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

// listAuthCmd represents the listAuth command
var listAuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "List who's on what host",
	Long:  `List who's in authorized_key on what host`,
	Run: func(_ *cobra.Command, _ []string) {
		conf := backend.ReadStorage()
		for alias, host := range conf.Hosts() {
			fmt.Printf("%-25s: %v\n", alias, host.GetUsers())
		}
	},
}

func init() {
	listCmd.AddCommand(listAuthCmd)
}
