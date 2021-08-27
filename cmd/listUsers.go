package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listUsersCmd represents the listUsers command
var listUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "List users",
	Long:  `Lists registered users`,
	Run: func(_ *cobra.Command, _ []string) {
		conf := readConfig()
		for _, user := range conf.Users {
			fmt.Printf("%-25s\t%v\n", user.Email, user.Groups)
		}
	},
}

func init() {
	listCmd.AddCommand(listUsersCmd)
}
