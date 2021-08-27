package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// initConfig reads in config file and ENV variables if set.
func readConfig() *config {
	var C config
	home, _ := os.UserHomeDir()
	b, err := ioutil.ReadFile(home + "/.ssh/.sshman")
	if err != nil {
		log.Printf("Error: unable to read .sshman, %v\n", err)
	}
	err = json.Unmarshal(b, &C)
	if err != nil {
		log.Printf("Error: unable to decode into struct, %v\n", err)
	}
	return &C
}

func writeConfig(c *config) {
	b, _ := json.MarshalIndent(c, "", "  ")
	home, _ := os.UserHomeDir()
	ioutil.WriteFile(home+"/.ssh/.sshman", b, 0644)
}
