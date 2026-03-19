package main

import (
	"testing"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func TestHandlerPause(t *testing.T) {
	gs := gamelogic.NewGameState("testuser")
	handler := handlerPause(gs)

	// Pause the game
	handler(routing.PlayingState{IsPaused: true})
	if !gs.Paused {
		t.Error("expected game to be paused")
	}

	// Resume the game
	handler(routing.PlayingState{IsPaused: false})
	if gs.Paused {
		t.Error("expected game to be resumed")
	}
}
