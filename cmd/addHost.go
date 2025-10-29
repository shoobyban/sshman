package cmd

import (
	"fmt"

	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

// addHostCmd represents the register host command
var addHostCmd = &cobra.Command{
	Use:   "host",
	Short: "Add a host",
	Long: `
	To register a host:
sshman register host {alias} {host_address:port} {user} {~/.ssh/working_keyfile.pub} [group1 group2 ...]
For example:
sshman register host google my.google.com:22 myuser ~/.ssh/google.pub deploy hosting google
`,
	Run: func(_ *cobra.Command, args []string) {
		if len(args) < 4 {
			fmt.Println(`Usage: sshman add host <alias> <host:port> <user> <keyfile> [groups...]`)
			return
		}
		cfg := backend.DefaultConfig()
		host := cfg.GetHost(args[0])
		oldgroups := []string{}
		if host != nil {
			fmt.Printf("Host with alias %s already exists. Overwrite? [y/n]: ", args[0])
			exitIfNo()
			oldgroups = host.GetGroups()
		}
		h, err := cfg.PrepareHost(args...)
		if err != nil {
			fmt.Printf("Error preparing host: %v\n", err)
			return
		}
		cfg.AddHost(h, true)
		h.UpdateGroups(cfg, oldgroups)
		cfg.Update(args[0])
	},
}

func init() {
	addCmd.AddCommand(addHostCmd)
}
