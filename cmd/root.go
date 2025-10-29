package cmd

import (
	"flag"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/shoobyban/sshman/backend"
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

// Add RBAC-related commands
var rolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "Manage roles and permissions",
	Long:  `Add, remove, or list roles and their permissions`,
}

var assignRoleCmd = &cobra.Command{
	Use:   "assign",
	Short: "Assign a role to a user or host",
	Long:  `Assign a specific role to a user or host to define their permissions`,
}

var listRolesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all roles and their permissions",
	Long:  `Display all roles and the permissions associated with them`,
	Run: func(_ *cobra.Command, _ []string) {
		cfg := backend.DefaultConfig()
		fmt.Println("Roles and their permissions:")
		for name, role := range cfg.Roles() {
			fmt.Printf("Role: %s\nPermissions: %v\n", name, role.Permissions)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initCorba)

	// Add RBAC commands to the root command
	// configure assignRoleCmd flags and handler
	assignRoleCmd.Flags().String("user", "", "Specify the user email to assign the role to")
	assignRoleCmd.Flags().String("host", "", "Specify the host alias to assign the role to")
	assignRoleCmd.Flags().String("role", "", "Specify the role to assign")
	assignRoleCmd.Run = func(cmd *cobra.Command, args []string) {
		userEmail, err := cmd.Flags().GetString("user")
		if err != nil {
			fmt.Printf("Error getting user: %v\n", err)
			return
		}
		hostAlias, err := cmd.Flags().GetString("host")
		if err != nil {
			fmt.Printf("Error getting host: %v\n", err)
			return
		}
		roleName, err := cmd.Flags().GetString("role")
		if err != nil {
			fmt.Printf("Error getting role: %v\n", err)
			return
		}

		if userEmail == "" && hostAlias == "" {
			fmt.Println("Error: --user or --host is required")
			return
		}
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

		if userEmail != "" {
			_, user := cfg.GetUserByEmail(userEmail)
			if user == nil {
				fmt.Printf("User %s not found\n", userEmail)
				return
			}
			user.Roles = append(user.Roles, roleName)
			cfg.UpdateUser(user)
			fmt.Printf("Assigned role %s to user %s\n", roleName, userEmail)
		} else if hostAlias != "" {
			fmt.Println("Error: Hosts cannot have roles. Use groups instead.")
		} else {
			fmt.Println("Error: Either --user or --host must be specified")
		}
	}

	rolesCmd.AddCommand(assignRoleCmd)
	rolesCmd.AddCommand(listRolesCmd)
	rootCmd.AddCommand(rolesCmd)
}

func initCorba() {
	flag.Parse()
}
