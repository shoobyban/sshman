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
		cfg := backend.ReadConfig(true)
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			port = 80
		}

		r := chi.NewMux()
		r.Use(middleware.Logger)
		api.Groups{Prefix: "/api/groups", Config: cfg}.AddRoutes(r)
		api.Hosts{Prefix: "/api/hosts", Config: cfg}.AddRoutes(r)
		api.Users{Prefix: "/api/users", Config: cfg}.AddRoutes(r)
		api.Logs{Prefix: "/api/logs", Config: cfg}.AddRoutes(r)
		api.Keys{Prefix: "/api/keys", Config: cfg}.AddRoutes(r)

		cfg.Log.Infof("Listening on http://localhost:%v", port)
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
