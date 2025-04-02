package mq

import (
	"fmt"
	"log"
	"room/internal/env"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func ConnectToMq() *amqp.Connection {
	user := env.GetMqUser()
	pass := env.GetMqPassword()

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@viksync_mq:5672/", user, pass))
	failOnError(err, "Failed to connect to RabbitMQ")

	return conn
}
