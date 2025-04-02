package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type OutputResponse struct {
	Output string
}

type MQNode struct {
	Code     string
	Jid      string
	Language string
}

func createJobID(lang, rid, userName string) string {
	uniqueID := uuid.New().String()
	data := lang + rid + userName + uniqueID
	hash := sha256.Sum256([]byte(data))

	return hex.EncodeToString(hash[:])[:16]
}

func pushToMq(ctx context.Context, mq *MQNode, channel *amqp.Channel) {
	mqJson, err := json.Marshal(mq)
	if err != nil {
		log.Println(err)
	}

	q, err := channel.QueueDeclare(
		"code", // name
		true,   // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		log.Println(err)
	}

	err = channel.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         mqJson,
		})
	if err != nil {
		log.Println(err)
	}

	log.Println(" [x] Sent")
}

func (app *Application) outputCode(w http.ResponseWriter, r *http.Request) {
	language := r.URL.Query().Get("language")
	roomId := r.URL.Query().Get("id")
	userName := r.URL.Query().Get("username")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	jobID := createJobID(language, roomId, userName)

	var codeData Code

	var err error

	if err = json.NewDecoder(r.Body).Decode(&codeData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	code := codeData.Code

	mq := &MQNode{
		Code:     code,
		Jid:      jobID,
		Language: language,
	}

	pushToMq(ctx, mq, app.MqChannel)

	w.Write([]byte(jobID))
}
