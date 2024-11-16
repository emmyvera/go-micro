package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declearExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type
		true,         // durable?
		false,        // auto-delete
		false,        // internal?
		false,        // no-wait?
		nil,          // arguement
	)
}

func declearQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused?
		true,  // exclusive
		false, // no-wait
		nil,   // arguement

	)
}
