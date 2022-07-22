// Package gy Card graveyard
package gy

import (
	"blackjack/card"
	"blackjack/cl"
	"strconv"
)

type Graveyard struct {
	DefunctCards []card.Card
}

// Add This card has died, unfortunately. It is out of the game, for now.
func (g *Graveyard) Add(c []card.Card) {
	// Research what's unpacking a slice
	g.DefunctCards = append(g.DefunctCards, c...)
	cl.Pfln("Added "+strconv.Itoa(len(c))+" just played cards to the graveyard.", cl.Blue)
}
