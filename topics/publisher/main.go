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

	ch, err := conn.Channel()
	handleError(err, "Can't connect to Rabbitmq")
	defer conn.Close()

	err = ch.ExchangeDeclare(
		"logs_topic",
		amqp.ExchangeTopic,
		true,
		false,
		false,
		false,
		nil,
	)
	handleError(err, "Failed to declare exchange")

	//
	routingKey := randomRoutingKey()
	// create message
	task := config.Task{Number1: rand.Intn(100), Number2: rand.Intn(999)}
	body, err := json.Marshal(task)
	if err != nil {
		handleError(err, "Error encoding JSON")
	}

	err = ch.Publish(
		"logs_topic",
		"black.",
		false,
		false,
		amqp.Publishing{
			ContentType: "json/application",
			Body:        body,
		},
	)

	handleError(err, "Can't send message")
	log.Printf("routingKey: %s , AddTask: %d+%d", routingKey, task.Number1, task.Number2)
}

func randomRoutingKey() string {
	num := rand.Intn(2)
	s := ""
	switch {
	case num == 0:
		s = "red"
		break
	case num == 1:
		s = "black"
		break
	case num == 2:
		s = "yellow"
		break
	}

	return s
}
