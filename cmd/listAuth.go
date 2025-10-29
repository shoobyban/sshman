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
		fmt.Println("DEPRECATED: 'sshman list auth' is deprecated. Use 'sshman host list' and 'sshman user list' for relevant info.")
		cfg := backend.DefaultConfig()
		for alias, host := range cfg.Hosts() {
			fmt.Printf("%-25s: %v\n", alias, host.GetUsers())
		}
	},
}

func init() {
	listCmd.AddCommand(listAuthCmd)
}
