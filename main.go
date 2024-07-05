package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game := &Game{
		player: &Player{Entity: Entity{name: "Warsinger", Object: Object{size: 16}, defense: 1, maxHealth: 100}, level: 1, Attacker: Attacker{attack: 6}},
		monsters: []*Monster{
			{Entity: Entity{name: "Gorgon", Object: Object{size: 32}, defense: 2, maxHealth: 75}, Attacker: Attacker{attack: 4}},
			{Entity: Entity{name: "Barbol", Object: Object{size: 16}, defense: 3, maxHealth: 40}, Attacker: Attacker{attack: 2}},
			{Entity: Entity{name: "Barbol1", Object: Object{size: 16}, defense: 3, maxHealth: 40}, Attacker: Attacker{attack: 2}},
			{Entity: Entity{name: "Barbol2", Object: Object{size: 16}, defense: 3, maxHealth: 40}, Attacker: Attacker{attack: 2}},
			{Entity: Entity{name: "Barbol3", Object: Object{size: 16}, defense: 3, maxHealth: 40}, Attacker: Attacker{attack: 2}},
		},
		items: []Usable{
			&Treasure{Item: Item{value: 100, Object: Object{15, 344, 15}}},
			&Health{Item: Item{value: 50, Object: Object{127, 65, 15}}},
			&Health{Item: Item{value: 50, Object: Object{324, 44, 15}}},
		},
	}
	game.Init()

	fmt.Println("Arrow Keys to move")
	fmt.Println("A key to attack monster when in range")
	fmt.Println("U key to pick up loot when in range")
	fmt.Println("R key to reset the game")
	fmt.Println("Q key to quit the game")

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Basic RPG")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
