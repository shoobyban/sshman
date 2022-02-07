package cmd

import (
	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "update",
	Short: "Read users into configuration",
	Long:  `Loop through all hosts, download all users from autorized_keys into configuration`,
	Run: func(_ *cobra.Command, _ []string) {
		cfg := backend.ReadStorage()
		cfg.Update()
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
}
