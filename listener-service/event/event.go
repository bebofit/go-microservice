package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", //name
		"topic",      //type
		true,         // is it durable
		false,        // auto deleted
		false,        // internal
		false,        // no weight
		nil,          // arguments
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // name,
		false, // durable
		false, // do not delete if unused
		true,  // exclusive
		false, // no wait
		nil,   // argument
	)
}
