package cmd

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/shoobyban/sshman/api"
	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

//go:embed dist/*
var dist embed.FS

// ReadConfig is a go-chi middleware that reads a fresh config and adds it to the request context
func ReadConfig(ilog *backend.ILog) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log.Printf("[DEBUG] Creating config")
		cfg := backend.DefaultConfig()
		cfg.SetLog(ilog)
		cfg.Storage.Watch(func() {
			//			cfg.Log().Infof("storage changed, reloading")
			cfg = backend.DefaultConfig()
		})
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), api.ConfigKey, cfg)
			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}

func isLoopbackBind(bind string) bool {
	switch bind {
	case "127.0.0.1", "localhost", "::1":
		return true
	default:
		return false
	}
}

func spaFileHandler(distfs fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(distfs))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		requestedPath := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if requestedPath == "." || requestedPath == "" {
			http.ServeFileFS(w, r, distfs, "index.html")
			return
		}

		if info, err := fs.Stat(distfs, requestedPath); err == nil && !info.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}

		http.ServeFileFS(w, r, distfs, "index.html")
	})
}

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Web UI",
	Long:  `Stays running and created a web UI.`,
	Run: func(cmd *cobra.Command, _ []string) {
		port, err := cmd.Flags().GetString("port")
		if err != nil || port == "dynamic" {
			port = ":0"
		}
		_, err = strconv.Atoi(port)
		if err == nil {
			port = ":" + port
		}
		r := chi.NewMux()
		r.Use(middleware.Logger)
		weblog := backend.NewLog(true)
		r.Use(ReadConfig(weblog))
		api.GroupsHandler{Prefix: "/api/groups"}.AddRoutes(r)
		api.HostsHandler{Prefix: "/api/hosts"}.AddRoutes(r)
		api.UsersHandler{Prefix: "/api/users"}.AddRoutes(r)
		api.LogsHandler{Prefix: "/api/logs"}.AddRoutes(r)
		enableKeysAPI, err := cmd.Flags().GetBool("enable-keys-api")
		if err != nil {
			log.Fatal(err)
		}
		if enableKeysAPI {
			api.KeysHandler{Prefix: "/api/keys"}.AddRoutes(r)
		}

		distfs, err := fs.Sub(dist, "dist")
		if err != nil {
			log.Fatal(err)
		}
		r.Handle("/*", spaFileHandler(distfs))
		listener, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatal(err)
		}
		if port == ":0" {
			fmt.Printf("Listening on http://localhost:%d\n", listener.Addr().(*net.TCPAddr).Port)
			portfile, _ := cmd.Flags().GetString("portfile")
			os.WriteFile(portfile, []byte(fmt.Sprint(listener.Addr().(*net.TCPAddr).Port)), 0644)
		}
		bind, err := cmd.Flags().GetString("bind")
		if err != nil {
			log.Fatal(err)
		}
		allowRemote, err := cmd.Flags().GetBool("allow-remote")
		if err != nil {
			log.Fatal(err)
		}
		if !allowRemote && !isLoopbackBind(bind) {
			log.Fatal("refusing to bind the admin web UI to a non-loopback address without --allow-remote")
		}
		server := &http.Server{Addr: bind + port, Handler: r}
		fmt.Println(server.Serve(listener))
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.PersistentFlags().StringP("bind", "b", "127.0.0.1", "Bind to IP address for web UI. Defaults to loopback for local-only access.")
	webCmd.PersistentFlags().StringP("port", "p", "dynamic", "Port for web UI. Can be a port number or 'dynamic' (without quotes). Defaults to dynamic address. Dynamic address will create a sshman.port file.")
	webCmd.PersistentFlags().StringP("portfile", "f", "sshman.port", "Port filename for dynamic address to check, can be relative or full path, don't use ~ or $HOME")
	webCmd.PersistentFlags().Bool("allow-remote", false, "Allow binding the admin web UI to a non-loopback address.")
	webCmd.PersistentFlags().Bool("enable-keys-api", false, "Expose the keys API over the web UI. Disabled by default.")
}
