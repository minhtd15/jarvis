package rabbitmq

import (
	"context"
	education_website "education-website"
	"education-website/client"
	"education-website/rabbitmq/response"
	"encoding/json"
	"fmt"
	_ "github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"

	"github.com/streadway/amqp"
)

type Message struct {
	// Định nghĩa các trường dữ liệu của message tương ứng với ứng dụng của bạn
	Content string `json:"content"`
}

func RabbitMqConsumer(redisClient client.RedisClient, classService education_website.ClassService) error {
	// Connect to RabbitMQ server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	//conn, err := amqp.Dial("amqp://guest:guest@34.80.130.47:5672/")
	//
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ server: %v", err)
	}
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare a queue
	q, err := ch.QueueDeclare(
		"queue", // Queue name
		true,    // Durable
		false,   // Delete when unused
		false,   // Exclusive
		false,   // No-wait
		nil,     // Arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %v", err)
	}

	// Consume messages from the queue
	msgs, err := ch.Consume(
		q.Name, // Queue name
		"",     // Consumer name
		true,   // Auto Acknowledge
		false,  // Exclusive
		false,  // No Local
		false,  // No Wait
		nil,    // Arguments
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %v", err)
	}

	// Create a channel to handle signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// Receive messages from the queue
	// Make a channel to receive messages into infinite loop.
	forever := make(chan bool)

	go func() {
		for message := range msgs {
			// For example, show received message in a console.
			log.Printf(" > Received message: %s\n", message.Body)
			err := handleMessageFromQueue(message.Body, redisClient, classService)
			if err != nil {
				log.Errorf("Error when handling message from queue: %v", err)
				return
			}
		}
	}()

	<-forever

	return nil
	//if err != nil {
	//	return err
	//}
	//
	//// Lặp qua các message và trả về message đầu tiên nhận được
	//for d := range msgs {
	//	var message Message
	//	err := json.Unmarshal(d.Body, &message)
	//	if err != nil {
	//		continue
	//	}
	//	log.Info("Received message: ", message.Content)
	//	return nil
	//}

	return nil
}

func handleMessageFromQueue(message []byte, redisClient client.RedisClient, classService education_website.ClassService) error {
	// Xử lý message nhận được từ hàng đợi
	log.Printf(" > Received message: %s\n", message)
	var rq response.YearlyResponse
	err := json.Unmarshal(message, &rq)
	if err != nil {
		log.WithError(err).Errorf("Error marshal request to fix course information")
		return err
	}
	stringValue := fmt.Sprintf("%.6f", rq.TotalYearlyRevenue)

	// save the year-revenue into redis
	err = redisClient.Save(rq.Year, stringValue, context.Background())
	if err != nil {
		log.WithError(err).Errorf("Error save yearly revenue to redis")
		return err
	}

	// save data into db
	err = classService.UpdateYearlyRevenueAndCourseRevenue(rq, context.Background())
	if err != nil {
		log.WithError(err).Errorf("Error save yearly revenue to db")
		return err
	}
	return err
}
