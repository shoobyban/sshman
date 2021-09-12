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

type httpError struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type httpResponse struct {
	Message string `json:"message,omitempty"`
	Data    string `json:"data,omitempty"`
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
			fmt.Printf("Frontend dev server proxy URL: %v\n", purl)
			if purl == "" {
				return
			}
		}

		http.Handle("/api/config", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			cfg := backend.ReadConfig()
			wResponse(w, cfg)
		}))
		http.Handle("/api/groups", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			cfg := backend.ReadConfig()
			wResponse(w, cfg.GetGroups())
		}))
		http.Handle("/api/user", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if postCheck(w, r) {
				return
			}
			conf := backend.ReadConfig()
			var req struct {
				Email  string   `json:"email"`
				Key    string   `json:"key"`
				Groups []string `json:"groups"`
			}
			err = decodeBody(w, r, &req)
			if err != nil {
				return
			}
			_, existing := conf.GetUserByEmail(req.Email)
			if existing != nil {
				wError(w, httpError{Error: "Email already exists"})
				return
			}
			u := []string{req.Email, req.Key}
			err = conf.RegisterUser([]string{}, append(u, req.Groups...)...)
			if err != nil {
				wError(w, httpError{Error: err.Error()})
			}
		}))
		http.Handle("/api/server", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if postCheck(w, r) {
				return
			}
			conf := backend.ReadConfig()
			var host backend.Hostentry
			err = decodeBody(w, r, &host)
			if err != nil {
				return
			}
			if _, ok := conf.Hosts[host.Alias]; ok {
				wError(w, httpError{Error: "Email already exists"})
				return
			}
			host.Config = conf
			err = conf.RegisterServer([]string{}, append([]string{host.Alias, host.Host, host.User, host.Key}, host.Groups...)...)
			if err != nil {
				wError(w, httpError{Error: err.Error()})
			}
		}))
		http.Handle("/api/groups/user", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if postCheck(w, r) {
				return
			}
			cfg := backend.ReadConfig()
			var req struct {
				Groups []string `json:"groups"`
				Emails []string `json:"emails"`
			}
			err = decodeBody(w, r, &req)
			if err != nil {
				return
			}
			if len(req.Emails) != 2 {
				wError(w, httpError{Error: "Invalid arguments"})
			}
			key, user := cfg.GetUserByEmail(req.Emails[0])
			if user == nil {
				key, user = cfg.GetUserByEmail(req.Emails[1])
			}
			if user == nil {
				wError(w, httpError{Error: "No such user"})
				return
			}
			old := user.GetGroups()
			user.SetGroups(req.Groups)
			user.UpdateGroups(cfg, old)
			cfg.Users[key] = *user
			cfg.Write()
			wResponse(w, httpResponse{Message: fmt.Sprintf("Changed groups for %s to %v\"}", user.Email, req.Groups)})
		}))
		http.Handle("/api/rename/user", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if postCheck(w, r) {
				return
			}
			cfg := backend.ReadConfig()
			var req struct {
				New string `json:"new"`
				Old string `json:"old"`
			}
			err = decodeBody(w, r, &req)
			if err != nil {
				return
			}
			key, user := cfg.GetUserByEmail(req.Old)
			if user != nil {
				user.Email = req.New
				cfg.Users[key] = *user
				cfg.Write()
				wResponse(w, httpResponse{Message: fmt.Sprintf("Renamed %s to %s\"}", req.Old, req.New)})
			}
		}))
		http.Handle("/api/groups/server", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if postCheck(w, r) {
				return
			}
			cfg := backend.ReadConfig()
			var req struct {
				Groups  []string `json:"groups"`
				Aliases []string `json:"aliases"`
			}
			err = decodeBody(w, r, &req)
			if err != nil {
				return
			}
			if len(req.Aliases) != 2 {
				wError(w, httpError{Error: "Invalid arguments"})
			}
			var host backend.Hostentry
			var ok bool
			var alias string
			if host, ok = cfg.Hosts[req.Aliases[0]]; !ok {
				if host, ok = cfg.Hosts[req.Aliases[1]]; !ok {
					wError(w, httpError{Error: "No such host"})
					return
				}
				alias = req.Aliases[1]
			} else {
				alias = req.Aliases[0]
			}
			old := host.GetGroups()
			host.SetGroups(req.Groups)
			host.UpdateGroups(cfg, old)
			cfg.Hosts[alias] = host
			cfg.Write()
			wResponse(w, httpResponse{Message: fmt.Sprintf("Changed groups for %s to %v\"}", alias, req.Groups)})
		}))
		http.Handle("/api/rename/server", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if postCheck(w, r) {
				return
			}
			cfg := backend.ReadConfig()
			var req struct {
				New string `json:"new"`
				Old string `json:"old"`
			}
			err = decodeBody(w, r, &req)
			if err != nil {
				return
			}
			if host, ok := cfg.Hosts[req.Old]; ok {
				cfg.Hosts[req.New] = host
				delete(cfg.Hosts, req.Old)
				cfg.Write()
				wResponse(w, httpResponse{Message: fmt.Sprintf("Renamed %s to %s\"}", req.Old, req.New)})
			}
		}))
		http.Handle("/", getHandler(cmd))

		// start server
		if port == 80 {
			fmt.Println("Listening on http://localhost/")
		} else {
			fmt.Printf("Listening on http://localhost:%d/\n", port)
		}
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

func postCheck(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != "POST" || r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return true
	}
	return false
}

func decodeBody(w http.ResponseWriter, r *http.Request, req interface{}) error {
	r.Body = http.MaxBytesReader(w, r.Body, 655360) // 640K ought to be enough for anybody
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	return err
}

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

func wError(w http.ResponseWriter, errorDetails interface{}) {
	w.WriteHeader(http.StatusBadRequest)
	b, _ := json.Marshal(errorDetails)
	w.Write(b)
}

func wResponse(w http.ResponseWriter, response interface{}) {
	b, _ := json.Marshal(response)
	w.Write(b)
}
