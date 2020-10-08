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
	ch, err := conn.Channel()
	handleError(err, "Can't create channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	handleError(err, "Failed to declare exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)

	err = ch.QueueBind(
		q.Name,       // queue name
		"#.black",    // routing key
		"logs_topic", // exchange
		false,
		nil)
	handleError(err, "Failed to bind a queue")

	messageQueue, err := ch.Consume(q.Name, "", false, false, false, false, nil)
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
