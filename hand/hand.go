/*
A hand is a stack of cards of a single player.
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

func (h *Hand) Add(c card.Card) {
	h.Cards = append(h.Cards, c)
	h.recalculate()
}

func (h *Hand) IsBlackjack() bool {
	return len(h.Cards) == 2 && h.Count == 21
}

func (h *Hand) recalculate() {
	h.Count = 0
	aceCount := 0

	for _, c := range h.Cards {
		if c.Name == "Ace" {
			aceCount++
		} else {
			h.Count += c.Value[0]
		}
	}

	for i := 0; i < aceCount; i++ {
		if h.Count+11 <= 21 {
			h.Count += 11
		} else {
			h.Count += 1
		}
	}
}

func (h Hand) Stats() {
	var names []string
	if h.Cards != nil {
		for _, c := range h.Cards {
			names = append(names, c.Name+"["+c.Suit+"]")
		}
		fmt.Printf("Player's hand: [%d cards]\n", len(h.Cards))
		fmt.Println(cl.Yellow, names, cl.Reset)
		fmt.Printf("Player's count: %d\n", h.Count)
	} else {
		fmt.Println(
			"Somehow the player stack is empty. This is not supposed to happen.",
		)
		os.Exit(1)
	}
}
