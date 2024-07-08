package main

import (
	"fmt"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := &GameInfo{}

	g.Init()

	fmt.Println("Arrow Keys to move")
	fmt.Println("A key to attack monster when in range")
	fmt.Println("U key to pick up loot when in range")
	fmt.Println("R key to reset the game")
	fmt.Println("F to toggle full screen")
	fmt.Println("Q key to quit the game")

	mX, mY := ebiten.Monitor().Size()
	xScale := float64(g.Board.GetWidth()) / float64(mX)
	yScale := float64(g.Board.GetHeight()) / float64(mY)
	scale := math.Max(xScale, yScale) * 1.1

	ebiten.SetWindowSize(int(float64(g.Board.GetWidth())/scale), int(float64(g.Board.GetHeight())/scale))
	ebiten.SetWindowTitle("Basic RPG")
	ebiten.SetWindowDecorated(false)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
