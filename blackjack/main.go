package main

const (
	stand = "s"
	hit   = "h"
)

func main() {
	var gs GameState
	gs = Shuffle(gs)
	gs = Deal(gs)

	for {
		if gs.Phase != PlayerTurn {
			break
		}
		gs = TakePlayerTurn(gs)
	}

	for {
		if gs.Phase != DealerTurn {
			break
		}
		gs = TakeDealerTurn(gs)
	}

	gs = EndHand(gs)
}
