package main

import (
	"fmt"

	"github.com/angusgmorrison/gophercises/blackjack_ai/blackjack"
	"github.com/angusgmorrison/gophercises/deck"
)

func main() {
	opts := blackjack.Options{NHands: 2}
	game := blackjack.New(opts)
	winnings := game.Play(basicAI{})
	fmt.Println(winnings)
}

type basicAI struct{}

func (ai basicAI) Bet(shuffled bool) int {
	return 100
}

func (ai basicAI) Play(hand []deck.Card, dealer deck.Card) blackjack.Move {
	score := blackjack.Score(hand...)
	if len(hand) == 2 {
		if score == 10 || score == 11 && !blackjack.Soft(hand...) {
			return blackjack.MoveDouble
		}
	}

	dScore := blackjack.Score(dealer)
	if dScore >= 5 && dScore <= 6 {
		return blackjack.MoveStand
	}
	if score < 13 {
		return blackjack.MoveHit
	}
	return blackjack.MoveStand
}

func (ai basicAI) Outcome(hand [][]deck.Card, dealer []deck.Card) {
	// noop
}
