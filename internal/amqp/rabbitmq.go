package amqp

import (
	"fmt"

	"github.com/streadway/amqp"
)

type Connection struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   *amqp.Queue
}

func NewConnection(amqpURI string) (*Connection, error) {
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	queue, err := channel.QueueDeclare(
		"events",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to declare a queue: %v", err)
	}

	return &Connection{
		conn:    conn,
		channel: channel,
		queue:   &queue,
	}, nil
}

func (c *Connection) ConsumeEvents(handlerFunc func([]byte)) error {
	msgs, err := c.channel.Consume(
		c.queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to consume events from queue: %v", err)
	}

	// Handle received events
	go func() {
		for event := range msgs {
			// Process the received event (e.g., save to the database, broadcast to WebSocket clients, etc.)
			handlerFunc(event.Body)
		}
	}()

	return nil
}

func (c *Connection) Close() error {
	err := c.channel.Close()
	if err != nil {
		return fmt.Errorf("failed to close channel: %v", err)
	}

	err = c.conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close connection: %v", err)
	}

	return nil
}
