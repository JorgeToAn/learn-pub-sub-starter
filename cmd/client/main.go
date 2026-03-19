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

	fmt.Println("Starting Peril client...")
	connection, err := amqp.Dial(rabbitConnUrl)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer connection.Close()
	fmt.Println("Connected to RabbitMQ")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("%v", err)
	}

	gameState := gamelogic.NewGameState(username)

	err = pubsub.SubscribeJSON(
		connection,                  // conn
		routing.ExchangePerilDirect, // exchange
		fmt.Sprintf("%s.%s", routing.PauseKey, username), // queueName
		routing.PauseKey,            // key
		pubsub.SimpleQueueTransient, // queueType
		handlerPause(gameState),     // handler
	)
	if err != nil {
		log.Fatalf("failed to subscribe to pause queue: %v", err)
	}

	running := true
	for running {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		command := words[0]
		switch command {
		case "move":
			_, err := gameState.CommandMove(words)
			if err != nil {
				fmt.Println(err)
			}
		case "spawn":
			err := gameState.CommandSpawn(words)
			if err != nil {
				fmt.Println(err)
			}
		case "status":
			gameState.CommandStatus()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "help":
			gamelogic.PrintClientHelp()
		case "quit":
			gamelogic.PrintQuit()
			running = false
		default:
			fmt.Printf("Unknown command")
		}
	}

	fmt.Println()
	fmt.Println("Closing Peril client...")
}
