package cmd

import (
	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a user or server",
	Long: `To register a user:
sshman register user {email} {sshkey.pub} [group1 group2 ...]
For example:
sshman register user email@test.com ~/.ssh/user1.pub production-team staging-servers

To register a server:
sshman register server {alias} {server_address:port} {user} {~/.ssh/working_keyfile.pub} [group1 group2 ...]
For example:
sshman register server google my.google.com:22 myuser ~/.ssh/google.pub deploy hosting google
`,
}

func init() {
	rootCmd.AddCommand(registerCmd)
}
