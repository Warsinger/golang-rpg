package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"gopkg.in/yaml.v3"
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

func LoadAssets() (AssetManager, error) {
	am := &AssetManagerInfo{}

	// as := []*AssetInfo{{Name: "bob", FrameCount: 1, Size: 44}, {Name: "bobie", FrameCount: 2, Size: 34}}
	// data, err := yaml.Marshal(as)
	// if err != nil {
	// 	return nil, err
	// }
	// err = os.WriteFile("game_assets.yml", data, 0644)
	// if err != nil {
	// 	return nil, err
	// }

	yamlFile, err := os.ReadFile("config/assets.yml")
	// fmt.Println(string(yamlFile))
	if err != nil {
		return nil, err
	}
	var assets []*AssetInfo
	err = yaml.Unmarshal(yamlFile, &assets)
	if err != nil {
		return nil, err
	}
	am.assets = make(map[string]map[AssetType]Asset)
	for _, a := range assets {
		fmt.Printf("asset info: %v\n", a)
		if am.assets[a.Name] == nil {
			am.assets[a.Name] = make(map[AssetType]Asset)

		}
		am.assets[a.Name][a.Type] = a
		err = loadImageAsset(a)
		if err != nil {
			return nil, err
		}
	}

	return am, nil
}

func loadImageAsset(a *AssetInfo) error {
	var typePath string
	switch a.Type {
	case "walk", "attack", "death":
		typePath = fmt.Sprintf("characters/%s", a.Name)
	case "item":
		typePath = "items"
	}

	filepath := fmt.Sprintf("assets/%s/%s", typePath, a.FilePath)
	img, _, err := ebitenutil.NewImageFromFile(filepath)
	if err != nil {
		log.Fatalf("failed to load asset %v: %v", a, err)
		return err
	}

	a.image = img
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
