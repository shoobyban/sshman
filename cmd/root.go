package cmd

import (
	"flag"
	"fmt"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sshman",
	Short: "SSH Key manager",
	Long:  `Deploy keys to remote hosts`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("Run sshman --help for more help")
	},
}

// RBAC commands are defined in cmd/role.go (single `role` command)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initCorba)

	// Global persistent flags
	rootCmd.PersistentFlags().StringP("config", "c", "", "Path to configuration file")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")

	// `role` command implementation lives in cmd/role.go and is registered there.

	// Deprecate old top-level commands (they have been replaced by resource-oriented commands)
	// Mark existing commands if present in this package by name. We cannot remove them immediately to preserve backward compatibility.
	// The source files still register commands like addCmd, removeCmd, renameCmd, readCmd, groupsCmd, listCmd, webCmd, delCmd.
	// Set the Deprecated field where possible.
	// Previously these top-level convenience commands were marked as Deprecated
	// to guide users toward resource-oriented commands (user/host/group subcommands).
	// We intentionally clear the Cobra Deprecated fields so that the CLI help output
	// no longer shows them as deprecated. The old commands still exist for
	// backward compatibility but there is no runtime deprecation banner.
}

func initCorba() {
	flag.Parse()
}
