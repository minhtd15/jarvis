package rabbitmq

import (
	"context"
	"education-website/api/response/qlda"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
)

func RabbitMQPublisher(fileData []byte, ctx context.Context, fileType string) (*qlda.AutoGenerated, error) {
	fmt.Println("RabbitMQPublisher")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	fmt.Println("Successfully connected to RabbitMQ")

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var q amqp.Queue
	if fileType == "application/pdf" {
		q, err = ch.QueueDeclare(
			"pdf",
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return nil, err
		}
	} else {
		q, err = ch.QueueDeclare(
			"docx",
			false,
			false,
			true,
			false,
			nil,
		)
		if err != nil {
			return nil, err
		}
	}

	msgs, err := ch.Consume(
		"receive", // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)

	if err != nil {
		return nil, err
	}
	//fmt.Println(q)

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        fileData,
			ReplyTo:     "receive",
		},
	)

	for d := range msgs {
		var data qlda.AutoGenerated
		err := json.Unmarshal(d.Body, &data)
		if err != nil {
			fmt.Println("Error decoding JSON: %s", err)
		}
		fmt.Printf("Decoded JSON: %v\n", data)
		//if msg, ok := data["message"]; ok && msg == "EOF" {
		//	break
		//}
		if data.Message == "EOF" || data.Message == "eof" {
			break
		}
		return &data, nil
	}

	fmt.Println("Successfully published message to RabbitMQ")
	return nil, err
}
