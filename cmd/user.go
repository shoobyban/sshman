package cmd

import (
	"fmt"

	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users",
}

var userAddCmd = &cobra.Command{
	Use:   "add <email> <sshkey.pub>",
	Short: "Add a user",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := backend.DefaultConfig()
		email := args[0]
		keyfile := args[1]
		groups, _ := cmd.Flags().GetStringArray("group")
		_, u := cfg.GetUserByEmail(email)
		if u != nil {
			fmt.Printf("User with email %s already exists. Overwrite? [y/n]: ", email)
			exitIfNo()
		}
		user, err := cfg.PrepareUser(email, keyfile, groups...)
		if err != nil {
			fmt.Printf("Error preparing user: %v\n", err)
			return
		}
		cfg.AddUser(user, "")
		cfg.AddUserToHosts(user)
		fmt.Printf("Added user %s\n", email)
	},
}

var userRemoveCmd = &cobra.Command{
	Use:   "remove <email>",
	Short: "Remove a user",
	Args:  cobra.MinimumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		cfg := backend.DefaultConfig()
		for _, email := range args {
			_, u := cfg.GetUserByEmail(email)
			if u != nil {
				cfg.RemoveUserFromHosts(u)
				fmt.Printf("Removed user %s\n", email)
			} else {
				fmt.Printf("User %s not found\n", email)
			}
		}
	},
}

var userRenameCmd = &cobra.Command{
	Use:   "rename <old_email> <new_email>",
	Short: "Rename a user",
	Args:  cobra.ExactArgs(2),
	Run: func(_ *cobra.Command, args []string) {
		cfg := backend.DefaultConfig()
		old := args[0]
		new := args[1]
		_, u := cfg.GetUserByEmail(old)
		if u == nil {
			fmt.Printf("User %s not found\n", old)
			return
		}
		u.Email = new
		cfg.UpdateUser(u)
		cfg.Write()
		fmt.Printf("Renamed %s to %s\n", old, new)
	},
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List users",
	Run: func(_ *cobra.Command, _ []string) {
		cfg := backend.DefaultConfig()
		for _, u := range cfg.Users() {
			fmt.Printf("%-25s\t%v\n", u.Email, u.GetGroups())
		}
	},
}

var userGroupsCmd = &cobra.Command{
	Use:   "groups <email> [groups...]",
	Short: "Set groups for a user",
	Args:  cobra.MinimumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		cfg := backend.DefaultConfig()
		email := args[0]
		groups := args[1:]
		_, u := cfg.GetUserByEmail(email)
		if u == nil {
			fmt.Printf("User %s not found\n", email)
			return
		}
		old := u.GetGroups()
		u.SetGroups(groups)
		cfg.UpdateUser(u)
		u.UpdateGroups(cfg, old)
		cfg.Write()
		fmt.Printf("Groups for %s updated: %v\n", email, groups)
	},
}

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(userAddCmd)
	userCmd.AddCommand(userRemoveCmd)
	userCmd.AddCommand(userRenameCmd)
	userCmd.AddCommand(userListCmd)
	userCmd.AddCommand(userGroupsCmd)
	userAddCmd.Flags().StringArrayP("group", "g", []string{}, "Group(s) to add the user to (repeatable)")
}
