package cmd

import (
	"fmt"

	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

// groupsServerCmd represents the user group editing command
var groupsServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Modify group assignments for a server",
	Long: `Modify server groups or remove from all groups:
$ ./sshman groups server serveralias group1 group2
`,
	Run: func(_ *cobra.Command, args []string) {
		cfg := backend.ReadConfig()
		if len(args) < 1 {
			return
		}
		email := args[0]
		groups := args[1:]
		if host, ok := cfg.Hosts[args[0]]; ok {
			host.SetGroups(groups)
			cfg.Hosts[args[0]] = host
			cfg.Write()
			fmt.Printf("Groups for %s edited: %v\n", email, host.GetGroups())
		}
	},
}

func init() {
	groupsCmd.AddCommand(groupsServerCmd)
}
