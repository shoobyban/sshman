package cmd

import (
	"fmt"
	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "update",
	Short: "Read users into configuration",
	Long:  `Loop through all hosts, download all users from autorized_keys into configuration`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("DEPRECATED: 'sshman update' is deprecated. Use 'sshman sync' instead.")
		cfg := backend.DefaultConfig()
		cfg.Update()
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
}
