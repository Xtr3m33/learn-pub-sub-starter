package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

const connectionString = "amqp://guest:guest@localhost:5672/"

func main() {
	fmt.Println("Starting Peril server...")

	connection, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatalf("error conencting to rabbitmq server: %v", err)
	}
	defer connection.Close()
	fmt.Println("Successfully connected to rabbitmq server")

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Shutting down Peril server...")
	os.Exit(0)
}
