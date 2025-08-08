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
	"time"
)

/*
	This is the core "game engine"
*/

var plays = 0
var gamePlayer player.Player
var housePlayer = house.House{}
var lastBet int
var consecutiveWins int
var consecutiveLosses int // Track consecutive player losses
var safetyNetActive bool  // Flag for the "Lucky Break" bonus

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
	gamePlayer.Hand = hand.Hand{}
	housePlayer = house.New(&gameDeck)

	gamePlayer.Hand.Add(gameDeck.PickNext())
	housePlayer.Hand.Add(gameDeck.PickNext())
	gamePlayer.Hand.Add(gameDeck.PickNext())
	housePlayer.Hand.Add(gameDeck.PickNext())

	cls()
	cl.Pfln(
		"A new hand has been dealt!",
		cl.Blue,
	)
	fmt.Println("--------------------------------------------------")
}

func endOfPlay(bet int, outcome string, playerBlackjack bool) {
	// Apply safety net if active
	if safetyNetActive && outcome == "loss" {
		outcome = "push"
		cl.Pfln("Safety Net activated! Your bet is safe.", cl.Green)
	}
	safetyNetActive = false // Always reset after use

	if playerBlackjack {
		payout := (bet * 3) / 2
		gamePlayer.Balance += payout
		consecutiveWins++
		consecutiveLosses = 0
		cl.Pfln(
			fmt.Sprintf("BLACKJACK! You win $%d (3:2 payout).", payout),
			cl.Green,
		)
	} else {
		switch outcome {
		case "win":
			gamePlayer.Balance += bet
			consecutiveWins++
			consecutiveLosses = 0
			cl.Pfln(fmt.Sprintf("You win $%d!", bet), cl.Green)
		case "loss":
			gamePlayer.Balance -= bet
			consecutiveWins = 0
			consecutiveLosses++
			cl.Pfln(fmt.Sprintf("You lose $%d.", bet), cl.Red)
		case "push":
			consecutiveWins = 0
			consecutiveLosses = 0
			cl.Pfln("Push! Your bet is returned.", cl.Yellow)
		}
	}

	fmt.Printf("Your new balance is $%d\n", gamePlayer.Balance)

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

func handleLuckyBreak(reader *bufio.Reader) {
	cl.Pfln(
		fmt.Sprintf(
			"Lucky you! Time for a Lucky Break!",
		),
		cl.Yellow,
	)
	fmt.Println("Choose your bonus for this hand:")
	fmt.Println("1: Peek at the House's hole card")
	fmt.Println("2: Activate Safety Net (loss becomes a push)")

	for {
		fmt.Print("Your choice (1 or 2): ")
		choiceInput, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(choiceInput)

		if choice == "1" {
			cl.Pfln("Lucky Break: Peeking at the hole card!", cl.Blue)
			housePlayer.Stats(false) // Reveal the full hand
			fmt.Println("Memorize it! The card will be hidden again in 5 seconds...")
			time.Sleep(5 * time.Second)
			break
		} else if choice == "2" {
			safetyNetActive = true
			cl.Pfln(
				"Lucky Break: Safety Net is ACTIVE for this hand!",
				cl.Blue,
			)
			time.Sleep(2 * time.Second)
			break
		} else {
			cl.Pfln("Invalid choice. Please enter 1 or 2.", cl.Red)
		}
	}
	consecutiveLosses = 0 // Reset counter after offering the bonus
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

		// Handle Lucky Break before showing hands
		if consecutiveLosses >= 4 {
			handleLuckyBreak(reader)
			// Redisplay screen after bonus
			cls()
			cl.Pfln("Let's play your Lucky Break hand!", cl.Blue)
			fmt.Println("--------------------------------------------------")
		}

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
	if consecutiveWins > 4 {
		cl.Pfln(
			fmt.Sprintf(
				"You're on a %d-win streak! The house is getting nervous. Shuffling 5 times...",
				consecutiveWins,
			),
			cl.Yellow,
		)
		for i := 0; i < 5; i++ {
			gameDeck.Shuffle("")
		}
		consecutiveWins = 0
		return
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
