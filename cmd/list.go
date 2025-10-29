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
To list hosts:
sshman list hosts
To list groups:
sshman list groups
To list who's on what host:
sshman list auth
`,
}

func init() {
	listCmd.Deprecated = "use 'sshman user list', 'sshman host list', or 'sshman group list' instead"
	rootCmd.AddCommand(listCmd)
}
