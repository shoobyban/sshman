package cmd

import (
	"fmt"
	"os"

	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

// registerUserCmd represents the registerUser command
var registerUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Register a user",
	Long: `To register a user:
	sshman register user email sshkey.pub {group1 group2 ...}`,
	Run: func(_ *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println(`To register a user:
			sshman register user email sshkey.pub {group1 group2 ...}`)
			os.Exit(0)
		}
		conf := backend.ReadConfig()
		_, u := conf.GetUserByEmail(args[0])
		if u != nil {
			fmt.Printf("User already exists with this email, overwrite [y/n]: ")
			exitIfNo()
		}
		conf.RegisterUser(u.GetGroups(), args...)
	},
}

func init() {
	registerCmd.AddCommand(registerUserCmd)
}
