package cmd

import (
	"log"
	"strings"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add user",
	Long:  `Add already existing user by email`,
	Run: func(_ *cobra.Command, args []string) {
		conf := readConfig()
		for _, email := range args {
			u := findByEmail(conf, email)
			if u != nil {
				//log.Printf("User %#v\n", u)
				addUser(conf, u)
			} else {
				log.Printf("No such user\n")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func addUser(C *config, newuser *user) {
	bar := progressbar.Default(int64(len(C.Hosts)))
	for alias, host := range C.Hosts {
		bar.Add(1)
		if len(newuser.Groups) != 0 && !match(newuser.Groups, host.Groups) {
			continue
		}
		key := host.Key
		if key == "" {
			key = C.Key
		}
		client, err := connect(key, host.Host, host.User)
		if err != nil {
			log.Printf("Error: error connecting %s: %v\n", alias, err)
			continue
		}
		b, err := client.Read()
		if err != nil {
			log.Printf("Error: error reading authorized keys on %s: %v\n", alias, err)
			continue
		}
		userlist := []string{}
		have := false
		sum := checksum(string(b))
		lines := strings.Split(string(b), "\n")
		for _, line := range lines {
			if len(line) == 0 {
				continue
			}
			parts := strings.Split(line, " ")
			if len(parts) != 3 {
				log.Printf("Error: Not good line: '%s'\n", line)
			}
			if parts[0] == newuser.KeyType &&
				parts[1] == newuser.Key &&
				parts[2] == newuser.Name {
				have = true
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

		if have {
			host.Checksum = sum
			host.Users = userlist
			C.Hosts[alias] = host
			continue
		} else {
			lines = deleteEmpty(append(lines, newuser.KeyType+" "+newuser.Key+" "+newuser.Name))

			err = client.Write(strings.Join(lines, "\n") + "\n")
			if err != nil {
				log.Printf("Error: error writing %s: %v\n", alias, err)
			}
			userlist = append(userlist, newuser.Email)
		}
		host.Checksum = sum
		host.Users = userlist
		C.Hosts[alias] = host
		client.Close()
	}
	writeConfig(C)
}
