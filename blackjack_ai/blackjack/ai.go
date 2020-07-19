package blackjack

import (
	"fmt"

	"github.com/angusgmorrison/gophercises/deck"
)

// AI specifies the methods required for a player of a blackjack game.
type AI interface {
	Bet(shuffled bool) int
	Play(hand []deck.Card, dealer deck.Card) Move
	Outcome(hand [][]deck.Card, dealer []deck.Card)
}

// dealerAI is the default implentation of the blackjack dealer.
type dealerAI struct{}

func (ai dealerAI) Bet(shuffled bool) int {
	// noop
	return 1
}

func (ai dealerAI) Play(hand []deck.Card, dealer deck.Card) Move {
	score := Score(hand...)
	if score <= 16 || score == 17 && Soft(hand...) {
		return MoveHit
	}
	return MoveStand
}

func (ai dealerAI) Outcome(hand [][]deck.Card, dealer []deck.Card) {
	// noop
}

// HumanAI conceals the implementation of a default human player
// for a blackjack game.
func HumanAI() AI {
	return humanAI{}
}

type humanAI struct{}

func (ai humanAI) Bet(shuffled bool) int {
	if shuffled {
		fmt.Println("Deck was just shuffled...")
	}
	fmt.Println("What would you like to bet? (min. 100)")
	var bet int
	fmt.Scanf("%d\n", &bet)
	return bet
}

// Accepted player inputs.
const (
	double = "d"
	hit    = "h"
	stand  = "s"
)

func (ai humanAI) Play(hand []deck.Card, dealer deck.Card) Move {
	for {
		var input string
		fmt.Println("AI:", hand)
		fmt.Println("Dealer:", dealer)
		fmt.Println("What will you do? (h)it, (s)tand, (d)ouble")
		fmt.Scanf("%s\n", &input)

		switch input {
		case double:
			return MoveDouble
		case hit:
			return MoveHit
		case stand:
			return MoveStand
		default:
			fmt.Println("Command not recognised: enter (h)it, (s)tand or (d)ouble")
		}
	}
}

func (ai humanAI) Outcome(hand [][]deck.Card, dealer []deck.Card) {
	fmt.Println("==FINAL HANDS==")
	fmt.Println("AI:", hand)
	fmt.Println("Dealer:", dealer)
	fmt.Println()
	// fmt.Printf("AI: %s\nScore: %d\n", ret.AI, pScore)
	// fmt.Printf("AI: %s\nScore: %d\n", ret.Dealer, dScore)
}
