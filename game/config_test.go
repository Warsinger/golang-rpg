package game

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"image"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang-rpg/assets"
	"golang-rpg/internal/configfile"

	"go.yaml.in/yaml/v3"
)

func copyScene(t *testing.T) {
	t.Helper()
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	destination := t.TempDir()
	for _, folder := range []string{"config", "assets"} {
		err := filepath.WalkDir(filepath.Join(root, folder), func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			relative, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			target := filepath.Join(destination, relative)
			if entry.IsDir() {
				return os.MkdirAll(target, 0755)
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			return os.WriteFile(target, data, 0644)
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	t.Chdir(destination)
}

func rewriteConfig[T any](t *testing.T, path string, change func(*T)) {
	t.Helper()
	var value T
	if err := configfile.Read(path, &value); err != nil {
		t.Fatal(err)
	}
	change(&value)
	data, err := yaml.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
}

func TestShippedSceneAndReload(t *testing.T) {
	copyScene(t)
	g, err := New()
	if err != nil {
		t.Fatal(err)
	}
	if g.Board.GetWidth() != 1920 || g.Board.GetHeight() != 1440 || len(g.Monsters) != 6 || len(g.Items) != 6 || g.Player.GetName() != "Warsinger" {
		t.Fatal("shipped scene changed")
	}
	g.Player.SetCurrentHealth(1)
	if err := g.Save(); err != nil {
		t.Fatal(err)
	}
	if err := g.Load(); err != nil {
		t.Fatal(err)
	}
	if g.Player.GetCurrentHealth() != 100 || len(g.Monsters) != 6 || len(g.Items) != 6 {
		t.Fatal("reload must start from configuration, not saved state, without accumulating objects")
	}
}

func TestFailedReloadPreservesScene(t *testing.T) {
	copyScene(t)
	g := &GameInfo{}
	if err := g.Init(); err != nil {
		t.Fatal(err)
	}
	board, player, manager := g.Board, g.Player, g.AssetManager
	if err := os.WriteFile("config/monsters.yml", []byte("- null\n"), 0644); err != nil {
		t.Fatal(err)
	}
	err := g.Init()
	if err == nil || !strings.Contains(err.Error(), "config/monsters.yml[0]") {
		t.Fatalf("expected contextual reload failure, got %v", err)
	}
	if g.Board != board || g.Player != player || g.AssetManager != manager || len(g.Items) != 6 || len(g.Monsters) != 6 {
		t.Fatal("failed reload replaced the active scene")
	}
}

func TestInvalidStartupConfigurations(t *testing.T) {
	tests := []struct {
		name, context string
		change        func(*testing.T)
	}{
		{"missing file", "config/player.yml", func(t *testing.T) {
			if err := os.Remove("config/player.yml"); err != nil {
				t.Fatal(err)
			}
		}},
		{"unknown field", "config/board.yml", func(t *testing.T) {
			if err := os.WriteFile("config/board.yml", []byte("gridwidht: 4\n"), 0644); err != nil {
				t.Fatal(err)
			}
		}},
		{"malformed YAML", "config/board.yml", func(t *testing.T) {
			if err := os.WriteFile("config/board.yml", []byte("[\n"), 0644); err != nil {
				t.Fatal(err)
			}
		}},
		{"multiple documents", "config/board.yml", func(t *testing.T) {
			f, err := os.OpenFile("config/board.yml", os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			if _, err = f.WriteString("\n---\ngridwidth: 5\n"); err != nil {
				t.Fatal(err)
			}
		}},
		{"negative board", "config/board.yml", func(t *testing.T) { rewriteConfig(t, "config/board.yml", func(b *BoardInfo) { b.GridWidth = -1 }) }},
		{"oversized board", "config/board.yml", func(t *testing.T) {
			rewriteConfig(t, "config/board.yml", func(b *BoardInfo) { b.GridWidth = int(^uint(0) >> 1) })
		}},
		{"zero interpolation step", "config/board.yml", func(t *testing.T) { rewriteConfig(t, "config/board.yml", func(b *BoardInfo) { b.GridSize = 11 }) }},
		{"missing tile", "config/board.yml", func(t *testing.T) { rewriteConfig(t, "config/board.yml", func(b *BoardInfo) { b.TileSet = "Missing" }) }},
		{"invalid health", "config/player.yml", func(t *testing.T) {
			rewriteConfig(t, "config/player.yml", func(p *PlayerInfo) { p.CurrentHealth = p.MaxHealth + 1 })
		}},
		{"zero attack", "config/player.yml", func(t *testing.T) { rewriteConfig(t, "config/player.yml", func(p *PlayerInfo) { p.AttackPower = 0 }) }},
		{"negative level", "config/player.yml", func(t *testing.T) { rewriteConfig(t, "config/player.yml", func(p *PlayerInfo) { p.Level = -1 }) }},
		{"offboard player", "config/player.yml", func(t *testing.T) { rewriteConfig(t, "config/player.yml", func(p *PlayerInfo) { p.GridX = 40 }) }},
		{"zero size player", "config/player.yml", func(t *testing.T) { rewriteConfig(t, "config/player.yml", func(p *PlayerInfo) { p.Size = 0 }) }},
		{"missing character", "config/player.yml", func(t *testing.T) {
			rewriteConfig(t, "config/player.yml", func(p *PlayerInfo) { p.Character = "Missing" })
		}},
		{"sprite too small for cells", "config/player.yml", func(t *testing.T) {
			rewriteConfig(t, "config/assets.yml", func(a *[]*assets.AssetInfo) { (*a)[0].FrameCount = 7; (*a)[0].Size = 1 })
		}},
		{"null treasure", "config/treasures.yml[0]", func(t *testing.T) {
			if err := os.WriteFile("config/treasures.yml", []byte("- null\n"), 0644); err != nil {
				t.Fatal(err)
			}
		}},
		{"negative treasure", "config/treasures.yml[0]", func(t *testing.T) {
			rewriteConfig(t, "config/treasures.yml", func(items *[]*TreasureInfo) { (*items)[0].Value = -1 })
		}},
		{"colliding treasure", "config/treasures.yml[0]", func(t *testing.T) {
			rewriteConfig(t, "config/treasures.yml", func(items *[]*TreasureInfo) { (*items)[0].GridX = 1; (*items)[0].GridY = 1 })
		}},
		{"null health pack", "config/healthpacks.yml[0]", func(t *testing.T) {
			if err := os.WriteFile("config/healthpacks.yml", []byte("- null\n"), 0644); err != nil {
				t.Fatal(err)
			}
		}},
		{"offboard health pack", "config/healthpacks.yml[0]", func(t *testing.T) {
			rewriteConfig(t, "config/healthpacks.yml", func(items *[]*HealthPackInfo) { (*items)[0].GridY = 30 })
		}},
		{"null monster", "config/monsters.yml[0]", func(t *testing.T) {
			if err := os.WriteFile("config/monsters.yml", []byte("- null\n"), 0644); err != nil {
				t.Fatal(err)
			}
		}},
		{"invalid monster attack", "config/monsters.yml[0]", func(t *testing.T) {
			rewriteConfig(t, "config/monsters.yml", func(m *[]*MonsterInfo) { (*m)[0].AttackPower = 0 })
		}},
		{"colliding monsters", "config/monsters.yml[1]", func(t *testing.T) {
			rewriteConfig(t, "config/monsters.yml", func(m *[]*MonsterInfo) { (*m)[1].GridX = (*m)[0].GridX; (*m)[1].GridY = (*m)[0].GridY })
		}},
		{"null asset", "config/assets.yml[0]", func(t *testing.T) {
			if err := os.WriteFile("config/assets.yml", []byte("- null\n"), 0644); err != nil {
				t.Fatal(err)
			}
		}},
		{"duplicate asset", "duplicate asset", func(t *testing.T) {
			rewriteConfig(t, "config/assets.yml", func(a *[]*assets.AssetInfo) { *a = append(*a, (*a)[0]) })
		}},
		{"bad asset type", "unsupported asset type", func(t *testing.T) {
			rewriteConfig(t, "config/assets.yml", func(a *[]*assets.AssetInfo) { (*a)[0].Type = "mystery" })
		}},
		{"zero frames", "framecount", func(t *testing.T) {
			rewriteConfig(t, "config/assets.yml", func(a *[]*assets.AssetInfo) { (*a)[0].FrameCount = 0 })
		}},
		{"asset traversal", "local names", func(t *testing.T) {
			rewriteConfig(t, "config/assets.yml", func(a *[]*assets.AssetInfo) { (*a)[0].FilePath = "../walk.png" })
		}},
		{"missing image", "walk.png", func(t *testing.T) {
			if err := os.Remove("assets/characters/Warrior/walk.png"); err != nil {
				t.Fatal(err)
			}
		}},
		{"corrupt image", "decode asset", func(t *testing.T) {
			if err := os.WriteFile("assets/characters/Warrior/walk.png", []byte("not PNG"), 0644); err != nil {
				t.Fatal(err)
			}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			copyScene(t)
			test.change(t)
			g := &GameInfo{}
			err := g.Init()
			if err == nil || !strings.Contains(err.Error(), test.context) {
				t.Fatalf("want %q error, got %v", test.context, err)
			}
		})
	}
}

func TestSaveFailureIncludesPath(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := os.Mkdir("game_state.yml", 0755); err != nil {
		t.Fatal(err)
	}
	if err := (&GameInfo{}).Save(); err == nil || !strings.Contains(err.Error(), "write game_state.yml") {
		t.Fatalf("expected save path context, got %v", err)
	}
}

func TestOversizedImageRejectedBeforeDecode(t *testing.T) {
	for _, size := range []image.Point{{8193, 1}, {4096, 4097}} {
		t.Run(size.String(), func(t *testing.T) {
			copyScene(t)
			var data bytes.Buffer
			if err := png.Encode(&data, image.NewRGBA(image.Rect(0, 0, 1, 1))); err != nil {
				t.Fatal(err)
			}
			header := data.Bytes()
			binary.BigEndian.PutUint32(header[16:20], uint32(size.X))
			binary.BigEndian.PutUint32(header[20:24], uint32(size.Y))
			binary.BigEndian.PutUint32(header[29:33], crc32.ChecksumIEEE(header[12:29]))
			if err := os.WriteFile("assets/characters/Warrior/walk.png", header, 0644); err != nil {
				t.Fatal(err)
			}
			if err := (&GameInfo{}).Init(); err == nil || !strings.Contains(err.Error(), "image must stay within") {
				t.Fatalf("expected pre-decode dimensions error, got %v", err)
			}
		})
	}
}
