package main

import (
	"encoding/json"
	"log"
	"math/rand"

	"github.com/streadway/amqp"
	"hoangduc02.01.1998/rabbitmq-go/config"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// connect to rmq
	conn, err := amqp.Dial(config.Config.AMQPConnectionURL)
	handleError(err, "Can't connect to AMQ")
	defer conn.Close()

	// craete channel
	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create channel")
	defer amqpChannel.Close()

	// config queue channel
	// queue, err := amqpChannel.QueueDeclare("213", true, false, false, false, nil)

	err = amqpChannel.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	handleError(err, "Could not declare `add` queue")

	// create message
	task := config.Task{Number1: rand.Intn(100), Number2: rand.Intn(999)}

	body, err := json.Marshal(task)
	if err != nil {
		handleError(err, "Error encoding JSON")
	}

	// publish
	err = amqpChannel.Publish("logs", "", false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})

	if err != nil {
		log.Fatalf("Error publishing message: %s", err)
	}

	log.Printf("AddTask: %d+%d", task.Number1, task.Number2)
}
