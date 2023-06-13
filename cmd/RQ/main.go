package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// Establish connection to RabbitMQ server
	conn, err := amqp.Dial("amqp://rmq:rmq@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare a queue
	q, err := ch.QueueDeclare(
		"events",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	// Ask for user input and send messages to the queue
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter messages (press Ctrl+C to exit):")

	for scanner.Scan() {
		message := scanner.Text()

		// Publish the message to the queue
		err = ch.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(message),
			})
		failOnError(err, "Failed to publish a message")

		fmt.Println("Message sent to the queue:", message)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
