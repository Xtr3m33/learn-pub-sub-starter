package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/messageutils"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

const connectionString = "amqp://guest:guest@localhost:5672/"

func main() {
	fmt.Println("Starting Peril server...")

	// Open connection to rabbitmq server
	connection, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatalf("error conencting to rabbitmq server: %v", err)
	}
	defer connection.Close()
	log.Print("Successfully connected to rabbitmq server")

	msgChannel, err := connection.Channel()
	if err != nil {
		log.Fatalf("error opening channel for rabbitmq server: %v", err)
	}

	_, _, err = pubsub.DeclareAndBind(
		connection,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		fmt.Sprintf("%s.*", routing.GameLogSlug),
		pubsub.Durable,
	)
	if err != nil {
		log.Fatalf("error declaring game log queue: %v", err)
	}

	gamelogic.PrintServerHelp()

	// Main loop
	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}

		switch strings.ToLower(input[0]) {
		case "pause":
			log.Print("pausing game")
			err = messageutils.SendPauseMessage(msgChannel, true)
			if err != nil {
				log.Fatalf("error pausing game: %v", err)
			}
		case "resume":
			log.Print("resuming game")
			err = messageutils.SendPauseMessage(msgChannel, false)
			if err != nil {
				log.Fatalf("error unpausing game: %v", err)
			}
		case "quit":
			log.Print("quitting game")
			fmt.Println("Shutting down Peril server...")
			os.Exit(0)
		default:
			fmt.Println("unknown command")
		}
	}

}
