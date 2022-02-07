package cmd

import (
	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

// forceupdateCmd represents the forceupdate command
var forceupdateCmd = &cobra.Command{
	Use:   "forceupdate",
	Short: "Write users on hosts",
	Long:  `Loop through all hosts, upload all users to autorized_keys files`,
	Run: func(_ *cobra.Command, _ []string) {
		cfg := backend.ReadStorage()
		var aliases []string
		for a := range cfg.Hosts() {
			aliases = append(aliases, a)
		}
		cfg.Regenerate(aliases...)
	},
}

func init() {
	rootCmd.AddCommand(forceupdateCmd)
}
