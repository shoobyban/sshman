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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initCorba)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initCorba() {
	flag.Parse()
}
