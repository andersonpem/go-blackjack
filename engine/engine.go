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
	"blackjack/player"
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
var gamePlayer player.Player
var housePlayer = house.House{}
var lastBet int
var consecutiveWins int // Track consecutive player wins

var reshuffleEveryPlays = 5
var reviveCardsWhenLessThan = 40
var gameDeck = deck.Deck{}
var gameDeckSize = 52
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
	var startingBalance = 100
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("   _____         ____  _            _    _            _    \n  / ____|       |  _ \\| |          | |  (_)          | |   \n | |  __  ___   | |_) | | __ _  ___| | ___  __ _  ___| | __\n | | |_ |/ _ \\  |  _ <| |/ _` |/ __| |/ / |/ _` |/ __| |/ /\n | |__| | (_) | | |_) | | (_| | (__|   <| | (_| | (__|   < \n  \\_____|\\___/  |____/|_|\\__,_|\\___|_|\\_\\ |\\__,_|\\___|_|\\_\\\n                                       _/ |                \n                                      |__/       \n")
	fmt.Println("Welcome to Go Blackjack!")
	fmt.Println("By @AndersonPEM")
	fmt.Println(
		"How many decks? [default is 1, recommended 2 for realism]: ",
	)
	fmt.Scanln(&amountOfDecks)
	if amountOfDecks <= 0 {
		amountOfDecks = 1
	}

	fmt.Println("Enter your starting balance [default is 100]: ")
	balanceInput, _ := reader.ReadString('\n')
	if val, err := strconv.Atoi(strings.TrimSpace(balanceInput)); err == nil &&
		val > 0 {
		startingBalance = val
	}

	gamePlayer = player.New(startingBalance)

	fmt.Printf(
		cl.Yellow+"Provisioning %d deck(s) of cards...\n"+cl.Reset,
		amountOfDecks,
	)
	gameDeck = deck.Deck{Cards: card.ProvisionDecks(amountOfDecks)}
	gameDeckSize = len(gameDeck.Cards)
	gameDeck.Shuffle("Executing first game's shuffle...")
	fmt.Println("Ready to start. Press enter to place your first bet.")
	_, _ = fmt.Scanln()
}

func newPlay() {
	// Reset player and house hands
	gamePlayer.Hand = hand.Hand{}
	housePlayer = house.New(&gameDeck)

	// Deal cards
	gamePlayer.Hand.Add(gameDeck.PickNext())
	housePlayer.Hand.Add(gameDeck.PickNext()) // Hole card
	gamePlayer.Hand.Add(gameDeck.PickNext())
	housePlayer.Hand.Add(gameDeck.PickNext()) // Up card

	cls()
	cl.Pfln(
		"A new hand has been dealt!",
		cl.Blue,
	)
	fmt.Println("--------------------------------------------------")
}

func endOfPlay(bet int, outcome string, playerBlackjack bool) {
	// Payout and consecutive win logic
	if playerBlackjack {
		payout := (bet * 3) / 2
		gamePlayer.Balance += payout
		consecutiveWins++
		cl.Pfln(
			fmt.Sprintf("BLACKJACK! You win $%d (3:2 payout).", payout),
			cl.Green,
		)
	} else {
		switch outcome {
		case "win":
			gamePlayer.Balance += bet
			consecutiveWins++
			cl.Pfln(fmt.Sprintf("You win $%d!", bet), cl.Green)
		case "loss":
			gamePlayer.Balance -= bet
			consecutiveWins = 0 // Streak broken
			cl.Pfln(fmt.Sprintf("You lose $%d.", bet), cl.Red)
		case "push":
			consecutiveWins = 0 // Streak broken
			cl.Pfln("Push! Your bet is returned.", cl.Yellow)
		}
	}

	fmt.Printf("Your new balance is $%d\n", gamePlayer.Balance)

	// Collect all cards from the table
	Grave.Add(gamePlayer.Hand.Cards)
	Grave.Add(housePlayer.Hand.Cards)

	if gamePlayer.Balance > 0 {
		fmt.Println("Press enter to start the next hand :)")
		_, _ = fmt.Scanln()
		plays++
	}
}

func getBet(reader *bufio.Reader) int {
	prompt := "Place your bet: "
	if lastBet > 0 {
		prompt = fmt.Sprintf("Place your bet [Enter for $%d]: ", lastBet)
	}

	for {
		cls()
		fmt.Printf("Your balance is $%d\n", gamePlayer.Balance)
		fmt.Print(prompt)

		betInput, _ := reader.ReadString('\n')
		trimmedInput := strings.TrimSpace(betInput)

		var currentBet int

		if trimmedInput == "" {
			currentBet = lastBet
		} else {
			parsedBet, err := strconv.Atoi(trimmedInput)
			if err != nil {
				cl.Pfln("Invalid bet. Please enter a number.", cl.Red)
				continue
			}
			currentBet = parsedBet
		}

		if currentBet <= 0 {
			cl.Pfln("You must enter a bet greater than $0.", cl.Red)
			continue
		}
		if currentBet > gamePlayer.Balance {
			cl.Pfln("You cannot bet more than your balance.", cl.Red)
			continue
		}

		lastBet = currentBet
		return currentBet
	}
}

