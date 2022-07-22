/*
	The entire game Deck
*/
package deck

import (
	"blackjack/card"
	"blackjack/cl"
	"blackjack/gy"
	"blackjack/rng"
	"fmt"
)

type Deck struct {
	Cards []card.Card
}

func (d *Deck) Shuffle(message string) {
	if message != "" {
		fmt.Printf(cl.Yellow + message + "\n" + cl.Reset)
	} else {
		fmt.Printf(cl.Yellow + "Shuffling the cards...\n" + cl.Reset)
	}

	for i := range d.Cards {
		newPosition := rng.Intn(len(d.Cards) - 1)
		d.Cards[i], d.Cards[newPosition] = d.Cards[newPosition], d.Cards[i]
	}
}

// PickRand Picks a random card out of the deck
func (d *Deck) PickRand() card.Card {
	pos := rng.Intn(len(d.Cards) - 1)
	c := d.Cards[pos]
	d.Cards = card.RemoveCardByIndex(d.Cards, pos)
	return c
}

// PickNext Picks the next card in the deck
func (d *Deck) PickNext() card.Card {
	c := d.Cards[0]
	d.Cards = card.RemoveCardByIndex(d.Cards, 0)

	return c
}

// MonsterReborn TODO: we stopped here last time
// MonsterReborn plays yu-gi-oh and revives all cards from the graveyard
func (d *Deck) MonsterReborn(g *gy.Graveyard) {
	d.Cards = append(d.Cards, g.DefunctCards...)
	g.DefunctCards = nil
}
