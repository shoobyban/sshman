package cmd

import (
	"fmt"
	"os"

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
			fmt.Print(`To register a host:
sshman register host {alias} {host_address:port} {user} {~/.ssh/working_keyfile.pub} [group1 group2 ...]
For example:
sshman register host google my.google.com:22 myuser ~/.ssh/google.pub deploy hosting google
`)
			os.Exit(0)
		}
		conf := backend.ReadConfig()
		host := conf.GetHost(args[0])
		oldgroups := []string{}
		if host != nil {
			fmt.Printf("Host already exists with this alias, overwrite [y/n]: ")
			exitIfNo()
			oldgroups = host.GetGroups()
		}
		h, err := conf.PrepareHost(args...)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		conf.AddHost(h)
		host.UpdateGroups(conf, oldgroups)
		conf.Update(args[0])
	},
}

func init() {
	addCmd.AddCommand(addHostCmd)
}
