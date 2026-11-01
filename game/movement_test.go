package game

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"golang-rpg/assets"
	"image"
	"testing"
)

func TestMovementInterpolationInEveryDirection(t *testing.T) {
	for _, size := range []int{1, 2, 3} {
		for _, grid := range []int{12, 48, 50} {
			for _, step := range []struct {
				name      string
				direction Direction
				delta     image.Point
			}{
				{"right", Right, image.Pt(1, 0)}, {"left", Left, image.Pt(-1, 0)}, {"down", Down, image.Pt(0, 1)}, {"up", Up, image.Pt(0, -1)},
			} {
				t.Run(fmt.Sprintf("%s/size%d/grid%d", step.name, size, grid), func(t *testing.T) {
					g := ruleScene()
					p := g.Player.(*PlayerInfo)
					b := g.Board.(*BoardInfo)
					b.RemoveObjectFromBoard(p)
					b.GridSize = grid
					p.GridX = 3
					p.GridY = 2
					p.Size = size
					b.AddObjectToBoard(p)
					origin := image.Pt(3*grid, 2*grid)
					want := origin.Add(step.delta.Mul(grid))
					if !p.Move(step.direction, b) {
						t.Fatal("open move rejected")
					}
					points := drainMoves(p)
					if len(points) != 12 {
						t.Fatalf("move took %d steps, want 12", len(points))
					}
					if points[len(points)-1] != want {
						t.Fatalf("last point %v, want exact endpoint %v", points[len(points)-1], want)
					}
					previous := origin
					for _, pt := range points {
						delta := pt.Sub(previous)
						distance := delta.X*step.delta.X + delta.Y*step.delta.Y
						if distance < 1 || distance > (grid+11)/12 || delta.X*step.delta.Y-delta.Y*step.delta.X != 0 {
							t.Fatalf("uneven/reversed step %v -> %v", previous, pt)
						}
						previous = pt
					}
				})
			}
		}
	}
}

type movementDrawAsset struct {
	*assets.AssetInfo
	img *ebiten.Image
}

func (a movementDrawAsset) GetImage() *ebiten.Image { return a.img }

func TestHeldMovementUsesUpdatesNotDraws(t *testing.T) {
	for _, size := range []int{1, 2, 3} {
		for _, draws := range []int{0, 1, 5} {
			for _, step := range []struct {
				name  string
				input gameInput
				delta image.Point
			}{
				{"right", gameInput{right: true}, image.Pt(1, 0)}, {"left", gameInput{left: true}, image.Pt(-1, 0)},
				{"down", gameInput{down: true}, image.Pt(0, 1)}, {"up", gameInput{up: true}, image.Pt(0, -1)},
			} {
				t.Run(fmt.Sprintf("%s/size%d/draws%d", step.name, size, draws), func(t *testing.T) {
					g := ruleScene()
					p := g.Player.(*PlayerInfo)
					b := &BoardInfo{GridWidth: 20, GridHeight: 20, GridSize: 48, occupied: make([]Object, 400)}
					g.Board = b
					p.GridX = 5
					p.GridY = 5
					p.Size = size
					b.AddObjectToBoard(p)
					img := ebiten.NewImage(48*6, 48)
					defer img.Deallocate()
					p.walkAsset = movementDrawAsset{&assets.AssetInfo{FrameCount: 6}, img}
					screen := ebiten.NewImage(960, 960)
					defer screen.Deallocate()
					origin := image.Pt(240, 240)
					for tick := 1; tick <= 24; tick++ {
						if err := g.applyInput(step.input); err != nil {
							t.Fatal(err)
						}
						want := origin.Add(step.delta.Mul(tick * 4))
						if p.movePosition == nil || *p.movePosition != want {
							t.Fatalf("tick %d position %v, want %v", tick, p.movePosition, want)
						}
						if p.walkFrame != (tick/4)%6 {
							t.Fatalf("tick %d frame %d, want %d", tick, p.walkFrame, (tick/4)%6)
						}
						point, frame, ticks := *p.movePosition, p.walkFrame, p.walkTicks
						for range draws {
							if err := p.Draw(screen, b); err != nil {
								t.Fatal(err)
							}
						}
						if *p.movePosition != point || p.walkFrame != frame || p.walkTicks != ticks {
							t.Fatal("drawing changed movement/animation")
						}
					}
					wantGrid := image.Pt(5, 5).Add(step.delta.Mul(2))
					if image.Pt(p.GridX, p.GridY) != wantGrid || !p.moveQueue.IsEmpty() {
						t.Fatal("two-cell duration differed by direction or draw frequency")
					}
					if err := g.applyInput(gameInput{}); err != nil {
						t.Fatal(err)
					}
					if p.movePosition != nil || p.walkFrame != 0 {
						t.Fatal("idle did not settle at destination and freeze frame")
					}
					for x := 0; x < 20; x++ {
						for y := 0; y < 20; y++ {
							occupied := x >= p.GridX && x < p.GridX+size && y >= p.GridY && y < p.GridY+size
							if (b.occupied[b.GridToIndex(x, y)] == p) != occupied {
								t.Fatal("large footprint occupancy drifted during interpolation")
							}
						}
					}
				})
			}
		}
	}
}

func TestHeldAttackPreservesFifteenActionsPerSecond(t *testing.T) {
	g := ruleScene()
	p := g.Player.(*PlayerInfo)
	p.CurrentHealth = 1000
	p.MaxHealth = 1000
	m := addMonster(g, 2, 1, 1000)
	for tick := 1; tick <= 60; tick++ {
		if err := g.applyInput(gameInput{attack: true}); err != nil {
			t.Fatal(err)
		}
		attacks := (tick + 3) / 4
		if m.CurrentHealth != 1000-4*attacks || p.CurrentHealth != 1000-2*attacks {
			t.Fatalf("tick %d attack/counterattack cadence changed: %d/%d", tick, m.CurrentHealth, p.CurrentHealth)
		}
	}
}

func TestResetDiscardsPendingMovementAndAnimation(t *testing.T) {
	g, err := NewWithResources(resourceSnapshot(t, "config"), resourceSnapshot(t, "assets"))
	if err != nil {
		t.Fatal(err)
	}
	p := g.Player.(*PlayerInfo)
	x, y := p.GridX, p.GridY
	if err := g.applyInput(gameInput{left: true}); err != nil {
		t.Fatal(err)
	}
	if p.movePosition == nil || p.moveQueue.IsEmpty() {
		t.Fatal("no pending movement to reset")
	}
	if err := g.applyInput(gameInput{reset: true}); err != nil {
		t.Fatal(err)
	}
	next := g.Player.(*PlayerInfo)
	if next == p || next.GridX != x || next.GridY != y || next.movePosition != nil || !next.moveQueue.IsEmpty() || next.walkFrame != 0 || next.walkTicks != 0 || next.facingLeft || g.attackCooldown != 0 {
		t.Fatal("reset retained runtime movement state")
	}
}
