package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

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
	fmt.Println("Successfully connected to rabbitmq server")

	msgChannel, err := connection.Channel()
	if err != nil {
		log.Fatalf("error opening channel for rabbitmq server: %v", err)
	}

	//TODO - Testing
	err = pubsub.PublishJSON(
		msgChannel,
		routing.ExchangePerilDirect,
		routing.PauseKey,
		routing.PlayingState{
			IsPaused: true,
		},
	)
	if err != nil {
		log.Fatalf("error sending message on channel: %v", err)
	}

	// Wait for ctrl+c (exit)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Shutting down Peril server...")
	os.Exit(0)
}
