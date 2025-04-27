package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Не удалось подключиться к RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Не удалось создать канал: %v", err)
	}
	defer ch.Close()

	queues := []string{"user_registered", "order_created"}

	for _, queue := range queues {
		_, err := ch.QueueDeclare(
			queue,
			true,  // durable
			false, // autoDelete
			false, // exclusive
			false, // noWait
			nil,   // arguments
		)
		if err != nil {
			log.Fatalf("Ошибка создания/подписки на очередь %s: %v", queue, err)
		}
	}

	msgsUser, _ := ch.Consume(
		"user_registered",
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)

	msgsOrder, _ := ch.Consume(
		"order_created",
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)

	go func() {
		for d := range msgsUser {
			var data map[string]string
			_ = json.Unmarshal(d.Body, &data)
			fmt.Printf("[UserRegistered] Новый пользователь: Email = %s, UserID = %s\n", data["email"], data["user_id"])
		}
	}()

	go func() {
		for d := range msgsOrder {
			var data map[string]string
			_ = json.Unmarshal(d.Body, &data)
			fmt.Printf("[OrderCreated] Новый заказ: OrderID = %s, UserID = %s\n", data["order_id"], data["user_id"])
		}
	}()

	log.Println(" [*] Ожидание сообщений. Для выхода нажмите CTRL+C")
	<-forever
}
