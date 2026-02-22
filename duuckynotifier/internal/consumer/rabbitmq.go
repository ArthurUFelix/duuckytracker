package consumer

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/arthurufelix/duuckynotifier/internal/notifier"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn *amqp.Connection
	channel *amqp.Channel
	queueName string
	discordNotifier *notifier.DiscordNotifier
}

func NewConsumer(rabbitmqURL, queueName string, discord *notifier.DiscordNotifier) (*Consumer, error) {
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel:: %w", err)
	}

	_, err = channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	err = channel.Qos(
		1, 0, false,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	log.Println("Connected to RabbitMQ")

	return &Consumer{
		conn: conn,
		channel: channel,
		queueName: queueName,
		discordNotifier: discord,
	}, nil
}

func (c *Consumer) Start() error {
	messages, err := c.channel.Consume(
		c.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Printf("Started consuming from queue: %s", c.queueName)

	go func() {
		for msg := range messages {
			c.handleMessage(msg)
		}
	}()

	return nil
}

func (c *Consumer) handleMessage(msg amqp.Delivery) {
	log.Printf("Received message: %s", string(msg.Body))

	var rabbitMsg notifier.RabbitMQMessage
	err := json.Unmarshal(msg.Body, &rabbitMsg)
	if err != nil {
		log.Printf("failed to parse message: %v", err)
		msg.Nack(false, false)
		return
	}

	err = c.discordNotifier.SendNotification(rabbitMsg)
	if err != nil {
		log.Printf("failed to send Discord notification: %v", err)
		msg.Nack(false, true)
		return
	}

	msg.Ack(false)
}

func (c *Consumer) Close() error {
	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			return err
		}
	}
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return err
		}
	}
	log.Println("Disconnected from RabbitMQ")
	return nil
}