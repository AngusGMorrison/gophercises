package main

import (
	"fmt"

	"github.com/angusgmorrison/gophercises/blackjack_ai/blackjack"
)

func main() {
	opts := blackjack.Options{NHands: 2}
	game := blackjack.New(opts)
	winnings := game.Play(blackjack.HumanAI())
	fmt.Println(winnings)
}
