package main

import (
	"encoding/json"
	"log"
	"os"

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
	// queue, err := amqpChannel.QueueDeclare("test", true, false, false, false, nil)
	// handleError(err, "Could not declare `add` queue")

	// err = amqpChannel.Qos(10, 0, false)
	// handleError(err, "Could not configure Qos")

	err = amqpChannel.ExchangeDeclare(
		"logs",              // name
		amqp.ExchangeFanout, // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	handleError(err, "Failed to declare an exchange")

	queue, err := amqpChannel.QueueDeclare(
		"test_lai", // name
		false,      // durable
		false,      // delete when unused
		true,       // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	handleError(err, "Failed to declare a queue")

	err = amqpChannel.QueueBind(
		queue.Name, // queue name
		"",         // routing key
		"logs",     // exchange
		false,
		nil,
	)

	// consume queue
	messageQueue, err := amqpChannel.Consume(queue.Name, "", false, false, false, false, nil)
	handleError(err, "Counld not register consumer")

	stopChan := make(chan bool)

	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())

		for d := range messageQueue {
			log.Printf("Received a message: %s", d.Body)
			task := &config.Task{}
			err := json.Unmarshal(d.Body, task)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}

			log.Printf("Result of %d + %d is : %d", task.Number1, task.Number2, task.Number1+task.Number2)
			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("Acknowledged message")
			}
		}
	}()

	<-stopChan
}
