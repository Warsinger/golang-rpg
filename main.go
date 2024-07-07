package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := &Game{
		Board:  &Board{Width: 800, Height: 600, GridSize: 20},
		Player: &Player{Entity: Entity{Name: "Warsinger", Object: Object{Size: 2}, Defense: 1, MaxHealth: 100}, Level: 1, Attacker: Attacker{AttackPower: 6}},
	}
	g.Init()

	fmt.Println("Arrow Keys to move")
	fmt.Println("A key to attack monster when in range")
	fmt.Println("U key to pick up loot when in range")
	fmt.Println("R key to reset the game")
	fmt.Println("Q key to quit the game")

	ebiten.SetWindowSize(g.Board.Width, g.Board.Height)
	ebiten.SetWindowTitle("Basic RPG")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
