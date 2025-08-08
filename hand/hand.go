/*
A hand is a stack of cards of a single player. If we have multiple players, each one has its own hand.
If we have a house logic, the house has a hand with 1 card visible and 1 card hidden in the start.
*/
package hand

import (
	"blackjack/card"
	"blackjack/cl"
	"fmt"
	"os"
)

type Hand struct {
	Cards []card.Card
	Count int
}

/*
Add: adds a new card to the player's hand.
*/
func (h *Hand) Add(c card.Card) {
	h.Cards = append(h.Cards, c)
	h.recalculate()
}

// IsBlackjack checks if the hand is a natural 21 with two cards.
func (h *Hand) IsBlackjack() bool {
	return len(h.Cards) == 2 && h.Count == 21
}

/*
recalculate: updates the hand's count, handling Aces correctly.
*/
func (h *Hand) recalculate() {
	h.Count = 0
	aceCount := 0

	// Sum non-ace cards first
	for _, c := range h.Cards {
		if c.Name == "Ace" {
			aceCount++
		} else {
			h.Count += c.Value[0]
		}
	}

	// Add aces, treating them as 11 unless it causes a bust
	for i := 0; i < aceCount; i++ {
		if h.Count+11 <= 21 {
			h.Count += 11
		} else {
			h.Count += 1
		}
	}
}

/*
Stats: Prints out a player's stats including the suit of each card.
*/
func (h Hand) Stats() {
	var names []string
	if h.Cards != nil {
		for _, c := range h.Cards {
			names = append(names, c.Name+"["+c.Suit+"]")
		}
		fmt.Printf("Your hand: [%d cards]\n", len(h.Cards))
		fmt.Println(cl.Yellow, names, cl.Reset)
		fmt.Printf("Your count: %d\n", h.Count)
	} else {
		fmt.Println(
			"Somehow your stack is empty. This is not supposed to happen. Call the admins!",
		)
		os.Exit(1)
	}
}
