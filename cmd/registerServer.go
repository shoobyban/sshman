package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

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
		conf := backend.ReadConfig(backend.NewSFTP())
		_, exists := conf.Hosts[args[0]]
		if exists {
			fmt.Printf("Host already exists with this alias, overwrite [y/n]: ")
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalf("error opening stdout %v\n", err)
			}
			response = strings.ToLower(strings.TrimSpace(response))
			if response != "y" && response != "yes" {
				os.Exit(0)
			}
		}
		conf.RegisterServer(args...)
		conf.Update(args[0])
	},
}

func init() {
	registerCmd.AddCommand(registerServerCmd)
}
