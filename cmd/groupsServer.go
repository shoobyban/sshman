package cmd

import (
	"fmt"

	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

// groupsHostCmd represents the user group editing command
var groupsHostCmd = &cobra.Command{
	Use:   "host",
	Short: "Modify group assignments for a host",
	Long: `Modify host groups or remove from all groups:
$ ./sshman groups host hostalias group1 group2
`,
	Run: func(_ *cobra.Command, args []string) {
		cfg := backend.ReadConfig()
		if len(args) < 1 {
			return
		}
		email := args[0]
		groups := args[1:]
		if host, ok := cfg.Hosts[args[0]]; ok {
			oldgroups := host.GetGroups()
			host.SetGroups(groups)
			cfg.Hosts[args[0]] = host
			host.UpdateGroups(cfg, oldgroups)
			cfg.Write()
			fmt.Printf("Groups for %s edited: %v\n", email, host.GetGroups())
		}
	},
}

func init() {
	groupsCmd.AddCommand(groupsHostCmd)
}
