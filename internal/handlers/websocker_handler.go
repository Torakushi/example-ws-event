package handlers

import (
	websocket2 "github.com/Torakushi/example-ws-events/internal/websocket"
	"log"
	"net/http"

	gorilla "github.com/gorilla/websocket"
)

var upgrader = gorilla.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for simplicity. You can customize this based on your needs.
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebsocketHandler(hub *websocket2.Hub, w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade WebSocket connection: ", err)
		return
	}

	client := websocket2.NewClient(conn, hub)

	// Register the client with the hub
	hub.RegisterClient(client)

	go client.Read()
	go client.Write()
}
