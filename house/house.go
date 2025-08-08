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

func (h *House) HitUntil17(d *deck.Deck) {
	for h.Hand.Count < 17 {
		fmt.Println(cl.Blue + "House hits..." + cl.Reset)
		c := d.PickNext()
		h.Hand.Add(c)
		h.Stats(false)
	}
}

func (h *House) IsUpcardAce() bool {
	if len(h.Hand.Cards) == 2 {
		return h.Hand.Cards[1].Name == "Ace"
	}
	return false
}

func (h *House) Stats(hideFirstCard bool) {
	var names []string
	if h.Hand.Cards != nil {
		fmt.Println("House's hand:")
		if hideFirstCard {
			names = append(names, "Hidden Card")
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

func New(d *deck.Deck) House {
	h := House{}
	h.Hand.Add(d.PickNext()) // Hole card
	h.Hand.Add(d.PickNext()) // Up card
	return h
}
