package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := &GameInfo{}

	g.Init()

	fmt.Println("Arrow Keys to move")
	fmt.Println("A key to attack monster when in range")
	fmt.Println("U key to pick up loot when in range")
	fmt.Println("R key to reset the game")
	fmt.Println("Q key to quit the game")

	ebiten.SetWindowSize(g.Board.GetWidth(), g.Board.GetHeight())
	ebiten.SetWindowTitle("Basic RPG")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
