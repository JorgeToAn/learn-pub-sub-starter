package main

import (
	"testing"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func TestHandlerPause(t *testing.T) {
	gs := gamelogic.NewGameState("testuser")
	handler := handlerPause(gs)

	// Pause the game
	acktype := handler(routing.PlayingState{IsPaused: true})
	if !gs.Paused {
		t.Error("expected game to be paused")
	}
	if acktype != pubsub.Ack {
		t.Error("expected AckType to be Ack")
	}

	// Resume the game
	acktype = handler(routing.PlayingState{IsPaused: false})
	if gs.Paused {
		t.Error("expected game to be resumed")
	}
	if acktype != pubsub.Ack {
		t.Error("expected AckType to be Ack")
	}
}
