package main

import (
	"fmt"
	"github.com/andre2ar/go-events/pkg/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	rabbitmqChannel, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer rabbitmqChannel.Close()

	messagesChannel := make(chan amqp.Delivery)

	go rabbitmq.Consume(rabbitmqChannel, messagesChannel, "orders")

	fmt.Println("Waiting for messages...")
	for msg := range messagesChannel {
		fmt.Println(string(msg.Body))
		msg.Ack(false)
	}
}
