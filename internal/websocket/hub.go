package websocket

import (
	"encoding/json"
	"github.com/Torakushi/example-ws-events/internal/models"
	"log"
)

const RetryPolicy = 10

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client

	lastEvents [][]byte
}

func NewHub(lastEvents []*models.Event) (*Hub, error) {
	// TODO: Better initialisation?
	events := make([][]byte, len(lastEvents))
	for i := 0; i < len(events); i++ {
		j, err := json.Marshal(lastEvents[i])
		if err != nil {
			return nil, err
		}

		events[i] = j
	}

	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		lastEvents: events,
	}, nil
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		case message := <-h.broadcast:
			h.updatePastEvent(message)

			for client := range h.clients {
				client.send <- message
			}
		}
	}
}

func (h *Hub) updatePastEvent(b []byte) {
	h.lastEvents = append(h.lastEvents, b)

	if len(h.lastEvents) == RetryPolicy+1 {
		h.lastEvents = h.lastEvents[1:]
	}
}

func (h *Hub) RegisterClient(client *Client) {
	log.Println("New Client")
	h.register <- client
}

func (h *Hub) UnregisterClient(client *Client) {
	log.Println("Bye Client")

	h.unregister <- client
}

func (h *Hub) Broadcast(message []byte) {
	log.Println("Broadcast")

	h.broadcast <- message
}

func (h *Hub) GetPastEvents() [][]byte {
	return h.lastEvents
}
