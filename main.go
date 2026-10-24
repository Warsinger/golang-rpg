package main

import (
	"fmt"
	"log"

	"golang-rpg/game"
	"golang-rpg/internal/platform"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	var g *game.GameInfo
	var err error
	if platform.Browser() {
		ebiten.SetRunnableOnUnfocused(true)
		configs, images := platform.Resources()
		g, err = game.NewWithResources(configs, images)
	} else {
		g, err = game.New()
	}
	if err != nil {
		platform.ReportError("Start RPG: " + err.Error())
		log.Fatalf("start RPG: %v", err)
	}

	fmt.Println("Arrow Keys to move")
	fmt.Println("A key to attack monster when in range")
	fmt.Println("U key to pick up loot when in range")
	fmt.Println("R key to reset the game")
	fmt.Println("F to toggle full screen")
	fmt.Println("Q key to quit the game")

	if !platform.Browser() {
		mX, mY := ebiten.Monitor().Size()
		xScale := float64(g.Board.GetWidth()) / float64(mX)
		yScale := float64(g.Board.GetHeight()) / float64(mY)
		scale := max(xScale, yScale) * 1.1

		ebiten.SetWindowSize(int(float64(g.Board.GetWidth())/scale), int(float64(g.Board.GetHeight())/scale))
		ebiten.SetWindowTitle("Basic RPG")
		ebiten.SetWindowDecorated(false)

	}
	ebiten.SetTPS(15)

	if err := ebiten.RunGame(g); err != nil {
		platform.ReportError(err.Error())
		log.Fatal(err)
	}
}
