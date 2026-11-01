package assets

import (
	"golang.org/x/image/font/sfnt"
	"io"
	"testing"
)

func TestGameFontCoverage(t *testing.T) {
	data, err := io.ReadAll(FontReader())
	if err != nil {
		t.Fatal(err)
	}
	f, err := sfnt.Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	var buffer sfnt.Buffer
	name, err := f.Name(&buffer, sfnt.NameIDFamily)
	if err != nil {
		t.Fatal(err)
	}
	if name != "Press Start 2P" {
		t.Fatalf("unexpected embedded font family %q", name)
	}
	// Labels are built from configured names, health, gold, levels and item units.
	// All printable ASCII used by the shipped scene must have a real glyph.
	for r := rune(32); r <= 126; r++ {
		index, err := f.GlyphIndex(&buffer, r)
		if err != nil {
			t.Fatalf("glyph %q: %v", r, err)
		}
		if index == 0 {
			t.Errorf("missing glyph %q", r)
		}
	}
	// Each caller gets its own position, including after a previous reader was consumed.
	again, err := io.ReadAll(FontReader())
	if err != nil {
		t.Fatal(err)
	}
	if len(again) != len(data) {
		t.Fatalf("font reader was not reset: got %d bytes, want %d", len(again), len(data))
	}
}
