package main

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"gopkg.in/yaml.v3"
)

type AnimState string

type PlayerInfo struct {
	EntityInfo
	AttackInfo
	Experience int
	Level      int
	Character  string

	walkAsset Asset
	walkFrame int

	deathAsset Asset

	animState AnimState
	moveQueue *Queue[image.Point]
}

type Player interface {
	Entity
	Attacker

	Draw(screen *ebiten.Image, b Board)

	GetExperience() int
	GetLevel() int
	AttackMonster(m Monster)
	AddExperience(xp int)
	LevelUp(newLevel int)
	Move(direction Direction, b Board) bool
	EnqueueMove(pt *image.Point)
	DequeueMove() *image.Point
	Idle()
	LoadImages(am AssetManager)
	UseItem(i Item)
}

func LoadPlayer(b Board, am AssetManager) (Player, error) {
	yamlFile, err := os.ReadFile("config/player.yml")
	if err != nil {
		return nil, err
	}
	var player PlayerInfo
	err = yaml.Unmarshal(yamlFile, &player)
	if err != nil {
		return nil, err
	}

	b.AddObjectToBoard(&player)

	player.LoadImages(am)
	player.moveQueue = NewQueue[image.Point]()

	return &player, nil
}

func (p *PlayerInfo) Draw(screen *ebiten.Image, b Board) {
	pt := p.DequeueMove()
	var x, y float32
	incFrame := false
	if pt != nil {
		x, y = float32(pt.X), float32(pt.Y)
		incFrame = true
	} else {
		x, y = b.GridToXY(p.GridX, p.GridY)
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

	size := float64(p.GetSize())
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(x)/size, float64(y)/size)
	opts.GeoM.Scale(size, size)

	frame := 0
	img := p.walkAsset.GetImage()
	if p.Alive() {
		switch p.animState {
		case "walk":
			if incFrame {
				incrementFrame(&p.walkFrame, p.walkAsset)
			}
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
	op.LineSpacing = mplusNormalFace.Size * 1.5
	text.Draw(screen, infoText, mplusNormalFace, op)
}

func (p *PlayerInfo) Move(direction Direction, b Board) bool {
	// don't move if we're already moving to let the queue drain
	if p.moveQueue.Size() > 0 {
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
	}

	p.animState = "walk"

	if b.CanOccupySpace(p, gx, gy) {
		// queue up several coordinates between old and new to draw intermeiate frames
		p.queuePoints(p.GridX, p.GridY, gx, gy, b)

		b.RemoveObjectFromBoard(p)
		p.GridX = gx
		p.GridY = gy
		b.AddObjectToBoard(p)

		incrementFrame(&p.walkFrame, p.walkAsset)
		return true
	}

	return false
}

func (p *PlayerInfo) queuePoints(gx1, gy1, gx2, gy2 int, b Board) {
	const stepsPerFrame = 12
	x1, y1 := b.GridToXY(gx1, gy1)
	x2, y2 := b.GridToXY(gx2, gy2)
	stepSize := b.GetGridSize() / stepsPerFrame
	var pt image.Point
	if gx1 != gx2 {
		for i := int(x1); i < int(x2); i += stepSize {
			pt = image.Pt(i, int(y1))
			p.EnqueueMove(&pt)
		}
	} else if gy1 != gy2 {
		for i := int(y1); i < int(y2); i += stepSize {
			pt = image.Pt(int(x1), i)
			p.EnqueueMove(&pt)
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

func (p *PlayerInfo) LoadImages(am AssetManager) {
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
