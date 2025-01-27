package pubsub

import (
	"bytes"
	"encoding/gob"
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
	decoder := func(data []byte) (T, error) {
		var target T
		err := json.Unmarshal(data, &target)
		return target, err
	}

	return subscribe(
		conn,
		exchange,
		queueName,
		key,
		simpleQueueType,
		handler,
		decoder,
	)
}

func SubscribeGob[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) AckType,
) error {
	decoder := func(data []byte) (T, error) {
		buf := bytes.NewBuffer(data)
		dec := gob.NewDecoder(buf)
		var target T

		err := dec.Decode(&target)
		return target, err
	}

	return subscribe(
		conn,
		exchange,
		queueName,
		key,
		simpleQueueType,
		handler,
		decoder,
	)
}

func subscribe[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) AckType,
	unmarshaler func([]byte) (T, error),
) error {
	declaredCh, queue, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
	if err != nil {
		return fmt.Errorf("error declaring and binding queue: %v", err)
	}

	declaredCh.Qos(10, 0, false)

	deliveryCh, err := declaredCh.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("error consuming messages: %v", err)
	}

	go func() {
		defer declaredCh.Close()
		for message := range deliveryCh {
			var target T
			target, err := unmarshaler(message.Body)
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
