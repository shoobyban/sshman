package cmd

import (
	"fmt"

	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

var roleCmd = &cobra.Command{
	Use:   "role",
	Short: "Manage roles",
}

var roleAssignCmd = &cobra.Command{
	Use:   "assign",
	Short: "Assign role to user or host",
	Run: func(cmd *cobra.Command, _ []string) {
		userEmail, _ := cmd.Flags().GetString("user")
		hostAlias, _ := cmd.Flags().GetString("host")
		roleName, _ := cmd.Flags().GetString("role")
		if roleName == "" {
			fmt.Println("Error: --role is required")
			return
		}
		cfg := backend.DefaultConfig()
		roles := cfg.Roles()
		if _, exists := roles[roleName]; !exists {
			fmt.Printf("Role %s does not exist\n", roleName)
			return
		}
		if userEmail == "" && hostAlias == "" {
			fmt.Println("Error: --user or --host is required")
			return
		}
		if userEmail != "" {
			_, u := cfg.GetUserByEmail(userEmail)
			if u == nil {
				fmt.Printf("User %s not found\n", userEmail)
				return
			}
			// add role to user and update
			u.Roles = append(u.Roles, roleName)
			if err := cfg.UpdateUser(u); err != nil {
				fmt.Printf("Error updating user with new role: %v\n", err)
				return
			}
			fmt.Printf("Assigned role %s to user %s\n", roleName, userEmail)
		}
		if hostAlias != "" {
			// backend currently doesn't support assigning roles to hosts
			fmt.Println("Error: Hosts cannot have roles in current backend. Use groups to control host access.")
		}
	},
}

var roleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List roles",
	Run: func(_ *cobra.Command, _ []string) {
		cfg := backend.DefaultConfig()
		fmt.Println("Roles and their permissions:")
		for name, role := range cfg.Roles() {
			fmt.Printf("Role: %s\nPermissions: %v\n", name, role.Permissions)
		}
	},
}

func init() {
	rootCmd.AddCommand(roleCmd)
	roleCmd.AddCommand(roleAssignCmd)
	roleCmd.AddCommand(roleListCmd)
	roleAssignCmd.Flags().String("user", "", "Specify the user email to assign the role to")
	roleAssignCmd.Flags().String("host", "", "Specify the host alias to assign the role to")
	roleAssignCmd.Flags().String("role", "", "Specify the role to assign")
}
