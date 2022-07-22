package main

import "blackjack/engine"

func main() {

	/*
		Some Todos:
		* For a realistic blackjack experience, keep the game going until it's reasonable another 4 players have lost
		* In this case, every new round has to reset the player stack
		* After N rounds (or x blackjacks), re-shuffle the deck no matter how many cards are there
		* Create a player type when on multiplayer mode
	*/

	// Starts the game engine and the game loop
	engine.Bootstrap()
}
