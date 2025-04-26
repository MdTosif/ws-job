package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/gorilla/websocket"
	"github.com/mdtosif/ws-job/pkg/ws/handler"
)

var upgrader = websocket.Upgrader{

}

func main() {
	
	// run this command only if development mode 
	/* go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil)) // pprof runs here
	}()

	go client.RunParallel() */



	// http server running on port 4444 and handling websocket connections
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pkg/html/index.html")
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
	
	log.Println("Server started on port 4444")

	err := http.ListenAndServe(":4444", nil)

	if err != nil {
		panic(err)
	}	


}