package services

import (
	"encoding/json"
	"fmt"
	"github.com/Torakushi/example-ws-events/internal/websocket"

	"github.com/Torakushi/example-ws-events/internal/models"
	"github.com/Torakushi/example-ws-events/internal/persistence"
)

type EventService struct {
	eventRepository *persistence.EventRepository
	wsHub           *websocket.Hub
}

func NewEventService(eventRepository *persistence.EventRepository, wsHub *websocket.Hub) *EventService {
	return &EventService{
		eventRepository: eventRepository,
		wsHub:           wsHub,
	}
}

func (s *EventService) ProcessEvent(body []byte) {
	event := &models.Event{}
	err := json.Unmarshal(body, event)
	if err != nil {
		fmt.Println(1)
		fmt.Println(err)
		return
	}

	// Save the event to the database
	err = s.eventRepository.SaveEvent(event)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Broadcast the event to all clients
	eventData, _ := json.Marshal(event)
	s.wsHub.Broadcast(eventData)
}
