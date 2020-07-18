package main

import (
	"fmt"
	"strings"

	"github.com/angusgmorrison/gophercises/deck"
)

type Hand []deck.Card

const (
	stand = "s"
	hit   = "h"
)

func main() {
	var gs GameState
	gs = Shuffle(gs)
	gs = Deal(gs)

	var input string
	// Ensure the most recent GameState is accessed on each iteration.
	for current := gs; current.Phase == PlayerTurn; current = gs {
		fmt.Println("Player:", gs.Player)
		fmt.Println("Dealer:", gs.Dealer.DealerString())
		fmt.Println("What will you do? (h)it, (s)tand")
		fmt.Scanf("%s\n", &input)
		switch input {
		case hit:
			gs = Hit(gs)
		case stand:
			gs = Stand(gs)
		default:
			fmt.Println("Command not recognised: enter (h)it or (s)tand")
		}
	}

	for gs.Phase == DealerTurn {
		if shouldHit(gs.Dealer.Score(), gs.Dealer.MinScore()) {
			gs = Hit(gs)
		} else {
			gs = Stand(gs)
		}
	}

	gs = EndHand(gs)
}

func shouldHit(score, minScore int) bool {
	return score <= 16 || score == 17 && minScore != 17
}

func (h Hand) DealerString() string {
	return h[0].String() + ", **HIDDEN**"
}

func (h Hand) Score() int {
	minScore := h.MinScore()
	if minScore > 11 {
		return minScore
	}
	for _, c := range h {
		if c.Rank == deck.Ace {
			// ace is currently worth 1, and we are changing it to be worth 11
			return minScore - 1 + 11
		}
	}
	return minScore
}

func (h Hand) MinScore() int {
	var score int
	for _, c := range h {
		score += min(int(c.Rank), 10)
	}
	return score
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (h Hand) String() string {
	strs := make([]string, len(h))
	for i, card := range h {
		strs[i] = card.String()
	}
	return strings.Join(strs, ", ")
}
