package cmd

import (
	"github.com/spf13/cobra"
)

// groupsCmd represents the groups command
var groupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "Modify user or host groups",
	Long: `Modify user's groups, or remove groups from user to allow global access:
$ ./sshman groups user email@host.com group1 group2
Modify host groups or remove from all groups:
$ ./sshman groups host hostalias group1 group2
`,
}

func init() {
	groupsCmd.Deprecated = "use 'sshman user groups' or 'sshman host groups' instead"
	rootCmd.AddCommand(groupsCmd)
}
