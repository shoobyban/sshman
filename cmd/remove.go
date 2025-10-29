package cmd

import (
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove user or host from configuration",
	Long:  `Remove user or host entry (unregister) from configuration.`,
}

func init() {
	removeCmd.Deprecated = "use 'sshman user remove' or 'sshman host remove' instead"
	rootCmd.AddCommand(removeCmd)
}
