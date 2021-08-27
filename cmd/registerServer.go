package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/apioapp/slog"
	"github.com/spf13/cobra"
)

// registerServerCmd represents the register server command
var registerServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Register a server",
	Long: `To register a user:
	sshman register server {alias} {server_address:port} {user} {~/.ssh/working_keyfile} [group1 group2 ...]`,
	Run: func(_ *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println(`To register a user:
			sshman register user email sshkey.pub {group1 group2 ...}`)
			os.Exit(0)
		}
		conf := readConfig()
		u := findByEmail(conf, args[0])
		if u != nil {
			fmt.Printf("User already exists with this email, overwrite [y/n]: ")
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				slog.Fatalf("error opening stdout %v", err)
			}
			response = strings.ToLower(strings.TrimSpace(response))
			if response != "y" && response != "yes" {
				os.Exit(0)
			}
		}
		registerUser(conf, args...)
	},
}

func init() {
	registerCmd.AddCommand(registerUserCmd)
}

func registerUser(C *config, args ...string) error {
	b, err := ioutil.ReadFile(args[1])
	if err != nil {
		slog.Errorf("error reading public key file: '%s' %v", args[1], err)
		return err
	}
	parts := strings.Split(strings.TrimSuffix(string(b), "\n"), " ")
	if len(parts) != 3 {
		slog.Errorf("not a proper public key file")
	}
	lsum := checksum(parts[1])
	newuser := user{
		KeyType: parts[0],
		Key:     parts[1],
		Name:    parts[2],
		Email:   args[0],
	}
	if len(args) > 2 {
		newuser.Groups = args[2:]
	}
	C.Users[lsum] = newuser
	slog.Infof("Registering %s %s %s %s", parts[0], parts[2], args[0], lsum)
	writeConfig(C)
	return nil
}
