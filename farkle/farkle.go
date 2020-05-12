package farkle

import (
	"math/rand"
	"time"
	"fmt"
)

func Seed() {
	rand.Seed(time.Now().UnixNano())
}

func RollDie() int {
	return rand.Intn(6) + 1
}

func RollDice(n int) []int {
	dice := make([]int, n)
	
	for i := range dice {
		dice[i] = RollDie()
	}

	return dice
}

type TurnState struct {
	Farkle    bool
	CurrScore int
	Dice      []int     
}

func newTurn() *TurnState {
	return &TurnState{false, 0, make([]int, 6)}
}

type GameState struct {
	Players []FarkleAgent
	Scores  []int
	Goal      int
	Turn     *TurnState
}

type FarkleAgent interface {
	Play(input *GameState) (diceToKeep int, score int, rollAgain bool) 
	RollOn(input *GameState) bool
	Name() string
}

func PlayGame(players []FarkleAgent, goal int) []int {
	//initialize game
	var theGame = GameState{}
	theGame.Players, theGame.Scores = players, make([]int, len(players))
	theGame.Goal = goal

	//Turn loop
	var currIdx = len(theGame.Players) - 1
	theGame.Turn = &TurnState{Farkle: true}
	for theGame.Scores[currIdx] < theGame.Goal {
		//Advance to next player
		currIdx = (currIdx + 1) % len(theGame.Players)
		//Renew dice unless rolling on
		if theGame.Turn.Farkle || ! theGame.Players[currIdx].RollOn(&theGame) {
			theGame.Turn = newTurn()
		}
		
		//Execute turn
		theGame.Scores[currIdx] += takeTurn(&theGame, theGame.Players[currIdx])
	}
	fmt.Println("Final Round!")
	fmt.Println("Scores are: ", theGame.Scores)
	//Some player has passed the goal, everyone else gets a final turn
	for i := 0; i < len(theGame.Players) - 1; i++ {
		//Advance to next player
		currIdx = (currIdx + 1) % len(theGame.Players)
		//Renew dice unless rolling on
		if theGame.Turn.Farkle || ! theGame.Players[currIdx].RollOn(&theGame) {
			theGame.Turn = newTurn()
		}
		
		//Execute turn
		theGame.Scores[currIdx] += takeTurn(&theGame, theGame.Players[currIdx])
	}
	return theGame.Scores
}

func takeTurn(theGame *GameState, currPlayer FarkleAgent) int {
	fmt.Println("It is now ", currPlayer.Name(), "'s turn!")
	//Roll dice until farkle or player quits
	for {
		//Roll dice
		theGame.Turn.Dice = RollDice(len(theGame.Turn.Dice))
		fmt.Println(currPlayer.Name(), " rolls: \n", theGame.Turn.Dice)
		//Get player decision
		var diceToKeep, rollScore, rollAgain = currPlayer.Play(theGame)

		if rollScore == 0 { //Farkle!
			theGame.Turn.Farkle = true
			theGame.Turn.CurrScore = 0
			break
		} else {
			theGame.Turn.CurrScore += rollScore
			fmt.Println(currPlayer.Name(), " is on: ", theGame.Turn.CurrScore)
		}
		//Take some dice out
		var diceToRoll = len(theGame.Turn.Dice) - diceToKeep
		if diceToRoll == 0 {diceToRoll = 6}
		theGame.Turn.Dice = make([]int, diceToRoll)
		//Player quit
		if ! rollAgain {
			break
		}
	}
	return theGame.Turn.CurrScore
}