func gameLoop() {
	reader := bufio.NewReader(os.Stdin)

	for {
		if gamePlayer.Balance <= 0 {
			cl.Pfln("You've run out of money! Game over.", cl.Red)
			os.Exit(0)
		}

		currentBet := getBet(reader)
		newPlay()

		housePlayer.Stats(true)
		gamePlayer.Hand.Stats()

		if gamePlayer.Hand.IsBlackjack() {
			if housePlayer.Hand.IsBlackjack() {
				endOfPlay(currentBet, "push", false)
			} else {
				endOfPlay(currentBet, "win", true)
			}
			checkReshuffleConditions()
			continue
		}

		if housePlayer.IsUpcardAce() {
			handleInsurance(reader, currentBet)
			if housePlayer.Hand.IsBlackjack() {
				cl.Pfln("House has Blackjack!", cl.Red)
				housePlayer.Stats(false)
				endOfPlay(currentBet, "loss", false)
				checkReshuffleConditions()
				continue
			}
		}

		if playerTurn(reader) {
			cl.Pfln(
				fmt.Sprintf("BUST! You lose with %d.", gamePlayer.Hand.Count),
				cl.Red,
			)
			endOfPlay(currentBet, "loss", false)
			checkReshuffleConditions()
			continue
		}

		houseTurn()

		outcome := determineWinner()
		endOfPlay(currentBet, outcome, false)

		checkReshuffleConditions()
	}
}

func handleInsurance(reader *bufio.Reader, mainBet int) {
	cl.Pfln("House shows an Ace. Insurance? (y/n)", cl.Yellow)
	input, _ := reader.ReadString('\n')
	if strings.TrimSpace(input) == "y" {
		insuranceBet := mainBet / 2
		if insuranceBet == 0 {
			insuranceBet = 1
		}
		if gamePlayer.Balance-mainBet < insuranceBet {
			cl.Pfln("Not enough balance for insurance bet.", cl.Red)
			return
		}

		if housePlayer.Hand.IsBlackjack() {
			payout := insuranceBet * 2
			gamePlayer.Balance += payout
			cl.Pfln(
				fmt.Sprintf(
					"House has Blackjack! Insurance pays $%d.",
					payout,
				),
				cl.Green,
			)
		} else {
			gamePlayer.Balance -= insuranceBet
			cl.Pfln(
				fmt.Sprintf(
					"House does not have Blackjack. Insurance bet of $%d lost.",
					insuranceBet,
				),
				cl.Red,
			)
			fmt.Printf("Your new balance is $%d\n", gamePlayer.Balance)
		}
	}
}

func playerTurn(reader *bufio.Reader) (busted bool) {
	for {
		cl.Pfln("Your action? (h)it or (s)tand", cl.Yellow)
		input, _ := reader.ReadString('\n')
		action := strings.TrimSpace(input)

		if action == "h" {
			gamePlayer.Hand.Add(gameDeck.PickNext())
			gamePlayer.Hand.Stats()
			if gamePlayer.Hand.Count > 21 {
				return true
			}
		} else if action == "s" {
			return false
		} else {
			cl.Pfln("Invalid input. Please enter 'h' or 's'.", cl.Red)
		}
	}
}

func houseTurn() {
	fmt.Println("--------------------------------------------------")
	cl.Pfln("Player stands. Revealing house's hand...", cl.Blue)
	housePlayer.Stats(false)
	housePlayer.HitUntil17(&gameDeck)
}

func determineWinner() string {
	playerCount := gamePlayer.Hand.Count
	houseCount := housePlayer.Hand.Count

	cl.Pfln(
		fmt.Sprintf(
			"Final Score -> Player: %d | House: %d",
			playerCount,
			houseCount,
		),
		cl.Blue,
	)

	if houseCount > 21 {
		return "win"
	}
	if playerCount > houseCount {
		return "win"
	}
	if houseCount > playerCount {
		return "loss"
	}
	return "push"
}

func checkReshuffleConditions() {
	// New check for consecutive wins
	if consecutiveWins > 4 {
		cl.Pfln(
			fmt.Sprintf(
				"You're on a %d-win streak! The house is getting nervous. Shuffling 5 times...",
				consecutiveWins,
			),
			cl.Yellow,
		)
		for i := 0; i < 5; i++ {
			gameDeck.Shuffle("") // Pass empty string for default message
		}
		consecutiveWins = 0 // Reset the counter after the special shuffle
		return              // No need to check other shuffle conditions
	}

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

	if plays > 0 && plays%reshuffleEveryPlays == 0 {
		gameDeck.Shuffle(
			strconv.Itoa(reshuffleEveryPlays) +
				" plays have passed. Re-shuffling...",
		)
	}
}

func percentOfNumber(totalNumber int, partialNumber int) int {
	if totalNumber == 0 {
		return 0
	}
	return (100 * partialNumber) / totalNumber
}
