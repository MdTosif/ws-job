package handler

import (
	"bufio"
	"io"
	"log"

	"slices"

	"github.com/gorilla/websocket"
	"github.com/mdtosif/ws-job/cmd/runner"
)

var openConnections = []*websocket.Conn{}

func HandleWsConn(conn *websocket.Conn) {
	defer conn.Close() // Ensure the connection is closed when the handler exits
	openConnections = append(openConnections, conn)

	var exec = runner.New()

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// Handle different message types
		switch msgType {
		case websocket.TextMessage:
			// Handle text message (UTF-8 encoded string)
			log.Printf("Received text message: %s", msg)
			stdoutPipe, stderrPipe, err := exec.Run(string(msg))

			if err != nil {
				log.Println("Error executing command:", err)
				// Send error message to client
				if err := conn.WriteMessage(websocket.TextMessage, []byte(err.Error())); err != nil {
					log.Println("Error sending error message:", err)
				}
				break
			}

			stream := func(conn *websocket.Conn, pipe io.ReadCloser) {

				scanner := bufio.NewScanner(pipe)
				for scanner.Scan() {
					if err := conn.WriteMessage(websocket.TextMessage, (scanner.Bytes())); err != nil {
						log.Println("Error sending message:", err)
						break
					}
				}

			}

			go stream(conn, stdoutPipe)
			go stream(conn, stderrPipe)

		case websocket.BinaryMessage:
			// Handle binary message (e.g., file upload, image)
			log.Printf("Received binary message: %v", msg)
			// Echo back the binary message
			// if err := conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
			// 	log.Println("Error sending binary message:", err)
			// 	break
			// }

		case websocket.CloseMessage:
			// Handle close message (client wants to close the connection)
			log.Println("Received close message")
			// Send close message to client (gracefully close the connection)
			if err := conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
				log.Println("Error sending close message:", err)
			}
			return

		case websocket.PingMessage:
			// Handle ping message (client is checking if the server is still alive)
			log.Println("Received ping message")
			// Send pong message in response
			if err := conn.WriteMessage(websocket.PongMessage, nil); err != nil {
				log.Println("Error sending pong message:", err)

			}

		case websocket.PongMessage:
			// Handle pong message (responding to a ping or just acknowledging the server's response)
			log.Println("Received pong message")
			// Typically, you don't need to send a pong message in response to an incoming pong message

		default:
			log.Printf("Received unknown message type: %d", msgType)
		}
	}

	conn.SetCloseHandler(func(code int, text string) error {
		log.Printf("Connection closed with code %d: %s", code, text)
		exec.Stop()
		// remove from openConnections
		for i, c := range openConnections {
			if c == conn {
				openConnections = slices.Delete(openConnections, i, i+1)
				break
			}
		}
		return nil
	})

}
