package main

import (
	"fmt"
	"log"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type ObjectInfo struct {
	GridX int
	GridY int
	Size  int
}

type Object interface {
	GetGridX() int
	GetGridY() int
	GetSize() int
}

type EntityInfo struct {
	ObjectInfo
	Name            string
	Defense         int
	MaxHealth       int
	CurrentHealth   int
	Gold            int
	ExperienceValue int
}
type Entity interface {
	Object
	GetName() string
	GetDefense() int
	GetMaxHealth() int
	GetCurrentHealth() int
	SetCurrentHealth(newHealth int)
	GetGold() int
	AddGold(newGold int) int
	GetExperienceValue() int
	Alive() bool
	Heal()
}
type AttackInfo struct {
	AttackPower int

	attackImg   *ebiten.Image
	attackFrame int
}
type Attacker interface {
	GetAttackPower() int
	Attack(d Entity)
}

func (o *ObjectInfo) GetGridX() int {
	return o.GridX
}
func (o *ObjectInfo) GetGridY() int {
	return o.GridY
}
func (o *ObjectInfo) GetSize() int {
	return o.Size
}

func (e *EntityInfo) Alive() bool {
	return e.CurrentHealth > 0
}
func (e *EntityInfo) Heal() {
	e.CurrentHealth = e.MaxHealth
}

func (e *EntityInfo) GetName() string {
	return e.Name
}

func (e *EntityInfo) GetDefense() int {
	return e.Defense
}

func (e *EntityInfo) GetMaxHealth() int {
	return e.MaxHealth
}

func (e *EntityInfo) GetCurrentHealth() int {
	return e.CurrentHealth
}
func (e *EntityInfo) SetCurrentHealth(newHealth int) {
	e.CurrentHealth = newHealth
}

func (e *EntityInfo) GetGold() int {
	return e.Gold
}

func (e *EntityInfo) AddGold(newGold int) int {
	e.Gold += newGold
	return e.Gold
}

func (e *EntityInfo) GetExperienceValue() int {
	return e.ExperienceValue
}

func (e *EntityInfo) DrawInfo(screen *ebiten.Image, x, y float32) {
	// Draw health inside the character
	var infoText string
	if e.CurrentHealth > 0 {
		infoText = fmt.Sprintf("%s\n%d/%d\n%dg", e.Name, e.CurrentHealth, e.MaxHealth, e.Gold)
	} else {
		infoText = fmt.Sprintf("Dead\n%s", e.Name)
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.LineSpacing = mplusNormalFace.Size * 1.5
	text.Draw(screen, infoText, mplusNormalFace, op)
}

func (a *AttackInfo) GetAttackPower() int {
	return a.AttackPower
}

func (a *AttackInfo) Attack(d Entity) {
	// calculate the attackers's attack value and subtract from defender's health
	pAttack := int(math.Max(float64(rand.IntN(a.AttackPower+1)-d.GetDefense()), 0))
	d.SetCurrentHealth(int(math.Max(float64(d.GetCurrentHealth()-pAttack), 0)))

	incrementFrame(&a.attackFrame)
}

func (a *AttackInfo) LoadAttackImage(path string) error {
	// TODO make a cache of image files for the different monsters

	attackImg, _, err := ebitenutil.NewImageFromFile(path + "attack1.png")
	if err != nil {
		log.Fatalf("failed to load attack sprite sheet: %v", err)
		return err
	}

	a.attackImg = attackImg
	return nil
}
