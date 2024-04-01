package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

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

		handleWebSocket(ws)
	})

	router.Run(":8080")
}

func handleWebSocket(ws *websocket.Conn) {
	defer ws.Close()

	// Read messages from the WebSocket connection
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			break
		}

		// Print received message to the console
		println("Received message:", string(msg))

		// Echo the received message back to the client
		err = ws.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}
