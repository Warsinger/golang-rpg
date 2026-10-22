package main

import (
	"errors"
	"fmt"
	"image"
	"os"
	"strings"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func ruleScene() *GameInfo {
	b := &BoardInfo{GridWidth: 8, GridHeight: 6, GridSize: 48, occupied: make([]Object, 48)}
	walk, attack, death := &AssetInfo{FrameCount: 3}, &AssetInfo{FrameCount: 3}, &AssetInfo{FrameCount: 3}
	p := &PlayerInfo{EntityInfo: EntityInfo{ObjectInfo: ObjectInfo{GridX: 1, GridY: 1, Size: 1}, Name: "Player", MaxHealth: 100, CurrentHealth: 100, Gold: 10, Defense: 1}, AttackInfo: AttackInfo{AttackPower: 6, attackAsset: attack, randomIntN: func(n int) int { return n - 1 }}, Level: 1, walkAsset: walk, deathAsset: death, moveQueue: NewQueue[image.Point]()}
	am := &AssetManagerInfo{assets: map[string]map[AssetType]Asset{"Treasure": {"item": &AssetInfo{FrameCount: 1}}, "HealthPack": {"item": &AssetInfo{FrameCount: 1}}}}
	b.AddObjectToBoard(p)
	return &GameInfo{Board: b, Player: p, AssetManager: am}
}

func addMonster(g *GameInfo, x, y, hp int) *MonsterInfo {
	m := &MonsterInfo{EntityInfo: EntityInfo{ObjectInfo: ObjectInfo{GridX: x, GridY: y, Size: 1}, Name: "Monster", MaxHealth: hp, CurrentHealth: hp, Defense: 2, Gold: 12, ExperienceValue: 16}, AttackInfo: AttackInfo{AttackPower: 3, attackAsset: &AssetInfo{FrameCount: 3}, randomIntN: func(n int) int { return n - 1 }}}
	g.Board.AddObjectToBoard(m)
	g.Monsters = append(g.Monsters, m)
	return m
}

func drainMoves(p *PlayerInfo) []image.Point {
	var points []image.Point
	for pt := p.DequeueMove(); pt != nil; pt = p.DequeueMove() {
		points = append(points, *pt)
	}
	return points
}

func TestPickupRemovesUsedItemsAndPreservesNewOccupant(t *testing.T) {
	g := ruleScene()
	p := g.Player.(*PlayerInfo)
	item := NewTreasure(g.Board, g.AssetManager, 7, 2, 1, 1)
	g.Items = []Item{item}
	if err := g.applyInput(gameInput{use: true}); err != nil {
		t.Fatal(err)
	}
	if p.Gold != 17 || !item.Used {
		t.Fatal("treasure was not applied exactly once")
	}
	if len(g.Items) != 0 {
		t.Errorf("used item retained in scene: %d", len(g.Items))
	}
	if !p.Move(Right, g.Board) {
		t.Fatal("picked-up item still blocks movement")
	}
	drainMoves(p)
	if err := g.applyInput(gameInput{use: true}); err != nil {
		t.Fatal(err)
	}
	if g.Board.CanOccupySpace(&ObjectInfo{Size: 1}, 2, 1) {
		t.Error("reusing an old pickup cleared the player's occupied cell")
	}
	if p.Gold != 17 {
		t.Fatal("repeated pickup awarded duplicate gold")
	}
}

func TestPickupKeepsDistantItemsAndConsumesAllNearby(t *testing.T) {
	g := ruleScene()
	p := g.Player.(*PlayerInfo)
	p.CurrentHealth = 90
	treasure := NewTreasure(g.Board, g.AssetManager, 7, 2, 1, 1)
	health := &HealthPackInfo{ItemInfo: ItemInfo{ObjectInfo: ObjectInfo{GridX: 1, GridY: 2, Size: 1}, Value: 50}}
	g.Board.AddObjectToBoard(health)
	distant := NewTreasure(g.Board, g.AssetManager, 20, 6, 4, 1)
	g.Items = []Item{treasure, health, distant}
	if err := g.applyInput(gameInput{use: true}); err != nil {
		t.Fatal(err)
	}
	if p.Gold != 17 || p.CurrentHealth != 100 || len(g.Items) != 1 || g.Items[0] != distant {
		t.Fatalf("pickup state: gold=%d hp=%d items=%v", p.Gold, p.CurrentHealth, g.Items)
	}
	if !g.Board.CanOccupySpace(&ObjectInfo{Size: 1}, 2, 1) || !g.Board.CanOccupySpace(&ObjectInfo{Size: 1}, 1, 2) || g.Board.CanOccupySpace(&ObjectInfo{Size: 1}, 6, 4) {
		t.Fatal("pickup occupancy disagrees with remaining items")
	}
}

func TestAttackCounterattackDeathAndLoot(t *testing.T) {
	g := ruleScene()
	p := g.Player.(*PlayerInfo)
	m := addMonster(g, 2, 1, 8)
	if err := g.applyInput(gameInput{attack: true}); err != nil {
		t.Fatal(err)
	}
	if m.CurrentHealth != 4 || p.CurrentHealth != 98 || p.Experience != 0 || len(g.Items) != 0 {
		t.Fatal("surviving target did not counterattack with expected damage")
	}
	if err := g.applyInput(gameInput{attack: true}); err != nil {
		t.Fatal(err)
	}
	if m.Alive() || p.CurrentHealth != 98 || p.Experience != 16 || len(g.Items) != 1 {
		t.Fatal("kill should award XP/loot without a counterattack")
	}
	loot := g.Items[0]
	if loot.GetValue() != 12 || loot.GetGridX() != 2 || loot.GetGridY() != 1 || loot.GetSize() != 1 || loot.GetUsed() {
		t.Fatal("loot did not inherit the monster's gold and footprint")
	}
	if g.Board.CanOccupySpace(&ObjectInfo{Size: 1}, 2, 1) {
		t.Fatal("loot should occupy the dead monster's space")
	}
	if err := g.applyInput(gameInput{attack: true}); err != nil {
		t.Fatal(err)
	}
	if len(g.Items) != 1 || p.Experience != 16 || p.CurrentHealth != 98 {
		t.Fatal("dead monster was rewarded/attacked twice")
	}
	if err := g.applyInput(gameInput{use: true}); err != nil {
		t.Fatal(err)
	}
	if p.Gold != 22 || len(g.Items) != 0 || !g.Board.CanOccupySpace(&ObjectInfo{Size: 1}, 2, 1) {
		t.Fatal("loot pickup was not removed/applied")
	}
}

func TestAttackSkipsDistantMonstersAndHitsAllInRange(t *testing.T) {
	g := ruleScene()
	addMonster(g, 2, 1, 20)
	addMonster(g, 1, 2, 20)
	far := addMonster(g, 6, 4, 20)
	if err := g.applyInput(gameInput{attack: true}); err != nil {
		t.Fatal(err)
	}
	if g.Monsters[0].GetCurrentHealth() != 16 || g.Monsters[1].GetCurrentHealth() != 16 || far.CurrentHealth != 20 || g.Player.GetCurrentHealth() != 96 {
		t.Fatal("attack range or all-target counterattack behavior changed")
	}
}

func TestInputPriorityAndDeath(t *testing.T) {
	for _, test := range []struct {
		name         string
		input        gameInput
		wantX, wantY int
	}{
		{"right first", gameInput{right: true, left: true, down: true, up: true, attack: true}, 2, 1},
		{"left before down", gameInput{left: true, down: true, up: true, attack: true}, 0, 1},
		{"down before up", gameInput{down: true, up: true, attack: true}, 1, 2},
		{"up before attack", gameInput{up: true, attack: true}, 1, 0},
	} {
		t.Run(test.name, func(t *testing.T) {
			g := ruleScene()
			m := addMonster(g, 2, 2, 20)
			if err := g.applyInput(test.input); err != nil {
				t.Fatal(err)
			}
			if g.Player.GetGridX() != test.wantX || g.Player.GetGridY() != test.wantY || m.CurrentHealth != 20 {
				t.Fatal("held-key precedence changed")
			}
		})
	}
	t.Run("blocked direction still suppresses attack", func(t *testing.T) {
		g := ruleScene()
		m := addMonster(g, 2, 1, 20)
		if err := g.applyInput(gameInput{right: true, attack: true}); err != nil {
			t.Fatal(err)
		}
		if g.Player.GetGridX() != 1 || m.CurrentHealth != 20 {
			t.Fatal("blocked movement should not fall through to attack")
		}
	})
	t.Run("dead player cannot act", func(t *testing.T) {
		g := ruleScene()
		p := g.Player.(*PlayerInfo)
		p.CurrentHealth = 0
		m := addMonster(g, 2, 1, 20)
		item := NewTreasure(g.Board, g.AssetManager, 7, 1, 2, 1)
		g.Items = []Item{item}
		if err := g.applyInput(gameInput{right: true, attack: true, use: true}); err != nil {
			t.Fatal(err)
		}
		if p.GridX != 1 || m.CurrentHealth != 20 || item.Used || p.Gold != 10 {
			t.Fatal("dead player acted")
		}
	})
	t.Run("movement and pickup may occur in one tick", func(t *testing.T) {
		g := ruleScene()
		g.Items = []Item{NewTreasure(g.Board, g.AssetManager, 7, 3, 1, 1)}
		if err := g.applyInput(gameInput{right: true, use: true}); err != nil {
			t.Fatal(err)
		}
		if g.Player.GetGridX() != 2 || g.Player.GetGold() != 17 {
			t.Fatal("pickup must use the post-move position")
		}
	})
}

func TestResetQuitAndErrorPaths(t *testing.T) {
	t.Run("reset rebuilds configured scene", func(t *testing.T) {
		copyScene(t)
		g := ruleScene()
		g.Player.SetCurrentHealth(0)
		if err := g.applyInput(gameInput{reset: true}); err != nil {
			t.Fatal(err)
		}
		if g.Player.GetName() != "Warsinger" || g.Player.GetCurrentHealth() != 100 || len(g.Monsters) != 6 || len(g.Items) != 6 {
			t.Fatal("reset did not replace the old state")
		}
	})
	t.Run("quit writes snapshot even when dead", func(t *testing.T) {
		t.Chdir(t.TempDir())
		g := ruleScene()
		g.Player.SetCurrentHealth(0)
		if err := g.applyInput(gameInput{quit: true, right: true}); !errors.Is(err, ebiten.Termination) {
			t.Fatalf("quit result: %v", err)
		}
		if _, err := os.Stat("game_state.yml"); err != nil {
			t.Fatal(err)
		}
		if g.Player.GetGridX() != 1 {
			t.Fatal("quit allowed later movement")
		}
	})
	t.Run("save failure returned", func(t *testing.T) {
		t.Chdir(t.TempDir())
		if err := os.Mkdir("game_state.yml", 0755); err != nil {
			t.Fatal(err)
		}
		err := ruleScene().applyInput(gameInput{quit: true})
		if err == nil || !strings.Contains(err.Error(), "save before quit: write game_state.yml") {
			t.Fatalf("save error: %v", err)
		}
	})
	t.Run("reset failure precedes quit", func(t *testing.T) {
		t.Chdir(t.TempDir())
		g := ruleScene()
		player := g.Player
		err := g.applyInput(gameInput{reset: true, quit: true})
		if !errors.Is(err, os.ErrNotExist) || !strings.Contains(err.Error(), "reset game") {
			t.Fatalf("reset error: %v", err)
		}
		if g.Player != player {
			t.Fatal("failed reset replaced state")
		}
		if _, err := os.Stat("game_state.yml"); !errors.Is(err, os.ErrNotExist) {
			t.Fatal("quit ran after failed reset")
		}
	})
}

func TestMovementAnimationAndQueueCharacterization(t *testing.T) {
	for _, test := range []struct {
		direction Direction
		want      image.Point
		count     int
		last      image.Point
	}{
		{Right, image.Pt(2, 1), 12, image.Pt(92, 48)},
		{Down, image.Pt(1, 2), 12, image.Pt(48, 92)},
		{Left, image.Pt(0, 1), 0, image.Point{}},
		{Up, image.Pt(1, 0), 0, image.Point{}},
	} {
		t.Run(fmt.Sprint(test.direction), func(t *testing.T) {
			g := ruleScene()
			p := g.Player.(*PlayerInfo)
			if !p.Move(test.direction, g.Board) {
				t.Fatal("valid move blocked")
			}
			if p.GridX != test.want.X || p.GridY != test.want.Y || p.walkFrame != 1 || p.animState != "walk" {
				t.Fatal("movement/animation state incorrect")
			}
			if !g.Board.CanOccupySpace(&ObjectInfo{Size: 1}, 1, 1) || g.Board.CanOccupySpace(&ObjectInfo{Size: 1}, p.GridX, p.GridY) {
				t.Fatal("movement did not transfer occupancy")
			}
			if test.count > 0 && p.Move(test.direction, g.Board) {
				t.Fatal("queued animation should block another movement")
			}
			points := drainMoves(p)
			if len(points) != test.count {
				t.Fatalf("got %d interpolation points, want %d", len(points), test.count)
			}
			if test.count > 0 && (points[0] != image.Pt(48, 48) || points[len(points)-1] != test.last) {
				t.Fatalf("interpolation endpoints %v", points)
			}
			if p.DequeueMove() != nil {
				t.Fatal("empty move queue returned a point")
			}
			p.Idle()
			if p.animState != "walk" {
				t.Fatal("existing Idle should not reset the animation")
			}
		})
	}
	t.Run("collision queues no movement", func(t *testing.T) {
		g := ruleScene()
		p := g.Player.(*PlayerInfo)
		addMonster(g, 2, 1, 10)
		if p.Move(Right, g.Board) || !p.moveQueue.IsEmpty() || p.GridX != 1 || p.walkFrame != 0 {
			t.Fatal("collision moved the player or queued frames")
		}
	})
	t.Run("invalid interpolation endpoint returns error", func(t *testing.T) {
		g := ruleScene()
		p := g.Player.(*PlayerInfo)
		if err := p.queuePoints(-1, 1, 2, 1, g.Board); err == nil {
			t.Fatal("bad origin accepted")
		}
		if err := p.queuePoints(1, 1, 99, 1, g.Board); err == nil {
			t.Fatal("bad destination accepted")
		}
	})
}

func TestCounterattackCanKillAndLaterTicksDoNotAct(t *testing.T) {
	g := ruleScene()
	p := g.Player.(*PlayerInfo)
	p.CurrentHealth = 1
	m := addMonster(g, 2, 1, 20)
	if err := g.applyInput(gameInput{attack: true}); err != nil {
		t.Fatal(err)
	}
	if p.CurrentHealth != 0 || m.CurrentHealth != 16 {
		t.Fatal("lethal counterattack was not clamped to zero")
	}
	if err := g.applyInput(gameInput{attack: true}); err != nil {
		t.Fatal(err)
	}
	if m.CurrentHealth != 16 || p.Experience != 0 {
		t.Fatal("dead player attacked in a later tick")
	}
}

func TestAlreadyUsedItemCannotClearOtherObject(t *testing.T) {
	g := ruleScene()
	old := &TreasureInfo{ItemInfo: ItemInfo{ObjectInfo: ObjectInfo{GridX: 1, GridY: 1, Size: 1}, Value: 7, Used: true}}
	g.Items = []Item{old}
	if err := g.applyInput(gameInput{use: true}); err != nil {
		t.Fatal(err)
	}
	if len(g.Items) != 0 || g.Player.GetGold() != 10 || g.Board.CanOccupySpace(&ObjectInfo{Size: 1}, 1, 1) {
		t.Fatal("cleanup of an already-used item removed the player or awarded gold")
	}
}

func TestIndependentScenesAndIdleTicks(t *testing.T) {
	first, second := ruleScene(), ruleScene()
	if err := first.applyInput(gameInput{right: true}); err != nil {
		t.Fatal(err)
	}
	if err := second.applyInput(gameInput{}); err != nil {
		t.Fatal(err)
	}
	if first.Player.GetGridX() != 2 || second.Player.GetGridX() != 1 || second.Player.GetCurrentHealth() != 100 || second.Player.GetGold() != 10 || !second.Player.(*PlayerInfo).moveQueue.IsEmpty() {
		t.Fatal("idle tick or another scene changed player state")
	}
	if err := second.applyInput(gameInput{use: true}); err != nil {
		t.Fatal(err)
	}
	if len(second.Items) != 0 {
		t.Fatal("empty pickup pass created items")
	}
}

func TestLargePlayerMovesOverOwnFootprint(t *testing.T) {
	g := ruleScene()
	p := g.Player.(*PlayerInfo)
	p.Size = 2
	g.Board.AddObjectToBoard(p)
	if !p.Move(Right, g.Board) {
		t.Fatal("multi-cell movement rejected overlap with its own old cells")
	}
	for _, pt := range []image.Point{{2, 1}, {3, 1}, {2, 2}, {3, 2}} {
		if g.Board.CanOccupySpace(&ObjectInfo{Size: 1}, pt.X, pt.Y) {
			t.Fatalf("new player cell not occupied: %v", pt)
		}
	}
	for _, pt := range []image.Point{{1, 1}, {1, 2}} {
		if !g.Board.CanOccupySpace(&ObjectInfo{Size: 1}, pt.X, pt.Y) {
			t.Fatalf("old player cell not cleared: %v", pt)
		}
	}
}
