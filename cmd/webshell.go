package main

import (
	"fmt"
	"github.com/andrewbackes/webshell/pkg/shell"
	"github.com/andrewbackes/webshell/pkg/websocket"
	"net/http"
)

func main() {
	fmt.Println("WebShell started.")
	// Handler for the react app. It must be built prior to running.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fs := http.FileServer(http.Dir("webapp/build"))
		fs.ServeHTTP(w, r)
	})
	var ws *websocket.Server
	ws = websocket.NewServer(websocket.Handler(func(m websocket.Message) {
		fmt.Println(m)
		ws.Write([]byte("$ " + string(m.Value)))
		go shell.Run(m.Value, ws)
	}))
	http.HandleFunc("/websocket", ws.UpgradeHandler)
	http.ListenAndServe(":8080", nil)
}
