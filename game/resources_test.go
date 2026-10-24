package game

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func resourceSnapshot(t *testing.T, folder string) fstest.MapFS {
	t.Helper()
	source := os.DirFS(filepath.Join("..", folder))
	files := fstest.MapFS{}
	if err := fs.WalkDir(source, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		data, err := fs.ReadFile(source, path)
		if err != nil {
			return err
		}
		files[path] = &fstest.MapFile{Data: data}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	return files
}
func TestInjectedResourcesLoadAndResetWithoutDisk(t *testing.T) {
	configs, images := resourceSnapshot(t, "config"), resourceSnapshot(t, "assets")
	t.Chdir(t.TempDir())
	g, err := NewWithResources(configs, images)
	if err != nil {
		t.Fatal(err)
	}
	if g.Board.GetWidth() != 1920 || len(g.Monsters) != 6 || len(g.Items) != 6 {
		t.Fatal("in-memory scene differs from shipped configuration")
	}
	g.Player.SetCurrentHealth(1)
	if err := g.applyInput(gameInput{reset: true}); err != nil {
		t.Fatal(err)
	}
	if g.Player.GetCurrentHealth() != 100 || len(g.Items) != 6 {
		t.Fatal("reset lost its injected resource sources")
	}
	if _, err := os.Stat("config"); !os.IsNotExist(err) {
		t.Fatal("test must not depend on working-directory configuration")
	}
}
func TestInjectedResourceFailuresPreserveActiveScene(t *testing.T) {
	configs, images := resourceSnapshot(t, "config"), resourceSnapshot(t, "assets")
	g, err := NewWithResources(configs, images)
	if err != nil {
		t.Fatal(err)
	}
	board, player := g.Board, g.Player
	delete(images, "characters/Warrior/walk.png")
	if err := g.Load(); err == nil {
		t.Fatal("missing injected image accepted")
	}
	if g.Board != board || g.Player != player {
		t.Fatal("failed resource reload replaced active scene")
	}
	if _, err := NewWithResources(nil, images); err == nil {
		t.Fatal("missing config source accepted")
	}
	if _, err := NewWithResources(configs, nil); err == nil {
		t.Fatal("missing image source accepted")
	}
}
