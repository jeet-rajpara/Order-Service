package amqp_helper

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return conn, ch, nil
}

func DeclareExchange(ch *amqp.Channel, name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	return ch.ExchangeDeclare(
		name,       // name
		kind,       // type
		durable,    // durable
		autoDelete, // auto-deleted
		internal,   // internal
		noWait,     // no-wait
		args,       // arguments
	)
}

func DeclareQueue(ch *amqp.Channel, name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return ch.QueueDeclare(
		name,       // name
		durable,    // durable
		autoDelete, // delete when unused
		exclusive,  // exclusive
		noWait,     // no-wait
		args,       // arguments
	)
}

func BindQueue(ch *amqp.Channel, queueName, routingKey, exchangeName string, noWait bool, args amqp.Table) error {
	return ch.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		noWait,
		args,
	)
}
func PublishMessage(ch *amqp.Channel, exchangeName string, body interface{}) error {

	// Declare the exchange
	err := ch.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}
	q, err := ch.QueueDeclare(
		"notificationQueue", // name
		false,               // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		map[string]interface{}{
			"x-max-length":              3,
			"x-dead-letter-exchange":    "my-dlx",
			"x-dead-letter-routing-key": "my-routing-key",
		},
		// nil, // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}
	err = ch.QueueBind(
		q.Name,           // queue name
		"my-routing-key", // routing key
		exchangeName,     // exchange
		false,            // no-wait
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind queue: %v", err)
	}
	messageBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = ch.Publish(
		exchangeName,     // exchange
		"my-routing-key", // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageBody,
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	return nil
}
