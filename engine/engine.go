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
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

/*
	This is the core "game engine"
*/

var firstPLay = true
var plays = 0
var responseTimes []float64
var playerHand = hand.Hand{
	Cards: nil,
	Count: 0,
}

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
	// Initializes the game engine
	newGame()

	// Starts the game loop
	gameLoop()
}

func guessPlay(d *deck.Deck, h *hand.Hand) any {
	if firstPLay == true {
		firstPLay = false
	} else {
		h.Add(d.PickNext())
	}
	var userInput string
	var input int
	// Record the current timestamp
	start := time.Now()
	h.Stats()

	fmt.Println("What's the value? (time's ticking...): ")
	fmt.Scanln(&userInput)

	for userInput == "" {
		fmt.Println("You can't give an empty answer! What's the value? (time's still ticking...): ")
		fmt.Scanln(&userInput)
	}
	// When user has pressed enter
	end := time.Now()
	rTime := end.Sub(start).Seconds()
	fmt.Printf("You took "+cl.Blue+"%.2f seconds"+cl.Reset+" to answer.\n", rTime)
	// Append to the average times array
	responseTimes = append(responseTimes, rTime)
	if number, err := strconv.ParseInt(userInput, 10, 32); err == nil {
		input = int(number)

		if h.Count == 21 && input != 21 {
			fmt.Println(cl.Yellow + "YOU MISSED A BLACKJACK =O" + cl.Reset)
			return nil
		}
		// User number matches
		if h.Count == input {
			fmt.Println(cl.Green + "That's right! Keep it going!" + cl.Reset)
		} else {
			// User missed the number
			fmt.Println(cl.Red + "You missed it :( keep it going!" + cl.Reset)
			fmt.Printf("Actual value was "+cl.Yellow+"%d \n"+cl.Reset, h.Count)
		}
	} else {
		// User placed a string instead of a number
		switch userInput {
		case "bj":
			// User asserted an incoming BlackJack. Let's check it.
			if h.Count == 21 {
				fmt.Println(cl.Blue + "Congrats! You guessed a blackjack." + cl.Reset)
				return true
			} else {
				fmt.Println(cl.Red + "Oopsie! That's not a blackjack!" + cl.Reset)
				fmt.Printf("Actual value was %d \n", h.Count)
				return true
			}
		case "b":
			// User asserted that this play is lost for him.
			if h.Count > 21 {
				fmt.Println(cl.Blue + "BUST! You guessed it right." + cl.Reset)
				fmt.Printf("Last sum was: %d \n", h.Count)
				return true
			} else {
				fmt.Println(cl.Red + "You missed. Game is still on." + cl.Reset)
			}
		case "mr":
			// Summons the Monster Reborn routine.
			// Artificially brings all the cards from the graveyard back to the deck. Debug
			fmt.Println(cl.Yellow + "MONSTER REBORN!" + cl.Reset)
			gameDeck.MonsterReborn(&Grave)
			gameDeck.Shuffle("Shuffling after a debug monster reborn request!")
		}
	}
	return nil
}

func avgResponseTime() {
	var total, average float64
	for _, responseTime := range responseTimes {
		total += responseTime
	}
	average = total / float64(len(responseTimes))
	fmt.Printf(cl.Blue+"Your average response time in this session was %.2f seconds.\n"+cl.Reset, average)
}

func newHand(d *deck.Deck) {
	playerHand = hand.Hand{
		Cards: nil,
		Count: 0,
	}
	// A new play always have two cards in your hand
	playerHand.Add(d.PickNext())
	playerHand.Add(d.PickNext())
}

func cls() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func newGame() {
	cls()
	var amountOfDecks = 1
	fmt.Printf("   _____         ____  _            _    _            _    \n  / ____|       |  _ \\| |          | |  (_)          | |   \n | |  __  ___   | |_) | | __ _  ___| | ___  __ _  ___| | __\n | | |_ |/ _ \\  |  _ <| |/ _` |/ __| |/ / |/ _` |/ __| |/ /\n | |__| | (_) | | |_) | | (_| | (__|   <| | (_| | (__|   < \n  \\_____|\\___/  |____/|_|\\__,_|\\___|_|\\_\\ |\\__,_|\\___|_|\\_\\\n                                       _/ |                \n                                      |__/       \n")
	fmt.Println("Let's see how fast you can calculate sums in a Blackjack game :)")
	fmt.Println("By @AndersonPEM")
	fmt.Println("Card suits: Ht (Hearts) Cb (Clubs) Dm (diamonds) Sp (spades)")
	fmt.Println("How many decks do you want in your game? [default is 1, recommended 2 for realism]: ")
	_, err := fmt.Scanln(&amountOfDecks)
	if err != nil {
		fmt.Println("Defaulting to a single deck")
	}
	// Let's generate the stack of decks for our game
	fmt.Printf(cl.Yellow+"Provisioning %d deck(s) of cards...\n"+cl.Reset, amountOfDecks)
	gameDeck = deck.Deck{Cards: card.ProvisionDecks(amountOfDecks)}
	gameDeckSize = len(gameDeck.Cards)
	// What good is a stack of cards without randomness?
	gameDeck.Shuffle("Executing first game's shuffle...")
	fmt.Println("Ready to start. Press enter to proceed.")
	fmt.Scanln()
}

func newPlay() {
	firstPLay = true
	newHand(&gameDeck)
	cls()
	// Todo: create a inttostr abstraction, this shit is ugly af
	cl.Pfln("A new play just started! Using the same deck. Cards available: "+strconv.Itoa(len(gameDeck.Cards)), cl.Blue)
	cl.Pfln("Graveyard count: "+strconv.Itoa(len(Grave.DefunctCards)), cl.Blue)
}

func gameLoop() {
	newPlay()
	// The game loop (first version)
	for {
		var eval = guessPlay(&gameDeck, &playerHand)
		if eval != nil {
			endOfPlay()
		} else {
			if playerHand.Count > 21 {
				fmt.Println(cl.Yellow+"Your hand exceeded 21: "+cl.Reset, playerHand.Count)
				endOfPlay()
			}
		}
		/*
			If the amount of cards is too low, we should add all the cards from the graveyard back and reshuffle.
			This avoids unfair advantages and possible card counting.
		*/
		if percentOfNumber(gameDeckSize, len(gameDeck.Cards)) <= reviveCardsWhenLessThan {
			fmt.Printf(cl.Yellow+"%d percent of the cards were already defuncted. Reviving all the cards... \n "+cl.Reset, percentOfNumber(gameDeckSize, len(gameDeck.Cards)))
			gameDeck.MonsterReborn(&Grave)
			gameDeck.Shuffle("Shuffling after refilling the game deck with used cards...")
		}
		// If N rounds have passed, reshuffle.
		if plays > 0 && plays%reshuffleEveryPlays == 0 {
			gameDeck.Shuffle(strconv.Itoa(reshuffleEveryPlays) + " plays have been passed. Re-shuffling (gotcha card counters!)...")
		}
	}
}

func endOfPlay() {
	avgResponseTime()
	Grave.Add(playerHand.Cards)
	fmt.Println("Press enter to start a new game :)")
	fmt.Scanln()
	plays++
	newPlay()
}

/*
Returns the percentage that partialNumber represents from totalNumber.
*/
func percentOfNumber(totalNumber int, partialNumber int) int {
	// X = (100 * vp ) / vt
	var percent = (100 * partialNumber) / totalNumber
	return percent
}
