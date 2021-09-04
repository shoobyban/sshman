package cmd

import (
	"github.com/spf13/cobra"
)

// groupsCmd represents the groups command
var groupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "Modify user or server groups",
	Long: `Modify user's groups, or remove groups from user to allow global access:
$ ./sshman groups user email@server.com group1 group2
Modify server groups or remove from all groups:
$ ./sshman groups server serveralias group1 group2
`,
}

func init() {
	rootCmd.AddCommand(groupsCmd)
}
