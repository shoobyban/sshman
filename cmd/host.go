package cmd

import (
	"fmt"

	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

var hostCmd = &cobra.Command{
	Use:   "host",
	Short: "Manage hosts",
}

var hostAddCmd = &cobra.Command{
	Use:   "add <alias> <host:port> <user> <keyfile>",
	Short: "Add a host",
	Args:  cobra.MinimumNArgs(4),
	Run: func(_ *cobra.Command, args []string) {
		cfg := backend.DefaultConfig()
		alias := args[0]
		if cfg.GetHost(alias) != nil {
			fmt.Printf("Host with alias %s already exists. Overwrite? [y/n]: ", alias)
			exitIfNo()
		}
		h, err := cfg.PrepareHost(args...)
		if err != nil {
			fmt.Printf("Error preparing host: %v\n", err)
			return
		}
		cfg.AddHost(h, true)
		cfg.Update(alias)
		fmt.Printf("Added host %s\n", alias)
	},
}

var hostRemoveCmd = &cobra.Command{
	Use:   "remove <alias>",
	Short: "Remove a host",
	Args:  cobra.MinimumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		cfg := backend.DefaultConfig()
		for _, alias := range args {
			if cfg.DeleteHost(alias) {
				fmt.Printf("Deleted host %s\n", alias)
			} else {
				fmt.Printf("Host %s not found\n", alias)
			}
		}
	},
}

var hostRenameCmd = &cobra.Command{
	Use:   "rename <old_alias> <new_alias>",
	Short: "Rename a host",
	Args:  cobra.ExactArgs(2),
	Run: func(_ *cobra.Command, args []string) {
		cfg := backend.DefaultConfig()
		old := args[0]
		new := args[1]
		h := cfg.GetHost(old)
		if h == nil {
			fmt.Printf("Host %s not found\n", old)
			return
		}
		cfg.SetHost(new, h)
		cfg.DeleteHost(old)
		cfg.Write()
		fmt.Printf("Renamed %s to %s\n", old, new)
	},
}

var hostListCmd = &cobra.Command{
	Use:   "list",
	Short: "List hosts",
	Run: func(_ *cobra.Command, _ []string) {
		cfg := backend.DefaultConfig()
		for alias, host := range cfg.Hosts() {
			fmt.Printf("%-25s\t%-50s\t%v\n", alias, host.Host, host.GetGroups())
		}
	},
}

var hostGroupsCmd = &cobra.Command{
	Use:   "groups <alias> [groups...]",
	Short: "Set groups for a host",
	Args:  cobra.MinimumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		cfg := backend.DefaultConfig()
		alias := args[0]
		groups := args[1:]
		h := cfg.GetHost(alias)
		if h == nil {
			fmt.Printf("Host %s not found\n", alias)
			return
		}
		old := h.GetGroups()
		h.SetGroups(groups)
		cfg.SetHost(alias, h)
		h.UpdateGroups(cfg, old)
		cfg.Write()
		fmt.Printf("Groups for %s updated: %v\n", alias, groups)
	},
}

func init() {
	rootCmd.AddCommand(hostCmd)
	hostCmd.AddCommand(hostAddCmd)
	hostCmd.AddCommand(hostRemoveCmd)
	hostCmd.AddCommand(hostRenameCmd)
	hostCmd.AddCommand(hostListCmd)
	hostCmd.AddCommand(hostGroupsCmd)
	hostAddCmd.Flags().StringArrayP("group", "g", []string{}, "Group(s) to add the host to (repeatable)")
}
