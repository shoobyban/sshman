package cmd

import (
	"context"
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

func ReadConfig(log *backend.ILog) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			cfg := backend.WebReadConfig(log)
			ctx := context.WithValue(r.Context(), "config", cfg)
			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}

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

		r := chi.NewMux()
		r.Use(middleware.Logger)
		weblog := backend.NewLog(true)
		r.Use(ReadConfig(weblog))
		api.Groups{Prefix: "/api/groups"}.AddRoutes(r)
		api.Hosts{Prefix: "/api/hosts"}.AddRoutes(r)
		api.Users{Prefix: "/api/users"}.AddRoutes(r)
		api.Logs{Prefix: "/api/logs"}.AddRoutes(r)
		api.Keys{Prefix: "/api/keys"}.AddRoutes(r)

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
