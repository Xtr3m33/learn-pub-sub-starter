package messageutils

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func SendPauseMessage(ch *amqp.Channel, isPaused bool) error {
	err := pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: isPaused})
	if err != nil {
		return fmt.Errorf("error sending message on channel %v: %v", ch, err)
	}
	return nil
}
