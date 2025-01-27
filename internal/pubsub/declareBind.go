package pubsub

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("error opening channel: %v", err)
	}

	// Set vals depending on queue type
	paramDurable := false
	paramAutodelete := false
	paramExclusive := false
	switch simpleQueueType {
	case SimpleQueueDurable:
		paramDurable = true
	case SimpleQueueTransient:
		paramAutodelete = true
		paramExclusive = true
	}

	queue, err := ch.QueueDeclare(
		queueName,
		paramDurable,
		paramAutodelete,
		paramExclusive,
		false,
		amqp.Table{
			"x-dead-letter-exchange": routing.ExchangePerilDeadLetter,
		},
	)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("error declaring queue: %v", err)
	}

	err = ch.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("error binding queue to channel: %v", err)
	}
	return ch, queue, nil
}
