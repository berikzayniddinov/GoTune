package events

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type EventPublisher struct {
	channel *amqp.Channel
}

func NewEventPublish(amqpURL string) *EventPublisher {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatalf("Не удалось подключиться к RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Не удалось создать канал RabbitMQ: %v", err)
	}

	return &EventPublisher{channel: ch}
}

func (p *EventPublisher) Publish(queue string, event any) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = p.channel.QueueDeclare(
		queue,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	return p.channel.Publish(
		"",    // default exchange
		queue, // routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
