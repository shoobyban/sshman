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
		cfg := backend.ReadConfig()
		var aliases []string
		for a := range cfg.Hosts {
			aliases = append(aliases, a)
		}
		cfg.Regenerate(aliases...)
	},
}

func init() {
	rootCmd.AddCommand(forceupdateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// forceupdateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// forceupdateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
