//go:generate stringer -type=Suit,Rank

package deck

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// Suit represents the suit of a card, including Jokers
type Suit uint8

const (
	Spades Suit = iota
	Diamonds
	Clubs
	Hearts
	Joker // special case
)

var suits = [...]Suit{Spades, Diamonds, Clubs, Hearts}

// Rank represents the numeric value of a card
type Rank uint8

const (
	_ Rank = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

const (
	minRank = Ace
	maxRank = King
)

// Card represents a playing card with both suit and rank.
type Card struct {
	Rank
	Suit
}

const cardsInDeck = 52

func (c Card) String() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}
	return fmt.Sprintf("%s of %s", c.Rank, c.Suit)
}

// An Option is a fuctional option provided to New to modify the deck returned.
type Option func([]Card) []Card

// New returns a new deck of cards as a slice with len cardsInDeck.
func New(opts ...Option) []Card {
	cards := make([]Card, 0, cardsInDeck)
	for _, suit := range suits {
		for rank := minRank; rank <= maxRank; rank++ {
			cards = append(cards, Card{Suit: suit, Rank: rank})
		}
	}
	for _, opt := range opts {
		cards = opt(cards)
	}
	return cards
}

// DefaultSort sorts cards in ascending order by absRank, from the Ace of Spades to the Jokers.
func DefaultSort(cards []Card) []Card {
	sort.Slice(cards, Less(cards))
	return cards
}

// Less returns the default less func for a deck of cards, comparing the absRank of each pair.
func Less(cards []Card) func(i, j int) bool {
	return func(i, j int) bool {
		return cards[i].absRank() < cards[j].absRank()
	}
}

// Sort returns an Option used by new to apply a custom sort as defined by less to a new deck of
// cards.
func SortOption(less func(cards []Card) func(i, j int) bool) Option {
	return func(cards []Card) []Card {
		sort.Slice(cards, less(cards))
		return cards
	}
}

func (c Card) absRank() int {
	return int(c.Suit)*int(maxRank) + int(c.Rank)
}

var shuffleRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// Shuffle is an Option returning a randomly shuffled deck of cards using the Fisher-Yates shuffle.
func Shuffle(cards []Card) []Card {
	for i := len(cards) - 1; i > 0; i-- {
		swapTo := shuffleRand.Intn(i + 1)
		cards[i], cards[swapTo] = cards[swapTo], cards[i]
	}
	// Returning cards despite the in-place change allows Shuffle to work as an Option.
	return cards
}

// Jokers returns an Option that adds n ranked Jokers to a deck.
func Jokers(n int) Option {
	return func(cards []Card) []Card {
		for i := 0; i < n; i++ {
			cards = append(cards, Card{Rank: Rank(i), Suit: Joker})
		}
		return cards
	}
}

// Filter accepts a predicate function and returns an Option that filters a deck, returning a new
// deck contain all cards for which the filter returns false. I.e. those that were not filtered out.
func Filter(f func(card Card) bool) Option {
	return func(cards []Card) []Card {
		var ret []Card
		for _, card := range cards {
			if !f(card) {
				ret = append(ret, card)
			}
		}
		return ret
	}
}

func Deck(n int) Option {
	return func(cards []Card) []Card {
		var ret []Card
		for i := 0; i < n; i++ {
			ret = append(ret, cards...)
		}
		return ret
	}
}
