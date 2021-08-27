package cmd

import (
	"strings"

	"github.com/apioapp/slog"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "update",
	Short: "Read users into configuration",
	Long:  `Loop through all servers, download all users from autorized_keys into configuration`,
	Run: func(cmd *cobra.Command, _ []string) {
		readUsers(cmd, readConfig())
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
}

func readUsers(_ *cobra.Command, C *config) {
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
		client.Close()
		sum := checksum(string(b))
		userlist := []string{}
		lines := strings.Split(string(b), "\n")
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
				C.Users[lsum] = user{
					KeyType: parts[0],
					Key:     parts[1],
					Name:    parts[2],
					Email:   parts[2] + "@" + alias,
				}
			}
			userlist = append(userlist, C.Users[lsum].Email)
		}
		host.Checksum = sum
		host.Users = userlist
		C.Hosts[alias] = host
	}
	writeConfig(C)
}
