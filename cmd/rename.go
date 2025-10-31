package cmd

import (
	"github.com/spf13/cobra"
)

// registerCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename a user (modify email) or host (modify alias)",
	Long: `Rename a user (modify email) or host (modify alias)
$ ./sshman rename user oldemail@host.com newemail@host.com
$ ./sshman rename host oldalias newalias
`,
}

func init() {
	rootCmd.AddCommand(renameCmd)
}
