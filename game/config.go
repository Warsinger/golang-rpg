package game

import (
	"fmt"
	"io/fs"
	"strings"

	"golang-rpg/assets"
	"golang-rpg/internal/configfile"
)

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

func requiredAsset(am assets.AssetManager, name string, kind assets.AssetType, gridSize int, context string) error {
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

func readSceneConfig(source fs.FS, name string, value any) error {
	if err := configfile.ReadFS(source, name, value); err != nil {
		return fmt.Errorf("config/%s: %w", name, err)
	}
	return nil
}
