package main

import (
	"fmt"
	"github.com/andrewbackes/webshell/pkg/websocket"
	"net/http"
)

func main() {
	fmt.Println("WebShell started.")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fs := http.FileServer(http.Dir("webapp/build"))
		fs.ServeHTTP(w, r)
	})
	var ws *websocket.Server
	ws = websocket.NewServer(websocket.MessageHandler(func(m websocket.Message) {
		fmt.Println(m)
		ws.Write(m.Value)
	}))
	http.HandleFunc("/websocket", ws.UpgradeHandler)
	http.ListenAndServe(":8080", nil)
}
