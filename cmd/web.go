package cmd

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/shoobyban/sshman/api"
	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

//go:embed "dist"
var dist embed.FS

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
		users := api.Users{
			Prefix: "/api/users",
			Config: backend.ReadConfig(),
		}
		hosts := api.Hosts{
			Prefix: "/api/hosts",
			Config: backend.ReadConfig(),
		}
		r := chi.NewMux()
		r.Use(middleware.Logger)
		r = users.Routers("/api/users", r)
		r = hosts.Routers("/api/hosts", r)
		log.Printf("Listening on http://localhost:%v", port)
		rs, err := dist.ReadDir(".")
		if err != nil {
			log.Fatal(err)
		}
		for _, r := range rs {
			log.Printf("File %v", r.Name())
		}
		http.Handle("/*", http.FileServer(http.FS(dist)))
		fmt.Println(http.ListenAndServe(":"+strconv.Itoa(port), r))
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.PersistentFlags().IntP("port", "p", 80, "Port for web UI. Defaults to 80")
}
