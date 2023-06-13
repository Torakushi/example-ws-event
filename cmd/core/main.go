package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Torakushi/example-ws-events/internal/amqp"
	"github.com/Torakushi/example-ws-events/internal/handlers"
	"github.com/Torakushi/example-ws-events/internal/persistence"
	"github.com/Torakushi/example-ws-events/internal/services"
	"github.com/Torakushi/example-ws-events/internal/websocket"
)

func main() {
	amqpConn, err := amqp.NewConnection("amqp://rmq:rmq@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer amqpConn.Close()

	pgConn, err := persistence.NewConnection("postgres://postgres:postgres@localhost:5432/nested?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer pgConn.Close()

	eventRepository := persistence.NewEventRepository(pgConn)

	lastEvents, err := eventRepository.GetLastEvents(websocket.RetryPolicy)
	if err != nil {
		log.Fatalf("Can t retrieve last events in DB: %v", err)
	}

	// Set up WebSocket hub
	wsHub, err := websocket.NewHub(lastEvents)
	if err != nil {
		log.Fatalf("Can t instantiate hub: %v", err)
	}
	go wsHub.Run()

	// Set up HTTP server
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handlers.WebsocketHandler(wsHub, w, r)
	})

	eventService := services.NewEventService(eventRepository, wsHub)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Closing connections...")
		pgConn.Close()
		os.Exit(0)
	}()

	// Start consuming events from RabbitMQ
	err = amqpConn.ConsumeEvents(eventService.ProcessEvent)
	if err != nil {
		log.Fatalf("Failed to start consuming events: %v", err)
	}

	log.Println("Waiting for events...")
	http.ListenAndServe(":8080", nil)
}
