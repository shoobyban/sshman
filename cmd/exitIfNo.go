package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func exitIfNo() {
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error: error opening stdout %v\n", err)
	}
	response = strings.ToLower(strings.TrimSpace(response))
	if response != "y" && response != "yes" {
		os.Exit(0)
	}
}
