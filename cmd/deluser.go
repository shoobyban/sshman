package cmd

import (
	"strings"

	"github.com/apioapp/slog"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

// delCmd represents the del command
var delCmd = &cobra.Command{
	Use:   "del",
	Short: "Delete user by email",
	Long:  `Check all servers and delete user with given email`,
	Run: func(_ *cobra.Command, args []string) {
		conf := readConfig()
		for _, email := range args {
			u := findByEmail(conf, email)
			if u != nil {
				//slog.Infof("User %#v", u)
				delUser(conf, u)
			} else {
				slog.Infof("No such user")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(delCmd)
}

func delUser(C *config, deluser *user) {
	bar := progressbar.Default(int64(len(C.Hosts)))
	for alias, host := range C.Hosts {
		bar.Add(1)
		key := host.Key
		if key == "" {
			key = C.Key
		}
		client, err := connect(key, host.Host, host.User)
		if err != nil {
			slog.Errorf("error connecting %s: %v", alias, err)
			continue
		}
		b, err := client.Read()
		if err != nil {
			slog.Errorf("error reading authorized keys on %s: %v", alias, err)
			continue
		}
		userlist := []string{}
		have := false
		sum := checksum(string(b))
		lines := strings.Split(string(b), "\n")
		newlines := []string{}
		for _, line := range lines {
			if len(line) == 0 {
				continue
			}
			parts := strings.Split(line, " ")
			if len(parts) != 3 {
				slog.Errorf("Not good line: '%s'", line)
			}
			lsum := checksum(parts[1])
			if _, ok := C.Users[lsum]; !ok {
				delete(C.Users, lsum)
			}
			if parts[1] == deluser.Key {
				have = true
				continue
			}
			newlines = append(newlines, line)
			userlist = append(userlist, C.Users[lsum].Email)
		}

		if have {
			newlines = deleteEmpty(newlines)
			err = client.Write(strings.Join(newlines, "\n") + "\n")
			if err != nil {
				slog.Errorf("error writing %s: %v", alias, err)
			}
		}
		host.Checksum = sum
		host.Users = userlist
		C.Hosts[alias] = host
		client.Close()
	}
	writeConfig(C)
}
