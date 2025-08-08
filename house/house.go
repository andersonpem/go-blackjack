package house

import (
	"blackjack/cl"
	"blackjack/deck"
	"blackjack/hand"
	"fmt"
)

type House struct {
	Hand hand.Hand
}

// HitUntil17 keeps drawing cards until the hand value is 17 or more.
func (h *House) HitUntil17(d *deck.Deck) {
	for h.Hand.Count < 17 {
		fmt.Println(cl.Blue + "House hits..." + cl.Reset)
		c := d.PickNext()
		h.Hand.Add(c)
		h.Stats(false)
	}
}

// IsUpcardAce checks if the house's visible card is an Ace.
// Assumes the second card (index 1) is the up-card.
func (h *House) IsUpcardAce() bool {
	if len(h.Hand.Cards) == 2 {
		return h.Hand.Cards[1].Name == "Ace"
	}
	return false
}

// Stats prints the house's hand. The first card is hidden if hideFirstCard is true.
func (h *House) Stats(hideFirstCard bool) {
	var names []string
	if h.Hand.Cards != nil {
		fmt.Println("House's hand:")
		if hideFirstCard {
			names = append(names, "Hidden Card")
			// Append the visible card (index 1)
			if len(h.Hand.Cards) > 1 {
				c := h.Hand.Cards[1]
				names = append(names, c.Name+"["+c.Suit+"]")
			}
		} else {
			for _, c := range h.Hand.Cards {
				names = append(names, c.Name+"["+c.Suit+"]")
			}
		}
		fmt.Println(cl.Yellow, names, cl.Reset)
		if !hideFirstCard {
			fmt.Printf("House's count: %d\n", h.Hand.Count)
		}
	}
}

// New creates a new house player and deals two cards from the deck.
func New(d *deck.Deck) House {
	h := House{}
	h.Hand.Add(d.PickNext()) // Hole card
	h.Hand.Add(d.PickNext()) // Up card
	return h
}
