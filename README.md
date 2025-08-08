# Go Blackjack

A Go-based Blackjack game engine where you play against the house.

This project started as a simple exercise to practice Go fundamentals and has evolved into a playable console-based Blackjack game. It implements standard casino rules, including multiple decks, insurance bets, and automatic reshuffling.

## Scope

This game engine has all the basic rules of Blackjack accounted for:

-   Play against a House dealer who must hit until 17.
-   One or more decks of 52 cards with distinct suits.
-   All non-numerical cards (Jack, Queen, King) are valued at 10.
-   The Ace value is handled dynamically (1 or 11) to the player's best advantage.
-   If the sum of cards surpasses 21, it's a "bust" (loss).
-   **Insurance**: When the dealer's up-card is an Ace, you can place an insurance bet.

## Some casino caveats

To simulate a real casino environment and prevent card counting:

-   After a set number of hands, the deck is automatically re-shuffled.
-   If the number of cards remaining in the shoe drops below 40%, all used cards from the "graveyard" are returned to the shoe, which is then re-shuffled.

## How to play

Build the game with its modules and run it.

You'll be asked how many decks you want to use (default is 1; 2 is recommended for a more realistic game).

The game will deal your first hand and the house's hand (with one card hidden).

You can then choose your action:

-   **(h)it**: Take another card.
-   **(s)tand**: End your turn and let the house play.

The game will guide you through the round, determine the winner, and deal the next hand.

## The graveyard

After a hand is played, all cards from the player and the house are discarded to the graveyard. The game will inform you of how many cards are left in the shoe and in the graveyard.

## A roadmap

The engine's core is done. The next steps to make it a more complete experience are:

-   **Betting**: Introduce a wallet/currency system for placing bets.
-   **Advanced Plays**: Implement "Double Down" and "Split" actions.
-   **Multiplayer**: Abstract the player into a struct to allow for multiple human players at the same table.

## Licensing

This code is for educational purposes only. MIT licensed. Do not take this as production-ready code.