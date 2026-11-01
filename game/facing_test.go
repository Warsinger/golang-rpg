package game

import (
	"fmt"
	"testing"
)

func TestPlayerFacingFollowsSuccessfulHorizontalMovement(t *testing.T) {
	g := ruleScene()
	p := g.Player.(*PlayerInfo)
	if p.facingLeft {
		t.Fatal("new player should face right like the source art")
	}
	for _, step := range []struct {
		direction Direction
		left      bool
	}{
		{Left, true}, {Up, true}, {Down, true}, {Right, false},
	} {
		if !p.Move(step.direction, g.Board) {
			t.Fatalf("move %v was blocked", step.direction)
		}
		drainMoves(p)
		if p.facingLeft != step.left {
			t.Fatalf("after %v: facingLeft=%v, want %v", step.direction, p.facingLeft, step.left)
		}
	}
	if !p.Move(Left, g.Board) {
		t.Fatal("left move was blocked")
	}
	p.Idle()
	if !p.facingLeft {
		t.Fatal("idle changed facing")
	}
	monster := addMonster(g, 1, 1, 100)
	p.AttackMonster(monster)
	if !p.facingLeft {
		t.Fatal("attacking changed facing")
	}
	p.UseItem(&HealthPackInfo{ItemInfo: ItemInfo{Value: 1}})
	if !p.facingLeft {
		t.Fatal("pickup changed facing")
	}
	p.CurrentHealth = 0
	if err := g.applyInput(gameInput{right: true}); err != nil {
		t.Fatal(err)
	}
	if !p.facingLeft {
		t.Fatal("death changed the last facing")
	}
}

func TestBlockedAndQueuedMovementPreserveFacing(t *testing.T) {
	g := ruleScene()
	p := g.Player.(*PlayerInfo)
	if !p.Move(Left, g.Board) {
		t.Fatal("left move was blocked")
	}
	if p.Move(Left, g.Board) || !p.facingLeft || p.GridX != 0 {
		t.Fatal("boundary attempt changed facing/position")
	}
	addMonster(g, 1, 1, 100)
	if p.Move(Right, g.Board) || !p.facingLeft {
		t.Fatal("collision turned the player without moving")
	}
	if !p.Move(Down, g.Board) || p.moveQueue.IsEmpty() {
		t.Fatal("expected queued vertical move")
	}
	if p.Move(Right, g.Board) || !p.facingLeft {
		t.Fatal("ignored input during interpolation changed facing")
	}
}

func TestSpriteReflectionStaysInsideScaledFrame(t *testing.T) {
	for _, size := range []int{1, 2, 3} {
		for _, left := range []bool{false, true} {
			for _, origin := range [][2]float32{{0, 0}, {113.5, 92.25}} {
				t.Run(fmt.Sprintf("size=%d/left=%v/origin=%v", size, left, origin), func(t *testing.T) {
					p := &PlayerInfo{EntityInfo: EntityInfo{ObjectInfo: ObjectInfo{Size: size}}, facingLeft: left}
					transform := p.spriteTransform(origin[0], origin[1], 48)
					x, y := float64(origin[0]), float64(origin[1])
					width := float64(48 * size)
					expectedLeft, expectedRight := x, x+width
					if left {
						expectedLeft, expectedRight = expectedRight, expectedLeft
					}
					for _, point := range [][4]float64{{0, 0, expectedLeft, y}, {48, 0, expectedRight, y}, {0, 48, expectedLeft, y + width}, {48, 48, expectedRight, y + width}, {24, 24, x + width/2, y + width/2}} {
						gotX, gotY := transform.Apply(point[0], point[1])
						if gotX != point[2] || gotY != point[3] {
							t.Fatalf("source (%v,%v) -> (%v,%v), want (%v,%v)", point[0], point[1], gotX, gotY, point[2], point[3])
						}
					}
				})
			}
		}
	}
}

func TestSceneRestartRestoresDefaultFacing(t *testing.T) {
	copyScene(t)
	g, err := New()
	if err != nil {
		t.Fatal(err)
	}
	p := g.Player.(*PlayerInfo)
	if !p.Move(Left, g.Board) || !p.facingLeft {
		t.Fatal("could not turn the shipped player left")
	}
	if err := g.applyInput(gameInput{reset: true}); err != nil {
		t.Fatal(err)
	}
	if g.Player.(*PlayerInfo).facingLeft {
		t.Fatal("restart retained the previous scene's facing")
	}
}
