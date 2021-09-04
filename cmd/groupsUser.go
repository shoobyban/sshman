package cmd

import (
	"log"

	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

// groupsUserCmd represents the user group editing command
var groupsUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Modify group assignments of a user",
	Long: `Modify user's groups, or remove groups from user to allow global access:
$ ./sshman groups user email@server.com group1 group2
`,
	Run: func(_ *cobra.Command, args []string) {
		cfg := backend.ReadConfig(backend.NewSFTP())
		if len(args) < 1 {
			return
		}
		email := args[0]
		groups := args[1:]
		key, user := cfg.GetUserByEmail(email)
		if user != nil {
			oldgroups := user.Groups
			user.Groups = groups
			cfg.Users[key] = *user
			cfg.Write()
			log.Printf("Groups for %s edited: %v\n", email, groups)
			user.UpdateGroups(cfg, oldgroups)
		}
	},
}

func init() {
	groupsCmd.AddCommand(groupsUserCmd)
}
