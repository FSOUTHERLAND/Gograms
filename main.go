package main

import (
	"fmt"
	"./farkle"
)

type cheater struct {
	points int
	name string
} 

func (cheat *cheater) RollOn(theGame *farkle.GameState) bool {
	return false
}

func (cheat *cheater) Play(theGame *farkle.GameState) (diceToKeep int, score int, rollAgain bool) {
	if theGame.Turn.CurrScore >= theGame.Goal {
		return 0, cheat.points, false
	}
	return 0, cheat.points, true
}

func (cheat *cheater) Name() string {
	return cheat.name
}

func main() {
	farkle.Seed()
	var player1 = cheater{100, "cheat1"}
	var player2 = cheater{50, "cheat2"}
	var players = []farkle.FarkleAgent{&player1, &player2}
	fmt.Println(farkle.PlayGame(players, 500))
}