/*
The game "engine"
*/
package engine

import (
	"blackjack/card"
	"blackjack/cl"
	"blackjack/deck"
	"blackjack/gy"
	"blackjack/hand"
	"blackjack/house"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

/*
	This is the core "game engine"
*/

var plays = 0
var playerHand = hand.Hand{}
var housePlayer = house.House{}

// After N plays, reshuffle the deck. This will be a deterrent for card counters haha
var reshuffleEveryPlays = 5

// If cards are less than N percent, add all the cards from the graveyard back to the game deck and reshuffle
var reviveCardsWhenLessThan = 40

// Initializer for the stack of decks we'll use in this session
var gameDeck = deck.Deck{}

// The full size of the game's stack of cards for calculations. Defaults to a single full deck.
var gameDeckSize = 52

// Grave An empty graveyard for cards
var Grave = gy.Graveyard{DefunctCards: nil}

// Bootstrap Initializes everything and gets the game loop going
func Bootstrap() {
	newGame()
	gameLoop()
}

func cls() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

func newGame() {
	cls()
	var amountOfDecks = 1
	fmt.Printf("   _____         ____  _            _    _            _    \n  / ____|       |  _ \\| |          | |  (_)          | |   \n | |  __  ___   | |_) | | __ _  ___| | ___  __ _  ___| | __\n | | |_ |/ _ \\  |  _ <| |/ _` |/ __| |/ / |/ _` |/ __| |/ /\n | |__| | (_) | | |_) | | (_| | (__|   <| | (_| | (__|   < \n  \\_____|\\___/  |____/|_|\\__,_|\\___|_|\\_\\ |\\__,_|\\___|_|\\_\\\n                                       _/ |                \n                                      |__/       \n")
	fmt.Println("Welcome to Go Blackjack!")
	fmt.Println("By @AndersonPEM")
	fmt.Println("Card suits: Ht (Hearts) Cb (Clubs) Dm (diamonds) Sp (spades)")
	fmt.Println(
		"How many decks do you want in your game? [default is 1, recommended 2 for realism]: ",
	)
	_, err := fmt.Scanln(&amountOfDecks)
	if err != nil || amountOfDecks <= 0 {
		amountOfDecks = 1
		fmt.Println("Defaulting to a single deck.")
	}
	// Let's generate the stack of decks for our game
	fmt.Printf(
		cl.Yellow+"Provisioning %d deck(s) of cards...\n"+cl.Reset,
		amountOfDecks,
	)
	gameDeck = deck.Deck{Cards: card.ProvisionDecks(amountOfDecks)}
	gameDeckSize = len(gameDeck.Cards)
	// What good is a stack of cards without randomness?
	gameDeck.Shuffle("Executing first game's shuffle...")
	fmt.Println("Ready to start. Press enter to deal the first hand.")
	_, _ = fmt.Scanln()
}

func newPlay() {
	// Reset player and house hands
	playerHand = hand.Hand{}
	housePlayer = house.New(&gameDeck)

	// A new play always has two cards in your hand
	playerHand.Add(gameDeck.PickNext())
	playerHand.Add(gameDeck.PickNext())

	cls()
	cl.Pfln(
		"A new play just started! Cards available: "+strconv.Itoa(len(gameDeck.Cards)),
		cl.Blue,
	)
	cl.Pfln(
		"Graveyard count: "+strconv.Itoa(len(Grave.DefunctCards)),
		cl.Blue,
	)
	fmt.Println("--------------------------------------------------")
}

func endOfPlay() {
	// Collect all cards from the table
	Grave.Add(playerHand.Cards)
	Grave.Add(housePlayer.Hand.Cards)

	fmt.Println("Press enter to start the next hand :)")
	_, _ = fmt.Scanln()
	plays++
	newPlay()
}

func gameLoop() {
	newPlay()
	reader := bufio.NewReader(os.Stdin)

	for {
		// Show initial hands
		housePlayer.Stats(true) // Hide house's first card
		playerHand.Stats()

		// Check for player blackjack
		if playerHand.IsBlackjack() {
			cl.Pfln("BLACKJACK! You win!", cl.Green)
			endOfPlay()
			continue
		}

		// Insurance Rule: If house's up-card is an Ace
		if housePlayer.IsUpcardAce() {
			cl.Pfln("House shows an Ace. Insurance? (y/n)", cl.Yellow)
			input, _ := reader.ReadString('\n')
			if strings.TrimSpace(input) == "y" {
				if housePlayer.Hand.IsBlackjack() {
					cl.Pfln(
						"House has Blackjack! You lose the hand but win the insurance bet. (Push)",
						cl.Blue,
					)
					housePlayer.Stats(false) // Reveal hand
					endOfPlay()
					continue
				} else {
					cl.Pfln("House does not have Blackjack. Insurance lost.", cl.Red)
				}
			}
		}

		// Player's turn
		playerBusted := playerTurn(reader)
		if playerBusted {
			cl.Pfln(
				fmt.Sprintf("BUST! You lose with %d.", playerHand.Count),
				cl.Red,
			)
			endOfPlay()
			continue
		}

		// House's turn
		houseTurn()

		// Determine winner
		determineWinner()
		endOfPlay()

		// Reshuffle logic
		checkReshuffleConditions()
	}
}

func playerTurn(reader *bufio.Reader) (busted bool) {
	for {
		cl.Pfln("Your action? (h)it or (s)tand", cl.Yellow)
		input, _ := reader.ReadString('\n')
		action := strings.TrimSpace(input)

		if action == "h" {
			playerHand.Add(gameDeck.PickNext())
			playerHand.Stats()
			if playerHand.Count > 21 {
				return true // Busted
			}
		} else if action == "s" {
			return false // Stood
		} else {
			cl.Pfln("Invalid input. Please enter 'h' or 's'.", cl.Red)
		}
	}
}

func houseTurn() {
	fmt.Println("--------------------------------------------------")
	cl.Pfln("Player stands. Revealing house's hand...", cl.Blue)
	housePlayer.Stats(false) // Reveal house's hand
	housePlayer.HitUntil17(&gameDeck)
}

func determineWinner() {
	cl.Pfln(
		fmt.Sprintf(
			"Final Score -> Player: %d | House: %d",
			playerHand.Count,
			housePlayer.Hand.Count,
		),
		cl.Blue,
	)

	if housePlayer.Hand.Count > 21 {
		cl.Pfln("House busts! You win!", cl.Green)
	} else if playerHand.Count > housePlayer.Hand.Count {
		cl.Pfln("You win!", cl.Green)
	} else if housePlayer.Hand.Count > playerHand.Count {
		cl.Pfln("House wins!", cl.Red)
	} else {
		cl.Pfln("Push! It's a tie.", cl.Yellow)
	}
}

func checkReshuffleConditions() {
	/*
		If the amount of cards is too low, we should add all the cards from the graveyard back and reshuffle.
		This avoids unfair advantages and possible card counting.
	*/
	if percentOfNumber(gameDeckSize, len(gameDeck.Cards)) <= reviveCardsWhenLessThan {
		fmt.Printf(
			cl.Yellow+"%d percent of the cards were already used. Reviving all cards... \n "+cl.Reset,
			100-percentOfNumber(gameDeckSize, len(gameDeck.Cards)),
		)
		gameDeck.MonsterReborn(&Grave)
		gameDeck.Shuffle(
			"Shuffling after refilling the game deck with used cards...",
		)
	}
	// If N rounds have passed, reshuffle.
	if plays > 0 && plays%reshuffleEveryPlays == 0 {
		gameDeck.Shuffle(
			strconv.Itoa(reshuffleEveryPlays) +
				" plays have passed. Re-shuffling (gotcha card counters!)...",
		)
	}
}

/*
Returns the percentage that partialNumber represents from totalNumber.
*/
func percentOfNumber(totalNumber int, partialNumber int) int {
	if totalNumber == 0 {
		return 0
	}
	return (100 * partialNumber) / totalNumber
}
