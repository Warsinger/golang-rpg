package game

import (
	"fmt"
	"image"
	"image/color"
	"io/fs"
	"math"

	"golang-rpg/assets"
	"golang-rpg/internal/collections"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type AnimState string

type PlayerInfo struct {
	EntityInfo
	AttackInfo
	Experience int
	Level      int
	Character  string

	walkAsset    assets.Asset
	walkFrame    int
	walkTicks    int
	movePosition *image.Point

	deathAsset assets.Asset

	// facingLeft is runtime-only; the source sheets and a fresh scene face right.
	facingLeft bool
	animState  AnimState
	moveQueue  *collections.Queue[image.Point]
}

type Player interface {
	Entity
	Attacker

	Draw(screen *ebiten.Image, b Board) error

	GetExperience() int
	GetLevel() int
	AttackMonster(m Monster)
	AddExperience(xp int)
	LevelUp(newLevel int)
	Move(direction Direction, b Board) bool
	EnqueueMove(pt *image.Point)
	DequeueMove() *image.Point
	AdvanceMovement()
	Idle()
	LoadImages(am assets.AssetManager)
	UseItem(i Item)
}

func LoadPlayer(configs fs.FS, b Board, am assets.AssetManager) (Player, error) {
	var player PlayerInfo
	if err := readSceneConfig(configs, "player.yml", &player); err != nil {
		return nil, err
	}

	if err := validateEntity(&player.EntityInfo, &player.AttackInfo, "config/player.yml"); err != nil {
		return nil, err
	}
	if player.Experience < 0 || player.Level < 0 {
		return nil, fmt.Errorf("config/player.yml: experience and level must be nonnegative")
	}
	if err := validatePlacement(b, &player, "config/player.yml"); err != nil {
		return nil, err
	}
	for _, kind := range []assets.AssetType{"walk", "attack", "death"} {
		if err := requiredAsset(am, player.Character, kind, b.GetGridSize(), "config/player.yml"); err != nil {
			return nil, err
		}
	}
	b.AddObjectToBoard(&player)

	player.LoadImages(am)
	player.moveQueue = collections.NewQueue[image.Point]()

	return &player, nil
}

func (p *PlayerInfo) Draw(screen *ebiten.Image, b Board) error {
	pt := p.movePosition
	var x, y float32
	if pt != nil {
		x, y = float32(pt.X), float32(pt.Y)
	} else {
		var err error
		x, y, err = b.GridToXY(p.GridX, p.GridY)
		if err != nil {
			return err
		}
	}

	var c color.Color
	if p.Alive() {
		if p.CurrentHealth < 25 {
			c = color.RGBA{255, 255, 0, 255}
		} else {
			c = color.RGBA{0, 255, 0, 255}
		}
	} else {
		c = color.RGBA{128, 128, 128, 255}
	}
	cx := x + float32(b.GetGridSize()*p.Size)/2
	cy := y + float32(b.GetGridSize()*p.Size)/2
	r := float32(p.Size*b.GetGridSize()) / 2
	vector.DrawFilledCircle(screen, cx, cy, r, c, true)

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM = p.spriteTransform(x, y, b.GetGridSize())

	frame := 0
	img := p.walkAsset.GetImage()
	if p.Alive() {
		switch p.animState {
		case "walk":
			frame = p.walkFrame
		case "attack":
			frame = p.attackFrame
			img = p.attackAsset.GetImage()
		}
	} else {
		img = p.deathAsset.GetImage()
		frame = p.deathAsset.GetFrameCount() - 1
	}
	rect := image.Rect(frame*b.GetGridSize(), 0, (frame+1)*b.GetGridSize(), b.GetGridSize())
	screen.DrawImage(img.SubImage(rect).(*ebiten.Image), opts)

	// debug := fmt.Sprintf("frame=%d\nrect=%s\nx,y=%f,%f\ngx,gy=%d,%d\ncx,cy=%f,%f", frame, rect.String(), x, y, p.GridX, p.GridY, cx, cy)
	// ebitenutil.DebugPrintAt(screen, debug, 400, 10)

	p.DrawInfo(screen, 4, 4)
	return nil
}

// spriteTransform mirrors one cropped frame inside its existing footprint.
// Apply the reflection before world scaling/translation so position and HUD stay put.
func (p *PlayerInfo) spriteTransform(x, y float32, frameWidth int) ebiten.GeoM {
	var transform ebiten.GeoM
	if p.facingLeft {
		transform.Scale(-1, 1)
		transform.Translate(float64(frameWidth), 0)
	}
	size := float64(p.GetSize())
	transform.Scale(size, size)
	transform.Translate(float64(x), float64(y))
	return transform
}

func (p *PlayerInfo) AttackMonster(m Monster) {
	p.animState = "attack"
	p.Attack(m)

	if m.Alive() {
		// if monster is still alive calculate the monster's attack value and subtract from player's health
		m.Attack(p)
	}
}

func (p *PlayerInfo) AddExperience(xp int) {
	p.Experience += xp
	newLevel := p.Experience / 25
	if newLevel > p.Level {
		p.LevelUp(newLevel)
	}
}

func (p *PlayerInfo) LevelUp(newLevel int) {
	p.AttackPower++
	p.Level = newLevel
	fmt.Println("Level Up!")
}

func (p *PlayerInfo) DrawInfo(screen *ebiten.Image, x, y float32) {
	var infoText string
	if p.CurrentHealth > 0 {
		infoText = fmt.Sprintf("%s(%d)\n%d/%d\n%dg\n", p.Name, p.Level, p.CurrentHealth, p.MaxHealth, p.Gold)
	} else {
		infoText = fmt.Sprintf("Dead\n%s", p.Name)
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.LineSpacing = gameFontFace.Size * 1.5
	text.Draw(screen, infoText, gameFontFace, op)
}

func (p *PlayerInfo) Move(direction Direction, b Board) bool {
	// Accept another grid step only after fixed updates finish this one.
	if !p.moveQueue.IsEmpty() {
		return false
	}

	gx, gy := p.GridX, p.GridY
	switch direction {
	case Up:
		gy -= 1
	case Down:
		gy += 1
	case Right:
		gx += 1
	case Left:
		gx -= 1
	default:
		return false
	}

	if b.CanOccupySpace(p, gx, gy) {
		// Validate and queue the visual path before committing occupancy.
		if err := p.queuePoints(p.GridX, p.GridY, gx, gy, b); err != nil {
			return false
		}
		switch direction {
		case Left:
			p.facingLeft = true
		case Right:
			p.facingLeft = false
		}
		p.animState = "walk"

		b.RemoveObjectFromBoard(p)
		p.GridX = gx
		p.GridY = gy
		b.AddObjectToBoard(p)

		return true
	}

	return false
}

// One cell takes 12 fixed updates (200ms at the default 60 TPS).
const movementTicks = 12
const walkFrameTicks = 4

func (p *PlayerInfo) queuePoints(gx1, gy1, gx2, gy2 int, b Board) error {
	x1, y1, err := b.GridToXY(gx1, gy1)
	if err != nil {
		return err
	}
	x2, y2, err := b.GridToXY(gx2, gy2)
	if err != nil {
		return err
	}
	// Exclude the origin and include the exact destination. A fixed count avoids
	// direction/size/divisibility changing the duration or creating an endpoint jump.
	for tick := 1; tick <= movementTicks; tick++ {
		fraction := float64(tick) / movementTicks
		pt := image.Pt(int(math.Round(float64(x1)+(float64(x2)-float64(x1))*fraction)), int(math.Round(float64(y1)+(float64(y2)-float64(y1))*fraction)))
		p.EnqueueMove(&pt)
	}
	return nil
}

// AdvanceMovement runs exactly once per game update; Draw only observes state.
func (p *PlayerInfo) AdvanceMovement() {
	p.movePosition = p.DequeueMove()
	if p.movePosition == nil {
		p.walkTicks = 0
		return
	}
	if p.Alive() && p.animState == "walk" {
		p.walkTicks++
		if p.walkTicks == walkFrameTicks {
			incrementFrame(&p.walkFrame, p.walkAsset)
			p.walkTicks = 0
		}
	}
}

func (p *PlayerInfo) UseItem(i Item) {
	i.Use(p)
	p.animState = "attack"
	incrementFrame(&p.attackFrame, p.attackAsset)
}

func (p *PlayerInfo) Idle() {
	// p.animState = "idle"
}

func (p *PlayerInfo) GetExperience() int {
	return p.Experience
}

func (p *PlayerInfo) GetLevel() int {
	return p.Level
}

func (p *PlayerInfo) LoadImages(am assets.AssetManager) {
	p.walkAsset = am.GetAssetInfo(p.Character, "walk")
	p.deathAsset = am.GetAssetInfo(p.Character, "death")

	p.LoadAttackImage(am, p.Character)
}
func (p *PlayerInfo) EnqueueMove(pt *image.Point) {
	p.moveQueue.Enqueue(*pt)
}
func (p *PlayerInfo) DequeueMove() *image.Point {
	pt, success := p.moveQueue.TryDequeue()
	if success {
		return &pt
	}

	return nil
}
