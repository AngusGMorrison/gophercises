package blackjack

import (
	"fmt"

	"github.com/angusgmorrison/gophercises/deck"
)

// Options allow the user to configure a new game.
type Options struct {
	NDecks             int
	NHands             int
	BlackjackPayout    float64
	ReshuffleThreshold int // the fraction of the deck below which to reshuffle (3 == 1/3)
}

// Option defaults
const (
	defaultNDecks             = 3
	defaultNHands             = 100
	defaultBlackjackPayout    = 1.5
	defaultReshuffleThreshold = 3
)

// New starts a new game with the specified options.
func New(opts Options) Game {
	g := Game{
		phase:    playerTurn,
		dealerAI: dealerAI{},
	}

	if opts.NDecks == 0 {
		opts.NDecks = defaultNDecks
	}
	if opts.NHands == 0 {
		opts.NHands = defaultNHands
	}
	if opts.BlackjackPayout == 0.0 {
		opts.BlackjackPayout = defaultBlackjackPayout
	}
	if opts.ReshuffleThreshold == 0 {
		opts.ReshuffleThreshold = defaultReshuffleThreshold
	}

	g.nDecks = opts.NDecks
	g.nHands = opts.NHands
	g.blackjackPayout = opts.BlackjackPayout
	g.minCards = (52 * g.nDecks) / opts.ReshuffleThreshold

	return g
}

// Game holds the current, mutable state of the game.
type Game struct {
	nDecks          int
	nHands          int
	minCards        int
	deck            []deck.Card
	phase           phase
	player          []deck.Card
	dealer          []deck.Card
	dealerAI        AI
	balance         int
	blackjackPayout float64
}

// Phase represents the current stage of gameplay.
type phase uint8

const (
	betting phase = iota
	playerTurn
	dealerTurn
	handOver
)

// Play begins the game, taking in a player AI, the number of decks
// to play with and the number of rounds to play, and returning the
// player's final balance.
func (g *Game) Play(player AI) int {
	for i := 0; i < g.nHands; i++ {
		if len(g.deck) < g.minCards {
			g.deck = deck.New(deck.Deck(g.nDecks), deck.Shuffle)
		}

		deal(g)

		for g.phase == playerTurn {
			hand := make([]deck.Card, len(g.player))
			copy(hand, g.player)
			move := player.Play(hand, g.dealer[0])
			move(g)
		}

		for g.phase == dealerTurn {
			hand := make([]deck.Card, len(g.dealer))
			copy(hand, g.dealer)
			move := g.dealerAI.Play(hand, hand[0])
			move(g)
		}

		endHand(g, player)
	}

	return g.balance
}

// Deal deals two cards from the top of the deck to all players in
// alternating order.
func deal(g *Game) {
	g.player = make([]deck.Card, 0, 5)
	g.dealer = make([]deck.Card, 0, 5)
	var card deck.Card
	for i := 0; i < 2; i++ {
		card, g.deck = draw(g.deck)
		g.player = append(g.player, card)
		card, g.deck = draw(g.deck)
		g.dealer = append(g.dealer, card)
	}
	g.phase = playerTurn
}

// Soft returns true if the score of the hand is a soft score. I.e.
// an ace is being counted as 11.
func Soft(hand ...deck.Card) bool {
	minScore := minScore(hand...)
	score := Score(hand...)
	return minScore != score
}

// Score returns the point value of the hand, converting aces between
// 11 and 1 as appropriate.
func Score(hand ...deck.Card) int {
	minScore := minScore(hand...)
	if minScore > 11 {
		return minScore
	}
	for _, c := range hand {
		if c.Rank == deck.Ace {
			return minScore - 1 + 11 // ace is currently worth 1; change it to be worth 11
		}
	}
	return minScore
}

// MinScore returns the point value of the hand, with any Aces counted
// as 1.
func minScore(hand ...deck.Card) int {
	var score int
	for _, c := range hand {
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

// Move is an action taken by players or the dealer on their term, as
// determined by their AI Play method.
type Move func(*Game)

// MoveHit draws a new card and adds it to the current player's hand.
func MoveHit(g *Game) {
	hand := g.CurrentHand()
	var card deck.Card
	card, g.deck = draw(g.deck)
	*hand = append(*hand, card)
	if Score(*hand...) >= 21 {
		MoveStand(g)
	}
}

// CurrentHand returns the hand of the player whose turn it is.
func (g *Game) CurrentHand() *[]deck.Card {
	switch g.phase {
	case playerTurn:
		return &g.player
	case dealerTurn:
		return &g.dealer
	default:
		panic("Current phase doesn't have a player turn")
	}
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

// MoveStand starts the next phase of gameplay.
func MoveStand(g *Game) {
	switch g.phase {
	case playerTurn:
		g.phase = dealerTurn
	case dealerTurn:
		g.phase = handOver
	}
}

// endHand compares and prints the player scores along with the
// outcome of the game, then clears the player hands.
func endHand(g *Game, ai AI) {
	pScore, dScore := Score(g.player...), Score(g.dealer...)

	switch {
	case pScore > 21:
		fmt.Println("You busted")
		g.balance--
	case dScore > 21:
		fmt.Println("Dealer busted")
		g.balance++
	case pScore > dScore:
		fmt.Println("You win!")
		g.balance++
	case dScore > pScore:
		fmt.Println("You lose!")
		g.balance--
	case dScore == pScore:
		fmt.Println("Draw")
	}
	fmt.Println()

	ai.Outcome([][]deck.Card{g.player}, g.dealer)
	g.player = nil
	g.dealer = nil
}
