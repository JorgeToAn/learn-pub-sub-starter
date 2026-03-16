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

func main() {
	const rabbitConnUrl = "amqp://guest:guest@localhost:5672/"

	fmt.Println("Starting Peril server...")
	connection, err := amqp.Dial(rabbitConnUrl)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer connection.Close()
	fmt.Println("Connection successful")

	publishCh, err := connection.Channel()
	if err != nil {
		log.Fatalf("failed to create RabbitMQ channel: %v", err)
	}
	err = pubsub.PublishJSON(
		publishCh,
		routing.ExchangePerilDirect,
		routing.PauseKey,
		routing.PlayingState{
			IsPaused: true,
		},
	)
	if err != nil {
		fmt.Printf("failed to publish messsage: %v\n", err)
	}
	fmt.Println("Published Pause message!")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println()
	fmt.Println("Closing Peril server...")
}
