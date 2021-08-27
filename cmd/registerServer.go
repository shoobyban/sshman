package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/apioapp/slog"
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
		conf := readConfig()
		u := findByEmail(conf, args[0])
		if u != nil {
			fmt.Printf("User already exists with this email, overwrite [y/n]: ")
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				slog.Fatalf("error opening stdout %v", err)
			}
			response = strings.ToLower(strings.TrimSpace(response))
			if response != "y" && response != "yes" {
				os.Exit(0)
			}
		}
		registerServer(conf, args...)
	},
}

func init() {
	registerCmd.AddCommand(registerServerCmd)
}

func registerServer(C *config, args ...string) error {
	alias := args[0]
	if _, err := os.Stat(args[3]); os.IsNotExist(err) {
		slog.Fatalf("no such file '%s'", args[3])
		return err
	}
	server := hostentry{
		Host:   args[1],
		User:   args[2],
		Key:    args[3],
		Users:  []string{},
		Groups: args[4:],
	}
	C.Hosts[alias] = server
	slog.Infof("Registering %s to server %s with %s user", alias, args[1], args[1])
	writeConfig(C)
	return nil
}
