package main

import (
	"fmt"
	"github.com/andrewbackes/webshell/pkg/auth/github"
	"github.com/andrewbackes/webshell/pkg/session"
	"github.com/andrewbackes/webshell/pkg/shell"
	"github.com/andrewbackes/webshell/pkg/websocket"
	"net/http"
	"os"
)

const (
	authConfigFile    = "auth_config.json"
	authEnabledEnvVar = "AUTH_ENABLED"
)

func main() {
	fmt.Println("WebShell started.")

	var ws *websocket.Server
	ws = websocket.NewServer(websocket.Handler(func(m websocket.Message) {
		fmt.Println(m)
		ws.Write([]byte("$ " + string(m.Value)))
		go shell.Run(m.Value, ws)
	}))

	if os.Getenv(authEnabledEnvVar) == "true" {
		fmt.Println("Auth enabled.")
		sessions := session.New()
		gh := github.New(authConfigFile, sessions)
		http.HandleFunc("/auth/login", gh.Login)
		http.HandleFunc("/auth/callback", gh.Callback)
		http.HandleFunc("/", sessions.Middleware(staticHandler))
		http.HandleFunc("/websocket", sessions.Middleware(ws.UpgradeHandler))
	} else {
		http.HandleFunc("/", staticHandler)
		http.HandleFunc("/websocket", ws.UpgradeHandler)
	}

	http.ListenAndServe(":8080", nil)
}

// Handler for the react app. It must be built prior to running.
func staticHandler(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("webapp/build"))
	fs.ServeHTTP(w, r)
}
