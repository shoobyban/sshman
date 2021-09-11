package cmd

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/shoobyban/sshman/backend"
	"github.com/spf13/cobra"
)

//go:embed dist
var embededFiles embed.FS
var devmode bool

func getHandler(cmd *cobra.Command) http.Handler {
	if devmode {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			purl, _ := cmd.Flags().GetString("devserver")
			fmt.Printf("Proxying %s%v\n", purl, r.URL.Path)
			devurl, _ := url.Parse(purl)
			proxy := httputil.NewSingleHostReverseProxy(devurl)
			r.URL.Host = devurl.Host
			r.URL.Scheme = devurl.Scheme
			r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
			r.Host = devurl.Host
			proxy.ServeHTTP(w, r)
		})
	}

	fsys, err := fs.Sub(embededFiles, "dist")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(fsys))
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
		if devmode {
			purl, _ := cmd.Flags().GetString("devserver")
			fmt.Printf("Proxy URL: %v\n", purl)
			if purl == "" {
				return
			}
		}
		http.Handle("/api/", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			cfg := backend.ReadConfig()
			b, _ := json.Marshal(cfg)
			w.Write(b)
		}))
		http.Handle("/", getHandler(cmd))
		fmt.Println(http.ListenAndServe(":"+strconv.Itoa(port), nil))
	},
}

func init() {
	devmode = strings.Contains(os.Args[0], "/exe/")
	rootCmd.AddCommand(webCmd)
	webCmd.PersistentFlags().IntP("port", "p", 80, "Port for web UI. Defaults to 80")
	if devmode {
		webCmd.PersistentFlags().StringP("devserver", "d", "http://localhost:3000/", "UI dev server URL to proxy static calls to")
	}
}
