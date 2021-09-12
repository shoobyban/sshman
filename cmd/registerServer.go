package cmd

import (
	"fmt"
	"os"

	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

// registerServerCmd represents the register server command
var registerServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Register a server",
	Long: `
	To register a server:
sshman register server {alias} {server_address:port} {user} {~/.ssh/working_keyfile.pub} [group1 group2 ...]
For example:
sshman register server google my.google.com:22 myuser ~/.ssh/google.pub deploy hosting google
`,
	Run: func(_ *cobra.Command, args []string) {
		if len(args) < 4 {
			fmt.Print(`To register a server:
sshman register server {alias} {server_address:port} {user} {~/.ssh/working_keyfile.pub} [group1 group2 ...]
For example:
sshman register server google my.google.com:22 myuser ~/.ssh/google.pub deploy hosting google
`)
			os.Exit(0)
		}
		conf := backend.ReadConfig()
		_, exists := conf.Hosts[args[0]]
		if exists {
			fmt.Printf("Host already exists with this alias, overwrite [y/n]: ")
			exitIfNo()
		}
		conf.RegisterServer([]string{}, args...)
		conf.Update(args[0])
	},
}

func init() {
	registerCmd.AddCommand(registerServerCmd)
}
