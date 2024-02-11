package main

import (
	"fmt"
	"github.com/andre2ar/go-events/pkg/rabbitmq"
)

func main() {
	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		"orders",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	err = ch.QueueBind(queue.Name, "", "amq.direct", false, nil)
	if err != nil {
		panic(err)
	}

	rabbitmq.Publish(ch, "Hello World!", "amq.direct")
	fmt.Println("Message published")
}
