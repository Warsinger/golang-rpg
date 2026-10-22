package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"go.yaml.in/yaml/v3"
)

// readConfig accepts one strict YAML document and keeps filesystem failures wrapped.
func readConfig(path string, value any) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("read %s: %w (run from the repository root)", path, err)
	}
	defer f.Close()
	const maxConfigBytes = 1 << 20
	data, err := io.ReadAll(io.LimitReader(f, maxConfigBytes+1))
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if len(data) > maxConfigBytes {
		return fmt.Errorf("%s: configuration exceeds 1 MiB", path)
	}
	decoder := yaml.NewDecoder(strings.NewReader(string(data)))
	decoder.KnownFields(true)
	if err := decoder.Decode(value); err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return fmt.Errorf("%s: expected a single YAML document", path)
		}
		return fmt.Errorf("decode %s: %w", path, err)
	}
	return nil
}

func validateBoard(b *BoardInfo) error {
	// Movement interpolation divides GridSize by 12; smaller cells would loop forever.
	const maxPixels = 8192
	if b.GridWidth <= 0 || b.GridHeight <= 0 || b.GridSize < 12 || b.GridSize > maxPixels {
		return fmt.Errorf("config/board.yml: gridwidth/gridheight must be positive and gridsize must be 12..%d", maxPixels)
	}
	if b.GridWidth > maxPixels/b.GridSize || b.GridHeight > maxPixels/b.GridSize {
		return fmt.Errorf("config/board.yml: board viewport must not exceed %d pixels per dimension", maxPixels)
	}
	if b.TileSet == "" {
		return fmt.Errorf("config/board.yml: tileset is required")
	}
	return nil
}

func validatePlacement(b Board, o Object, context string) error {
	if o.GetSize() <= 0 || o.GetSize() > b.GetGridWidth() || o.GetSize() > b.GetGridHeight() ||
		o.GetGridX() < 0 || o.GetGridY() < 0 || o.GetGridX() > b.GetGridWidth()-o.GetSize() || o.GetGridY() > b.GetGridHeight()-o.GetSize() {
		return fmt.Errorf("%s: object must have positive size and fit inside the board", context)
	}
	if !b.CanOccupySpace(o, o.GetGridX(), o.GetGridY()) {
		return fmt.Errorf("%s: object overlaps an occupied cell", context)
	}
	return nil
}

func validateEntity(e *EntityInfo, a *AttackInfo, context string) error {
	if strings.TrimSpace(e.Name) == "" || e.MaxHealth <= 0 || e.CurrentHealth < 0 || e.CurrentHealth > e.MaxHealth || e.Defense < 0 || e.Gold < 0 || e.ExperienceValue < 0 || a.AttackPower <= 0 {
		return fmt.Errorf("%s: name, positive maxhealth/attackpower, currenthealth within 0..maxhealth, and nonnegative defense/gold/experiencevalue are required", context)
	}
	return nil
}

func requiredAsset(am AssetManager, name string, kind AssetType, gridSize int, context string) error {
	a := am.GetAssetInfo(name, kind)
	if a == nil || a.GetImage() == nil {
		return fmt.Errorf("%s: missing asset %q of type %q", context, name, kind)
	}
	// Character drawing crops frames using the board cell size, not asset.Size.
	if kind == "walk" || kind == "attack" || kind == "death" {
		bounds := a.GetImage().Bounds()
		if a.GetFrameCount() <= 0 || gridSize > bounds.Dy() || gridSize > bounds.Dx()/a.GetFrameCount() {
			return fmt.Errorf("%s: asset %q/%s is too small for %d frames of %d-pixel cells", context, name, kind, a.GetFrameCount(), gridSize)
		}
	}
	return nil
}
