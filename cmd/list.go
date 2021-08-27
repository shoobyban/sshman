package cmd

import (
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List resources",
	Long: `To list available users:
sshman list users
To list servers:
sshman list servers
To list groups:
sshman list groups
To list who's on what server:
sshman list auth
`,
}

func init() {
	rootCmd.AddCommand(listCmd)
}
