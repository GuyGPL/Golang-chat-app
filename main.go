package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Message struct {
	Username string
	Message  string
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}, ReadBufferSize: 1024, WriteBufferSize: 1024}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

func main() {
	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Welcome to chat application"})
	})

	router.GET("/ws", func(ctx *gin.Context) {
		ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
			return
		}

		handleConnections(ws)
	})

	go handleMessages()

	router.Run(":8080")
}

func handleConnections(ws *websocket.Conn) {
	defer ws.Close()

	clients[ws] = true

	// Read messages from the WebSocket connection
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			delete(clients, ws)
			break
		}

		message := Message{Username: ws.RemoteAddr().String(), Message: string(msg)}
		// Print received message to the console
		println("Received message:", message.Message)

		broadcast <- message
	}
}

func handleMessages() {
	for {
		msg := <-broadcast

		// format message
		messageToSend := fmt.Sprintf("[%s]: %s", msg.Username, msg.Message)

		for client := range clients {
			// Echo the received message back to all clients
			err := client.WriteMessage(websocket.TextMessage, []byte(messageToSend))

			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}
