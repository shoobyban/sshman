package cmd

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/shoobyban/sshman/api"
	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

//go:embed dist/*
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
		groups := api.Groups{
			Prefix: "/api/groups",
			Config: backend.ReadConfig(),
		}
		r := chi.NewMux()
		r.Use(middleware.Logger)
		r = users.Routers("/api/users", r)
		r = hosts.Routers("/api/hosts", r)
		r = groups.Routers("/api/groups", r)
		log.Printf("Listening on http://localhost:%v", port)
		distfs, err := fs.Sub(dist, "dist")
		if err != nil {
			log.Fatal(err)
		}
		r.Handle("/*", http.FileServer(http.FS(distfs)))
		fmt.Println(http.ListenAndServe(":"+strconv.Itoa(port), r))
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.PersistentFlags().IntP("port", "p", 80, "Port for web UI. Defaults to 80")
}
