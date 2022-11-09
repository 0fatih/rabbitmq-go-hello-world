package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	exchangeName := "delayed-exchange"

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "failed to connect to rabbitmq")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "failed to open a channel")
	defer ch.Close()

	args := make(amqp.Table)
	args["x-delayed-type"] = "direct"

	err = ch.ExchangeDeclare(
		exchangeName,
		"x-delayed-message",
		true,
		false,
		false,
		false,
		args,
	)
	failOnError(err, "delayed exchange can't define")

	headers := make(amqp.Table)
	headers["x-delay"] = 3000

	body := "Sorry I am late!"
	err = ch.PublishWithContext(
		context.Background(),
		exchangeName,
		"hello",
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			ContentType:  "text/plain",
			Body:         []byte(body),
			Headers:      headers,
		},
	)

	failOnError(err, "failed to publish a message")
	log.Printf("[x] Sent %s\n", body)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
