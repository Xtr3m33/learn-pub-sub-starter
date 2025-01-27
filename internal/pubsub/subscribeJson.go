package pubsub

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) AckType,
) error {
	declaredCh, queue, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
	if err != nil {
		return fmt.Errorf("error declaring and binding queue: %v", err)
	}

	deliveryCh, err := declaredCh.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("error consuming messages: %v", err)
	}

	go func() {
		defer declaredCh.Close()
		for message := range deliveryCh {
			var target T
			err := json.Unmarshal(message.Body, &target)
			if err != nil {
				log.Printf("error unmarshaling message: %v", err)
				continue
			}
			ack := handler(target)
			switch ack {
			case AckTypeAck:
				message.Ack(false)
			case AckTypeNackRequeue:
				message.Nack(false, true)
			case AckTypeNackDiscard:
				message.Nack(false, false)
			}
		}
	}()
	return nil
}
