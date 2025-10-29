package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

const commandTree = `sshman
├── user
│   ├── add <email> <sshkey.pub> [flags]
│   ├── remove <email>
│   ├── rename <old_email> <new_email>
│   ├── list
│   └── groups <email> [groups...]
├── host
│   ├── add <alias> <host:port> <user> <keyfile> [flags]
│   ├── remove <alias>
│   ├── rename <old_alias> <new_alias>
│   ├── list
│   └── groups <alias> [groups...]
├── group
│   └── list
├── role
│   ├── assign --user <email> --role <role>
│   └── list
├── sync
├── web
├── tree - this command
└── version`

// treeCmd represents the tree command
var treeCmd = &cobra.Command{
	Use:   "tree",
	Short: "Shows the full command tree structure.",
	Long:  `Displays a tree structure of all available commands, subcommands, and their arguments.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(commandTree)
	},
}

func init() {
	rootCmd.AddCommand(treeCmd)
}
