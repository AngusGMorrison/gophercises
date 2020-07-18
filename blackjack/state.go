package main

import (
	"fmt"

	"github.com/angusgmorrison/gophercises/deck"
)

// GameState holds the state of a game at a point in time.
type GameState struct {
	Deck   []deck.Card
	Phase  Phase
	Player Hand
	Dealer Hand
}

// Phase represents the current stage of gameplay.
type Phase uint8

const (
	PlayerTurn Phase = iota
	DealerTurn
	HandOver
)

// CurrentPlayer returns the hand of the player whose turn it is.
func (gs *GameState) CurrentPlayer() *Hand {
	switch gs.Phase {
	case PlayerTurn:
		return &gs.Player
	case DealerTurn:
		return &gs.Dealer
	default:
		panic("Current phase doesn't have a player turn")
	}
}

// Shuffle returns a new shuffled deck.
func Shuffle(gs GameState) GameState {
	ret := clone(gs)
	ret.Deck = deck.New(deck.Deck(3), deck.Shuffle)
	return ret
}

// Deal deals two cards from the top of the deck to all players in
// alternating order.
func Deal(gs GameState) GameState {
	ret := clone(gs)
	ret.Player = make(Hand, 0, 5)
	ret.Dealer = make(Hand, 0, 5)
	var card deck.Card
	for i := 0; i < 2; i++ {
		card, ret.Deck = draw(ret.Deck)
		ret.Player = append(ret.Player, card)
		card, ret.Deck = draw(ret.Deck)
		ret.Dealer = append(ret.Dealer, card)
	}
	ret.Phase = PlayerTurn
	return ret
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

// TakePlayerTurn prompts the user for an action, processes it, and
// returns the resulting GameState.
func TakePlayerTurn(gs GameState) GameState {
	var input string
	fmt.Println("Player:", gs.Player)
	fmt.Println("Dealer:", gs.Dealer.DealerString())
	fmt.Println("What will you do? (h)it, (s)tand")
	fmt.Scanf("%s\n", &input)

	switch input {
	case hit:
		return Hit(gs)
	case stand:
		return Stand(gs)
	default:
		fmt.Println("Command not recognised: enter (h)it or (s)tand")
		return gs
	}
}

// TakeDealerTurn invokes shouldHit to determine whether the dealer
// should hit or stand, then returns the resulting GameState.
func TakeDealerTurn(gs GameState) GameState {
	if shouldHit(gs.Dealer.Score(), gs.Dealer.MinScore()) {
		return Hit(gs)
	}
	return Stand(gs)
}

// shouldHit determines whether the dealer AI should hit or stand.
func shouldHit(score, minScore int) bool {
	return score <= 16 || score == 17 && minScore != 17
}

// Hit draws a card from the top of the deck, adds it to the current
// player's hand and checks whether they're bust.
func Hit(gs GameState) GameState {
	ret := clone(gs)
	hand := ret.CurrentPlayer()
	var card deck.Card
	card, ret.Deck = draw(ret.Deck)
	*hand = append(*hand, card)
	if hand.Score() >= 21 {
		return Stand(ret)
	}
	return ret
}

// Stand triggers the next phase of gameplay.
func Stand(gs GameState) GameState {
	ret := clone(gs)
	switch ret.Phase {
	case PlayerTurn:
		ret.Phase = DealerTurn
	case DealerTurn:
		ret.Phase = HandOver
	}
	return ret
}

// EndHand compares and prints the player scores along with the
// outcome of the game, then clears the player hands.
func EndHand(gs GameState) GameState {
	ret := clone(gs)
	pScore, dScore := ret.Player.Score(), ret.Dealer.Score()

	fmt.Println("==FINAL HANDS==")
	fmt.Printf("Player: %s\nScore: %d\n", ret.Player, pScore)
	fmt.Printf("Player: %s\nScore: %d\n", ret.Dealer, dScore)

	switch {
	case pScore > 21:
		fmt.Println("You busted")
	case dScore > 21:
		fmt.Println("Dealer busted")
	case pScore > dScore:
		fmt.Println("You win!")
	case dScore > pScore:
		fmt.Println("You lose!")
	case dScore == pScore:
		fmt.Println("Draw")
	}
	fmt.Println()

	ret.Player = nil
	ret.Dealer = nil
	return ret
}

func clone(gs GameState) GameState {
	newState := GameState{
		Deck:   make([]deck.Card, len(gs.Deck)),
		Phase:  gs.Phase,
		Player: make(Hand, len(gs.Player)),
		Dealer: make(Hand, len(gs.Dealer)),
	}
	copy(newState.Deck, gs.Deck)
	copy(newState.Player, gs.Player)
	copy(newState.Dealer, gs.Dealer)
	return newState
}
