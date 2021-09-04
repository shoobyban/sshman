package cmd

import (
	"github.com/spf13/cobra"
)

// registerCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename a user (modify email) or server (modify alias)",
	Long: `Rename a user (modify email) or server (modify alias)
$ ./sshman rename user oldemail@server.com newemail@server.com
$ ./sshman rename server oldalias newalias
`,
}

func init() {
	rootCmd.AddCommand(renameCmd)
}
