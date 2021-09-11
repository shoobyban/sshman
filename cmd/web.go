package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// go:embed ../frontend/dist

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Web UI",
	Long:  `Stays running and created a web UI.`,
	Run: func(cmd *cobra.Command, _ []string) {
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			port = 80
		}
		http.Handle("/", http.FileServer(http.FS(os.DirFS("../frontend/dist"))))
		fmt.Println(http.ListenAndServe(":"+strconv.Itoa(port), nil))
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.PersistentFlags().IntP("port", "p", 80, "Port for web UI. Defaults to 80")
}
