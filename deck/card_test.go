package deck

import (
	"fmt"
	"math/rand"
	"testing"
)

func ExampleCard() {
	fmt.Println(Card{Rank: Ace, Suit: Hearts})
	fmt.Println(Card{Rank: Three, Suit: Diamonds})
	fmt.Println(Card{Rank: Nine, Suit: Spades})
	fmt.Println(Card{Rank: Queen, Suit: Clubs})
	fmt.Println(Card{Suit: Joker})

	// Output:
	// Ace of Hearts
	// Three of Diamonds
	// Nine of Spades
	// Queen of Clubs
	// Joker
}

func TestNew(t *testing.T) {
	cards := New()
	if len(cards) != cardsInDeck {
		t.Errorf("wrong number of cards in deck: got %d, want %d", len(cards), cardsInDeck)
	}
}

func TestDefaultSort(t *testing.T) {
	cards := New(DefaultSort)
	wantFirstCard := Card{Rank: Ace, Suit: Spades}
	wantLastCard := Card{Rank: King, Suit: Hearts}
	if cards[0] != wantFirstCard {
		t.Errorf("wanted first card to be %s, got %s", wantFirstCard, cards[0])
	}
	if cards[cardsInDeck-1] != wantLastCard {
		t.Errorf("wanted last card to be %s, got %s", wantLastCard, cards[cardsInDeck-1])
	}
}

func TestSort(t *testing.T) {
	opt := SortOption(Less)
	cards := New(opt)
	wantFirstCard := Card{Rank: Ace, Suit: Spades}
	wantLastCard := Card{Rank: King, Suit: Hearts}
	if cards[0] != wantFirstCard {
		t.Errorf("wanted first card to be %s, got %s", wantFirstCard, cards[0])
	}
	if cards[cardsInDeck-1] != wantLastCard {
		t.Errorf("wanted last card to be %s, got %s", wantLastCard, cards[cardsInDeck-1])
	}
}

func TestShuffle(t *testing.T) {
	// make shuffleRand deterministically generate the results of a Fisher-Yates shuffle with
	// seed 0:
	//   Two of Spades
	//   Two of Hearts
	//   Five of Spades
	//	 ...
	oldRand := shuffleRand
	shuffleRand = rand.New(rand.NewSource(0))

	want := []Card{
		{Two, Spades},
		{Two, Hearts},
		{Five, Spades},
	}

	shuffled := New(Shuffle)
	for i := 0; i < 3; i++ {
		if shuffled[i] != want[i] {
			t.Errorf("card %d is %s, want %s", i+1, shuffled[i], want[i])
		}
	}
	shuffleRand = oldRand
}

func TestJokers(t *testing.T) {
	wantJokers := 3
	cards := New(Jokers(wantJokers))
	var gotJokers int
	for _, card := range cards {
		if card.Suit == Joker {
			gotJokers++
		}
	}
	if gotJokers != wantJokers {
		t.Errorf("found %d Jokers, want %d", gotJokers, wantJokers)
	}
}

func TestFilter(t *testing.T) {
	filter := func(card Card) bool {
		return card.Rank == 2 || card.Rank == 3
	}
	cards := New(Filter(filter))
	for _, card := range cards {
		if card.Rank == 2 || card.Rank == 3 {
			t.Fatalf("found %s, do not want cards with Rank %d or %d", card, 2, 3)
		}
	}
}

func TestDeck(t *testing.T) {
	nDecks := 3
	cards := New(Deck(nDecks))
	wantLen := cardsInDeck * nDecks
	if len(cards) != wantLen {
		t.Errorf("deck has len %d, want %d", len(cards), wantLen)
	}
}
