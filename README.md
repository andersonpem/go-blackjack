# Go Blackjack

A minimalistic blackjack game engine focused on increasing speed of calculation.

A simple exercise in go to implement the concept of pointers, references, and structs with functions in something palpable.

This game's purpose isn't to be a fully fledged Blackjack game engine with all the fluff, but, only to allow the player to calculate the sums of cards and predict when a blackjack is about to occur. This makes a player's ability to sum cards to increase in speed.

Most online games of Blackjack have way too much fluff, so I've decided to take what I needed to do into my own hands: calculate faster and faster and faster without the advertising in between each play.

## Scope

This game engine has all the basic rules of Blackjack accounted for:

- One or more decks of 52 cards (not including the jokers) with distinct suits

- All the non-numerical cards equal to 10 (Jack, Queen and King)

- The ace (1) value depends on the following rule: if the value in hand will surpass 21, its value is one. Otherwise, it's 11.

- If the sum of the cards in hand surpasses 21, the player loses.


## Some casino caveats

After a certain amount of plays, the deck is shuffled to avoid card counting. (Not that I imagine that someone would train card counting on this hahah)

If the amount of cards available to play is lower than 40%, all the cards that were already used are brought back to the stack, and the stack is re-shuffled immediately. This fixes the problem of a unfair game when the amount of cards is too low. Casinos do that too.

## How to play

Build the game with its modules and run it.

You'll be asked for an amount of decks to be used. The default value is 1. For a better experience, choose 2 decks.

The game will start and show your hand. You have to answer the sum of your hand as soon as possible. Every second counts.

After each input, you'll be shown how much time it took for you to answer. When you reach a blackjack, or you lose, you'll be shown your average response time in the current session. Then the game continues with the next cards in the stack.

Special inputs:

- bj (BlackJack): you asserted your hand is a blackjack. The game will check if that's true.

- l (lost): you asserted your hand is a loser. The value is bigger than 21.

- mr (debug function: MONSTER REBORN!): artificially triggers the removal of all cards from the graveyard and back into the game.


## The graveyard

After a hand is played, it's discarded to the graveyard. You'll be warned of how many cards are left in the stack, in the graveyard and in your hand.

## A roadmap

The engine's core is done. With some graphical upgrades, and some new modules, a multiplayer mode is possible, with one player being the house (a player that shows only one of the 2 initial cards) and other common players. The players other than yourself are to give a sense of real blackjack going.

Maybe we can add some currency to the game to give a sense of loss. We all know how humans are susceptible to the "loss aversion" phenomenom. Why not playing with that a bit? :)

Double down would be an interesting ability to explore too.

## Licensing

This code is for educational purposes only. MIT licensed. Do not take this as a production-ready code.
