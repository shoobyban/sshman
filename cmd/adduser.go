package cmd

import (
	"fmt"
	"os"

	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

// addUserCmd represents the addUser command
var addUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Add a user to the configuration",
	Long: `To add a user:
	sshman add user email sshkey.pub {group1 group2 ...}`,
	Run: func(_ *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println(`To add a user:
			sshman add user email sshkey.pub {group1 group2 ...}`)
			os.Exit(0)
		}
		conf := backend.ReadStorage()
		_, u := conf.GetUserByEmail(args[0])
		if u != nil {
			fmt.Printf("User already exists with this email, overwrite [y/n]: ")
			exitIfNo()
		}
		u, err := conf.PrepareUser(args...)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		conf.AddUser(u)
		conf.AddUserToHosts(u)
	},
}

func init() {
	addCmd.AddCommand(addUserCmd)
}
