package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.AckType {
	return func(ps routing.PlayingState) pubsub.AckType {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
		return pubsub.AckTypeAck
	}
}

func handlerMove(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.ArmyMove) pubsub.AckType {
	return func(move gamelogic.ArmyMove) pubsub.AckType {
		defer fmt.Print("> ")
		outcome := gs.HandleMove(move)
		switch outcome {
		case gamelogic.MoveOutcomeMakeWar:
			err := pubsub.PublishJSON(
				ch,
				routing.ExchangePerilTopic,
				fmt.Sprintf("%s.%s", routing.WarRecognitionsPrefix, gs.GetUsername()),
				gamelogic.RecognitionOfWar{
					Attacker: move.Player,
					Defender: gs.GetPlayerSnap(),
				},
			)
			if err != nil {
				log.Printf("error acknowledging war: %v", err)
				return pubsub.AckTypeNackRequeue
			}
			return pubsub.AckTypeAck
		case gamelogic.MoveOutComeSafe:
			return pubsub.AckTypeAck
		case gamelogic.MoveOutcomeSamePlayer:
			fallthrough
		default:
			return pubsub.AckTypeNackDiscard
		}
	}
}

func handlerWar(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.RecognitionOfWar) pubsub.AckType {
	return func(rw gamelogic.RecognitionOfWar) pubsub.AckType {
		defer fmt.Print("> ")
		outcome, winner, loser := gs.HandleWar(rw)

		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.AckTypeNackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.AckTypeNackDiscard
		case gamelogic.WarOutcomeOpponentWon, gamelogic.WarOutcomeYouWon, gamelogic.WarOutcomeDraw:
			var msg string
			switch outcome {
			case gamelogic.WarOutcomeOpponentWon:
				msg = "%s won a war against %s"
			case gamelogic.WarOutcomeYouWon:
				msg = "%s won a war against %s"
			case gamelogic.WarOutcomeDraw:
				msg = "A war between %s and %s resulted in a draw"
			}
			err := pubsub.PublishGob(
				ch,
				routing.ExchangePerilTopic,
				fmt.Sprintf("%s.%s", routing.GameLogSlug, gs.GetUsername()),
				routing.GameLog{
					CurrentTime: time.Now().UTC(),
					Message:     fmt.Sprintf(msg, winner, loser),
					Username:    gs.GetUsername(),
				},
			)
			if err != nil {
				log.Printf("error publishing logs: %v", err)
				return pubsub.AckTypeNackRequeue
			}
			return pubsub.AckTypeAck
		default:
			log.Printf("unknown war outcome: %v", outcome)
			return pubsub.AckTypeNackDiscard
		}
	}
}
