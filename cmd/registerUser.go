package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

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
		conf := backend.ReadConfig(backend.NewSFTP())
		var oldgroups []string
		_, u := conf.GetUserByEmail(args[0])
		if u != nil {
			fmt.Printf("User already exists with this email, overwrite [y/n]: ")
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("Error: error opening stdout %v\n", err)
			}
			response = strings.ToLower(strings.TrimSpace(response))
			if response != "y" && response != "yes" {
				os.Exit(0)
			}
		}
		conf.RegisterUser(oldgroups, args...)
	},
}

func init() {
	registerCmd.AddCommand(registerUserCmd)
}
