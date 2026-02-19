package main

import (
	"fmt"

	"github.com/RevoTale/turn-based-game-engine/examples/tic_tac_toe/tictactoe"
)

func main() {
	game, err := tictactoe.NewGame()
	if err != nil {
		panic(err)
	}

	for _, index := range []int{0, 3, 1, 4, 2} {
		if err := game.Play(index); err != nil {
			fmt.Println("play failed:", err)
			break
		}
	}

	fmt.Println("Board:")
	fmt.Println(game.BoardString())
	fmt.Println("Log:")
	for _, entry := range game.Logs() {
		fmt.Println("-", entry)
	}
}
