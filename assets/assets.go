// Package assets owns the image registry and embedded font resource.
package assets

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"io/fs"
	"path"
	"strings"

	"golang-rpg/internal/configfile"

	"github.com/hajimehoshi/ebiten/v2"
)

type AssetType string
type AssetInfo struct {
	image      *ebiten.Image
	FrameCount int
	Size       int
	Type       AssetType
	FilePath   string
	Name       string
}
type Asset interface {
	GetImage() *ebiten.Image
	GetFrameCount() int
	GetSize() int
	GetType() AssetType
	GetFilePath() string
	GetName() string
}

type AssetManagerInfo struct {
	assets map[string]map[AssetType]Asset
}

type AssetManager interface {
	GetAssetInfo(name string, assetType AssetType) Asset
}

func (am *AssetManagerInfo) GetAssetInfo(name string, assetType AssetType) Asset {
	return am.assets[name][assetType]
}

func LoadAssets(configs, images fs.FS) (AssetManager, error) {
	am := &AssetManagerInfo{}

	var assets []*AssetInfo
	if err := configfile.ReadFS(configs, "assets.yml", &assets); err != nil {
		return nil, fmt.Errorf("config/assets.yml: %w", err)
	}
	am.assets = make(map[string]map[AssetType]Asset)
	for index, a := range assets {
		if err := validateAsset(a, index); err != nil {
			return nil, err
		}
		if am.assets[a.Name][a.Type] != nil {
			return nil, fmt.Errorf("config/assets.yml[%d]: duplicate asset %q/%s", index, a.Name, a.Type)
		}
		if am.assets[a.Name] == nil {
			am.assets[a.Name] = make(map[AssetType]Asset)

		}
		am.assets[a.Name][a.Type] = a
		err := loadImageAsset(images, a)
		if err != nil {
			return nil, err
		}
	}

	return am, nil
}

func loadImageAsset(images fs.FS, a *AssetInfo) error {
	var typePath string
	switch a.Type {
	case "walk", "attack", "death":
		typePath = fmt.Sprintf("characters/%s", a.Name)
	case "item":
		typePath = "items"
	case "tile":
		typePath = "tiles"
	}

	assetPath := path.Join(typePath, a.FilePath)
	data, err := fs.ReadFile(images, assetPath)
	if err != nil {
		return fmt.Errorf("asset %q/%s: %w", a.Name, a.Type, err)
	}
	metadata, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("decode asset %s: %w", assetPath, err)
	}
	if metadata.Width > 8192 || metadata.Height > 8192 || metadata.Width*metadata.Height > 16_000_000 {
		return fmt.Errorf("asset %s: image must stay within 8192 pixels per dimension and 16000000 pixels total", assetPath)
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("load asset %s: %w", assetPath, err)
	}

	a.image = ebiten.NewImageFromImage(img)
	return nil
}

func (a *AssetInfo) GetImage() *ebiten.Image {
	return a.image
}

func (a *AssetInfo) GetFrameCount() int {
	return a.FrameCount
}

func (a *AssetInfo) GetSize() int {
	return a.Size
}

func (a *AssetInfo) GetType() AssetType {
	return a.Type
}

func (a *AssetInfo) GetFilePath() string {
	return a.FilePath
}

func (a *AssetInfo) GetName() string {
	return a.Name
}

func validateAsset(a *AssetInfo, index int) error {
	context := fmt.Sprintf("config/assets.yml[%d]", index)
	if a == nil {
		return fmt.Errorf("%s: null asset", context)
	}
	if a.Name == "" || a.Name == "." || a.Name == ".." || strings.ContainsAny(a.Name, "/\\:") || a.FilePath == "" || a.FilePath == "." || a.FilePath == ".." || strings.ContainsAny(a.FilePath, "/\\:") {
		return fmt.Errorf("%s: name and filepath must be nonempty local names without directory separators", context)
	}
	if a.Size <= 0 || a.FrameCount <= 0 {
		return fmt.Errorf("%s: size and framecount must be positive", context)
	}
	switch a.Type {
	case "walk", "attack", "death", "item", "tile":
	default:
		return fmt.Errorf("%s: unsupported asset type %q", context, a.Type)
	}
	return nil
}
