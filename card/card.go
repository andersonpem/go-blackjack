package card

type Card struct {
	Deck  int
	Name  string
	Suit  string
	Value []int
}

// Get the full name of the current card
func (c Card) name() string {
	return c.Name + " of " + c.Suit
}

var suits = []string{"Ht", "Dm", "Cb", "Sp"}

var cards = []Card{
	{Deck: 1, Name: "Ace", Value: []int{1, 11}},
	{Deck: 1, Name: "2", Value: []int{2}},
	{Deck: 1, Name: "3", Value: []int{3}},
	{Deck: 1, Name: "4", Value: []int{4}},
	{Deck: 1, Name: "5", Value: []int{5}},
	{Deck: 1, Name: "6", Value: []int{6}},
	{Deck: 1, Name: "7", Value: []int{7}},
	{Deck: 1, Name: "8", Value: []int{8}},
	{Deck: 1, Name: "9", Value: []int{9}},
	{Deck: 1, Name: "10", Value: []int{10}},
	{Deck: 1, Name: "Jack", Value: []int{10}},
	{Deck: 1, Name: "Queen", Value: []int{10}},
	{Deck: 1, Name: "King", Value: []int{10}},
}

func getMeADeck(deckId int) []Card {
	if deckId == 0 {
		deckId = 1
	}
	var deck []Card
	for i := 0; i < len(suits); i++ {
		for x := 0; x < len(cards); x++ {
			currentCard := cards[x]
			currentCard.Deck = deckId
			currentCard.Suit = suits[i]
			deck = append(deck, currentCard)
		}
	}

	return deck
}

func ProvisionDecks(amount int) []Card {
	var decks []Card
	for i := 0; i < amount; i++ {
		decks = append(decks, getMeADeck(i+1)...)
	}
	return decks
}

func RemoveCardByIndex(c []Card, index int) []Card {
	ret := make([]Card, 0)
	ret = append(ret, c[:index]...)
	return append(ret, c[index+1:]...)
}
