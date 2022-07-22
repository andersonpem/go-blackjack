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
	stack := append(h.Cards, c)

	h.Cards = stack
	h.Count = 0

	// Let's leave the aces to process lastly.
	var aces []card.Card

	// Re-calculate the stack value
	// Firstly, handle all non-special cards (second condition)
	for _, c := range h.Cards {
		if c.Name == "Ace" {
			// Keep track of which indices are aces in the hand
			aces = append(aces, c)
		} else {
			// All cards but the ace have only a single defined value
			h.Count += c.Value[0]
		}
	}

	/*
		Now, let's handle the aces
		Bear in mind that we can have up to 4 aces in a 4 deck game, so let's take no assumptions of value
	*/
	if len(aces) > 0 {
		var amount = len(aces)
		/*
			Let's fetch the first card just to have the reference values.
			Let's assume we don't know the variable values
		*/
		var sample = aces[0]

		if h.Count+(amount*sample.Value[1]) <= 21 {
			h.Count += amount * sample.Value[1]
		} else {
			// No assumptions. I know the value is 1, however let's make it dynamic.
			h.Count += amount * sample.Value[0]
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
	} else {
		fmt.Println("Somehow your stack is empty. This is not supposed to happen. Call the admins!")
		os.Exit(1)
	}
}
