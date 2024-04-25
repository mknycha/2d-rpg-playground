package assets

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var assets map[string]*ebiten.Image

func init() {
	tilesFileContent, err := ioutil.ReadFile("./assets/files/zoria_msx.png")
	if err != nil {
		log.Fatal(err)
	}
	tilesImg, _, err := image.Decode(bytes.NewReader(tilesFileContent))
	if err != nil {
		log.Fatal("failed to decode:", err)
	}
	zoriaImage := ebiten.NewImageFromImage(tilesImg)
	fontFileContent, err := ioutil.ReadFile("./assets/files/font-1.png")
	if err != nil {
		log.Fatal(err)
	}
	fontImg, _, err := image.Decode(bytes.NewReader(fontFileContent))
	if err != nil {
		log.Fatal("failed to decode:", err)
	}
	fontImage := ebiten.NewImageFromImage(fontImg)
	assets = map[string]*ebiten.Image{
		"zoria_msx": zoriaImage,
		"font":      fontImage,
	}
}

func GetAsset(name string) (*ebiten.Image, error) {
	asset, ok := assets[name]
	if !ok {
		return nil, fmt.Errorf("asset '%s' could not be found", name)
	}
	return asset, nil
}
