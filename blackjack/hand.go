package main

import (
	"strings"

	"github.com/angusgmorrison/gophercises/deck"
)

// Hand represents a single player's hand of cards.
type Hand []deck.Card

// DealerString returns a Hand as a string with the second card
// hidden.
func (h Hand) DealerString() string {
	return h[0].String() + ", **HIDDEN**"
}

// Score returns the point value of the hand, converting aces between
// 11 and 1 as appropriate.
func (h Hand) Score() int {
	minScore := h.MinScore()
	if minScore > 11 {
		return minScore
	}
	for _, c := range h {
		if c.Rank == deck.Ace {
			// Ace is currently worth 1, and we are changing it to be worth 11.
			return minScore - 1 + 11
		}
	}
	return minScore
}

// MinScore returns the point value of the hand, with any Aces counted
// as 1.
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
