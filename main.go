package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game := &Game{
		Player: &Player{Entity: Entity{Name: "Warsinger", Object: Object{Size: 16}, Defense: 1, MaxHealth: 100}, Level: 1, Attacker: Attacker{AttackPower: 6}},
		// Monsters: []*Monster{
		// 	{Entity: Entity{Name: "Gorgon", Object: Object{Size: 32}, Defense: 2, MaxHealth: 75}, Attacker: Attacker{AttackPower: 4}},
		// 	{Entity: Entity{Name: "Barbol", Object: Object{Size: 16}, Defense: 3, MaxHealth: 40}, Attacker: Attacker{AttackPower: 2}},
		// 	{Entity: Entity{Name: "Barbol1", Object: Object{Size: 16}, Defense: 3, MaxHealth: 40}, Attacker: Attacker{AttackPower: 2}},
		// 	{Entity: Entity{Name: "Barbol2", Object: Object{Size: 16}, Defense: 3, MaxHealth: 40}, Attacker: Attacker{AttackPower: 2}},
		// 	{Entity: Entity{Name: "Barbol3", Object: Object{Size: 16}, Defense: 3, MaxHealth: 40}, Attacker: Attacker{AttackPower: 2}},
		// },
		// Items: []Usable{
		// 	&Treasure{Item: Item{Value: 100, Object: Object{15, 344, 15}}},
		// 	&Health{Item: Item{Value: 50, Object: Object{127, 65, 15}}},
		// 	&Health{Item: Item{Value: 50, Object: Object{324, 44, 15}}},
		// },
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
