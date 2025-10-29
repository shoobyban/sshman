package cmd

import (
	"fmt"

	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

// groupsUserCmd represents the user group editing command
var groupsUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Modify group assignments of a user",
	Long: `Modify user's groups, or remove groups from user to allow global access:
$ ./sshman groups user email@host.com group1 group2
`,
	Run: func(_ *cobra.Command, args []string) {
		fmt.Println("DEPRECATED: 'sshman groups user' is deprecated. Use 'sshman user groups' instead.")
		cfg := backend.DefaultConfig()
		if len(args) < 1 {
			return
		}
		email := args[0]
		groups := args[1:]
		_, user := cfg.GetUserByEmail(email)
		if user != nil {
			oldGroups := user.GetGroups()
			user.SetGroups(groups)
			cfg.UpdateUser(user)
			user.UpdateGroups(cfg, oldGroups)
			cfg.Write()
			fmt.Printf("Groups for %s edited: %v\n", email, groups)
		}
	},
}

func init() {
	groupsCmd.AddCommand(groupsUserCmd)
}
