package cmd

import (
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a user or host to the configuration",
	Long: `To add a user:
sshman add user {email} {sshkey.pub} [group1 group2 ...]
For example:
sshman add user email@test.com ~/.ssh/user1.pub production-team staging-hosts

To add a host:
sshman add host {alias} {host_address:port} {user} {~/.ssh/working_keyfile.pub} [group1 group2 ...]
For example:
sshman add host google my.google.com:22 myuser ~/.ssh/google.pub deploy hosting google
`,
}

func init() {
	rootCmd.AddCommand(addCmd)
}
