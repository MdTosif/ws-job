package main

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mdtosif/ws-job/pkg/ws/handler"
)

var upgrader = websocket.Upgrader{

}

func main() {
	

	// http server running on port 4444 and handling websocket connections
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// websocket handler
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			println(err.Error())
		}
		
		// handle the of all websocket
		handler.HandleWsConn(conn)
		
	})
	http.ListenAndServe(":4444", nil)

}