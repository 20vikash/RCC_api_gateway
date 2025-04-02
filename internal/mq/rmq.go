package mq

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MQ struct {
	User string
	Pass string
	Port string
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func (mq *MQ) ConnectToMq() *amqp.Channel {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@viksync_mq:%s/", mq.User, mq.Pass, mq.Port))
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	return ch
}
