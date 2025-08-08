package player

import (
	"blackjack/hand"
)

type Player struct {
	Hand    hand.Hand
	Balance int
}

// New creates a new player with a starting balance.
func New(startingBalance int) Player {
	return Player{
		Balance: startingBalance,
	}
}
