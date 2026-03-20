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
	log.Println("Connected to RabbitMQ")

	publishCh, err := connection.Channel()
	if err != nil {
		log.Fatalf("failed to create RabbitMQ channel: %v", err)
	}

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("%v", err)
	}

	gs := gamelogic.NewGameState(username)

	err = pubsub.SubscribeJSON(
		connection,                            // conn
		routing.ExchangePerilDirect,           // exchange
		routing.PauseKey+"."+gs.GetUsername(), // queueName
		routing.PauseKey,                      // key
		pubsub.SimpleQueueTransient,           // queueType
		handlerPause(gs),                      // handler
	)
	if err != nil {
		log.Fatalf("failed to subscribe to pause messages: %v", err)
	}

	err = pubsub.SubscribeJSON(
		connection,                 // conn
		routing.ExchangePerilTopic, // exchange
		routing.ArmyMovesPrefix+"."+gs.GetUsername(), // queueName
		routing.ArmyMovesPrefix+".*",                 // key
		pubsub.SimpleQueueTransient,                  // queueType
		handlerMove(publishCh, gs),                   // handler
	)
	if err != nil {
		log.Fatalf("failed to subscribe to move messages: %v", err)
	}

	err = pubsub.SubscribeJSON(
		connection,                         // conn
		routing.ExchangePerilTopic,         // exchange
		routing.WarRecognitionsPrefix,      // queue name
		routing.WarRecognitionsPrefix+".*", // routing key
		pubsub.SimpleQueueDurable,          // queue type
		handlerWar(gs),                     // handler
	)

	running := true
	for running {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		command := words[0]
		switch command {
		case "move":
			move, err := gs.CommandMove(words)
			if err != nil {
				fmt.Println(err)
				continue
			}

			err = pubsub.PublishJSON(
				publishCh,                  // channel
				routing.ExchangePerilTopic, // exchange
				routing.ArmyMovesPrefix+"."+gs.GetUsername(), // key
				move, // val
			)
			if err != nil {
				log.Printf("failed to publish move: %v", err)
				continue
			}
			log.Printf("Moved %v unit(s) to %s\n", len(move.Units), move.ToLocation)
		case "spawn":
			err := gs.CommandSpawn(words)
			if err != nil {
				fmt.Println(err)
			}
		case "status":
			gs.CommandStatus()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "help":
			gamelogic.PrintClientHelp()
		case "quit":
			gamelogic.PrintQuit()
			running = false
		default:
			fmt.Println("Unknown command")
		}
	}

	fmt.Println()
	fmt.Println("Closing Peril client...")
}
