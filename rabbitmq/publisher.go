package rabbitmq

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
)

func RabbitMQPublisher(fileData []byte, ctx context.Context) error {
	fmt.Println("RabbitMQPublisher")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}

	defer conn.Close()

	fmt.Println("Successfully connected to RabbitMQ")

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer conn.Close()

	q, err := ch.QueueDeclare(
		"queue",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}
	fmt.Println(q)

	err = ch.Publish(
		"",        // Chuyển empty string để gửi message trực tiếp vào queue mà không thông qua exchange
		"myQueue", // Tên của hàng đợi
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        fileData,
		},
	)
	fmt.Println("Successfully published message to RabbitMQ")
	return nil
}
