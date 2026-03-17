package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
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
	fmt.Println("Connected to RabbitMQ")

	publishCh, err := connection.Channel()
	if err != nil {
		log.Fatalf("failed to create RabbitMQ channel: %v", err)
	}

	running := true
	gamelogic.PrintServerHelp()
	for running {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}

		command := input[0]
		switch command {
		case "pause":
			log.Println("sending pause messsage")
			err = pubsub.PublishJSON(
				publishCh,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{
					IsPaused: true,
				},
			)
			if err != nil {
				log.Printf("failed to publish pause messsage: %v\n", err)
			}
		case "resume":
			log.Println("sending resume message")
			err = pubsub.PublishJSON(
				publishCh,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{
					IsPaused: false,
				},
			)
			if err != nil {
				log.Printf("failed to publish resume message: %v\n", err)
			}
		case "quit":
			log.Println("closing server")
			running = false
		case "help":
			gamelogic.PrintServerHelp()
		default:
			log.Printf("unknown command '%s'", command)
		}
	}

	fmt.Println("Closing Peril server...")
}
