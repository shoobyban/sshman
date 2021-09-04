package cmd

import (
	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "update",
	Short: "Read users into configuration",
	Long:  `Loop through all servers, download all users from autorized_keys into configuration`,
	Run: func(_ *cobra.Command, _ []string) {
		cfg := backend.ReadConfig(backend.NewSFTP())
		cfg.Update()
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
}
