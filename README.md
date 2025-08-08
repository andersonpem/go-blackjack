# Go Blackjack

A Go-based Blackjack game engine where you play against the house.

This project started as a simple exercise to practice Go fundamentals and has evolved into a playable console-based Blackjack game. It implements standard casino rules, including multiple decks, betting, insurance, and automatic reshuffling.

## Scope

This game engine has all the basic rules of Blackjack accounted for:

-   Play against a House dealer who must hit until 17.
-   **Betting**: Start with a balance, place bets on each hand, and try not to go broke!
-   **Payouts**: Normal wins pay 1:1, and a natural Blackjack pays 3:2.
-   One or more decks of 52 cards with distinct suits.
-   The Ace value is handled dynamically (1 or 11) to the player's best advantage.
-   If the sum of cards surpasses 21, it's a "bust" (loss).
-   **Insurance**: When the dealer's up-card is an Ace, you can place an insurance side-bet.

## Some casino caveats

To simulate a real casino environment and prevent card counting:

-   After a set number of hands, the deck is automatically re-shuffled.
-   If the number of cards remaining in the shoe drops below 40%, all used cards from the "graveyard" are returned to the shoe, which is then re-shuffled.

## How to play

Build the game with its modules and run it.

You'll be asked for a starting balance and the number of decks to use.

Before each hand, you must enter a bet. The game will then deal the cards.

You can then choose your action:

-   **(h)it**: Take another card.
-   **(s)tand**: End your turn and let the house play.

The game will guide you through the round, handle your winnings or losses, and continue to the next hand. If your balance reaches zero, the game is over.

## A roadmap

The engine's core is solid. The next steps to make it a more complete experience are:

-   **Advanced Plays**: Implement "Double Down" and "Split" actions.
-   **Multiplayer**: The new `Player` struct makes it easier to add multiple human players at the same table.
-   **UI**: Refine the console user interface for a cleaner presentation.

## Licensing

This code is for educational purposes only. MIT licensed. Do not take this as production-ready code.