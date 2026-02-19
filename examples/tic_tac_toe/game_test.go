package main

import (
	"errors"
	"testing"

	"github.com/RevoTale/turn-based-game-engine/examples/tic_tac_toe/tictactoe"
	"github.com/RevoTale/turn-based-game-engine/turnbased"
)

func TestGameWinnerPath(t *testing.T) {
	game, err := tictactoe.NewGame()
	if err != nil {
		t.Fatalf("NewGame failed: %v", err)
	}

	for _, index := range []int{0, 3, 1, 4, 2} {
		if err := game.Play(index); err != nil {
			t.Fatalf("Play(%d) failed: %v", index, err)
		}
	}

	_, winner, result := game.Snapshot()
	if winner != tictactoe.PlayerX {
		t.Fatalf("expected winner X, got %c", winner)
	}
	if result != turnbased.MatchResultWinner {
		t.Fatalf("expected winner result, got %v", result)
	}

	if err := game.Play(5); !errors.Is(err, tictactoe.ErrGameFinished) {
		t.Fatalf("expected ErrGameFinished, got %v", err)
	}
}

func TestGameDrawPath(t *testing.T) {
	game, err := tictactoe.NewGame()
	if err != nil {
		t.Fatalf("NewGame failed: %v", err)
	}

	for _, index := range []int{0, 1, 2, 4, 3, 5, 7, 6, 8} {
		if err := game.Play(index); err != nil {
			t.Fatalf("Play(%d) failed: %v", index, err)
		}
	}

	_, winner, result := game.Snapshot()
	if winner != 0 {
		t.Fatalf("expected no winner value, got %c", winner)
	}
	if result != turnbased.MatchResultDraw {
		t.Fatalf("expected draw result, got %v", result)
	}
	if err := game.Play(0); !errors.Is(err, tictactoe.ErrGameFinished) {
		t.Fatalf("expected ErrGameFinished, got %v", err)
	}
}

func TestGameInvalidMove(t *testing.T) {
	game, err := tictactoe.NewGame()
	if err != nil {
		t.Fatalf("NewGame failed: %v", err)
	}

	err = game.Play(11)
	if !errors.Is(err, tictactoe.ErrMoveBounds) {
		t.Fatalf("expected ErrMoveBounds, got %v", err)
	}

	if err := game.Play(0); err != nil {
		t.Fatalf("first play failed: %v", err)
	}
	err = game.Play(0)
	if !errors.Is(err, tictactoe.ErrCellBusy) {
		t.Fatalf("expected ErrCellBusy, got %v", err)
	}
}

func TestNilGame(t *testing.T) {
	var game *tictactoe.Game
	err := game.Play(0)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, tictactoe.ErrNilGame) {
		t.Fatalf("unexpected error: %v", err)
	}
}
