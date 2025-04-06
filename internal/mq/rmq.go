package mq

import (
	"fmt"
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MQ struct {
	User string
	Pass string
	Port string
}

var (
	conn *amqp.Connection
	once sync.Once
	err  error
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func (mq *MQ) ConnectToMq() *amqp.Connection {
	once.Do(func() {
		conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@viksync_mq:%s/", mq.User, mq.Pass, mq.Port))
		failOnError(err, "Failed to connect to RabbitMQ")
	})

	return conn
}

func (mq *MQ) CreateChannel(con *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Scheduler: Failed to open channel: %v", err)
	}

	return ch
}
