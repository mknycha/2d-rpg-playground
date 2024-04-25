package main

import (
	"bytes"
	"image"
	"image/color"
	_ "image/png"
	"io/ioutil"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mknycha/rpg/assets"
	"github.com/mknycha/rpg/level_map"
)

const (
	ScreenWidth                     = 256
	ScreenHeight                    = 200
	playerVerticalMoveAcc           = 0.06
	playerHorizontalMoveAccOnGround = 0.011
	playerHorizontalMoveAccInAir    = 0.005
	playerVerticalVelocityMax       = 1.0
	playerHorizontalDrag            = playerHorizontalMoveAccOnGround * 4
	playerHorizontalVelocityMax     = 0.1
	clampHorizontalVelocityBelow    = 0.01
	gravity                         = 0.012
	spriteWidth                     = 16
	spriteHeight                    = 16
)

var cameraPosX = 0.0
var cameraPosY = 0.0

var playerPosX = 0.0
var playerPosY = 0.0
var playerVelX = 0.0
var playerVelY = 0.0

var playerFacingRight bool
var playerCurrentFrame image.Image
var time int

var (
	yellow           = color.NRGBA{0xff, 0xff, 0x0, 0xff}
	red              = color.NRGBA{0xff, 0x0, 0x0, 0xff}
	lightBlue        = color.NRGBA{0x0, 0xff, 0xff, 0xff}
	green            = color.NRGBA{0x0, 0xff, 0x0, 0xff}
	greenTransparent = color.NRGBA{0x0, 150, 0x0, 250}
)

func DrawText(screen *ebiten.Image, text string, x int, y int) {
	fontImage, err := assets.GetAsset("font")
	if err != nil {
		log.Fatal(err)
	}
	for i, r := range []byte(text) {
		r = r - 32
		sx := int(r % 18)
		sy := int(r / 18)
		characterWidth := 7
		characterHeight := 9
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x+(i*characterWidth)), float64(y))
		img := fontImage.SubImage(image.Rect(characterWidth*sx, characterHeight*sy, characterWidth*(sx+1), characterHeight*(sy+1)))
		screen.DrawImage(ebiten.NewImageFromImage(img), op)
	}
}

type Game struct {
	levelMap *level_map.Map
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// playerVelY = 0
	// playerVelX = 0

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		playerVelY += -playerHorizontalMoveAccOnGround
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		playerVelY += playerHorizontalMoveAccOnGround
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		playerFacingRight = false
		playerVelX += -playerHorizontalMoveAccOnGround
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		playerFacingRight = true
		playerVelX += playerHorizontalMoveAccOnGround
	}

	// drag
	playerVelX += -playerHorizontalDrag * playerVelX
	if math.Abs(playerVelX) < clampHorizontalVelocityBelow {
		playerVelX = 0
	}
	playerVelY += -playerHorizontalDrag * playerVelY
	if math.Abs(playerVelY) < clampHorizontalVelocityBelow {
		playerVelY = 0
	}

	newPlayerPosX := playerPosX + playerVelX
	newPlayerPosY := playerPosY + playerVelY

	// Collission
	if playerVelX <= 0 { // going left
		if g.levelMap.GetSolid(int(newPlayerPosX+0.0), int(playerPosY+0.0)) || g.levelMap.GetSolid(int(newPlayerPosX+0.0), int(playerPosY+0.9)) {
			newPlayerPosX = float64(int(newPlayerPosX + 1))
			playerVelX = 0
		}
	} else if playerVelX > 0 {
		if g.levelMap.GetSolid(int(newPlayerPosX+1.0), int(playerPosY+0.0)) || g.levelMap.GetSolid(int(newPlayerPosX+1.0), int(playerPosY+0.9)) {
			newPlayerPosX = float64(int(newPlayerPosX))
			playerVelX = 0
		}
	}

	if playerVelY <= 0 {
		if g.levelMap.GetSolid(int(newPlayerPosX), int(newPlayerPosY+0.0)) || g.levelMap.GetSolid(int(newPlayerPosX+0.9), int(newPlayerPosY+0.0)) {
			newPlayerPosY = float64(int(newPlayerPosY) + 1)
			playerVelY = 0
		}
	} else if playerVelY > 0 {
		if g.levelMap.GetSolid(int(newPlayerPosX), int(newPlayerPosY+1.0)) || g.levelMap.GetSolid(int(newPlayerPosX+0.9), int(newPlayerPosY+1.0)) {
			newPlayerPosY = float64(int(newPlayerPosY))
			playerVelY = 0
		}
	}

	// clamp velocities
	if playerVelX > playerHorizontalVelocityMax {
		playerVelX = playerHorizontalVelocityMax
	}
	if playerVelX < -playerHorizontalVelocityMax {
		playerVelX = -playerHorizontalVelocityMax
	}
	if playerVelY > playerVerticalVelocityMax {
		playerVelY = playerVerticalVelocityMax
	}
	if playerVelY < -playerVerticalVelocityMax {
		playerVelY = -playerVerticalVelocityMax
	}

	// animation
	time++
	if playerVelX != 0 || playerVelY != 0 {
		playerCurrentFrame = characterRunning[time/3%len(characterRunning)]
	} else {
		playerCurrentFrame = characterStanding
	}

	playerPosX = newPlayerPosX
	playerPosY = newPlayerPosY

	cameraPosX = playerPosX
	cameraPosY = playerPosY
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	var tileWidth = 16
	var tileHeight = 16
	visibleTilesX := ScreenWidth / tileWidth
	visibleTilesY := ScreenHeight / tileHeight

	// Calculate top-leftmost visible tile
	offsetX := cameraPosX - float64(visibleTilesX)/2.0
	offsetY := cameraPosY - float64(visibleTilesY)/2.0
	// Clamp camera close to the boundaries
	if offsetX < 0 {
		offsetX = 0
	}
	if offsetY < 0 {
		offsetY = 0
	}
	if offsetX > float64(g.levelMap.Width-visibleTilesX) {
		offsetX = float64(g.levelMap.Width - visibleTilesX)
	}
	if offsetY > float64(g.levelMap.Height-visibleTilesY) {
		offsetY = float64(g.levelMap.Height - visibleTilesY)
	}

	// Calculate tile offests for smooth movement (partial tiles to display)
	tileOffsetX := (offsetX - float64(int(offsetX))) * float64(tileWidth)
	tileOffsetY := (offsetY - float64(int(offsetY))) * float64(tileHeight)
	screen.Fill(lightBlue)
	// img := ebiten.NewImage(tileWidth, tileHeight)
	// Draw one tile more from left and right to avoid weird glitches on the edges
	for x := -1; x < visibleTilesX+1; x++ {
		for y := -1; y < visibleTilesY+1; y++ {
			tileIndex := g.levelMap.GetIndex(x+int(offsetX), y+int(offsetY))
			sx := tileIndex % 7
			sy := tileIndex / 7
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*tileWidth)-tileOffsetX, float64(y*tileHeight)-tileOffsetY)
			img := g.levelMap.Sprite.SubImage(image.Rect(spriteWidth*sx, spriteHeight*sy, spriteWidth*(sx+1), spriteHeight*(sy+1)))
			screen.DrawImage(ebiten.NewImageFromImage(img), op)
		}
	}
	// Draw player
	op := &ebiten.DrawImageOptions{}
	if !playerFacingRight {
		// flip horizontally
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(tileWidth), 0)
	}
	op.GeoM.Translate(float64(playerPosX-offsetX)*float64(tileWidth), float64(playerPosY-offsetY)*float64(tileHeight))
	// img.Fill(green)
	// img.Bounds()
	// screen.DrawImage(img, op)

	screen.DrawImage(ebiten.NewImageFromImage(playerCurrentFrame), op)

	DrawText(screen, "Hello world! 2022", 30, 30)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

var charactersAtlas *ebiten.Image
var (
	characterStanding image.Image
	characterRunning  []image.Image
)

func init() {
	fileContent, err := ioutil.ReadFile("./assets/files/pheasant.png")
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(bytes.NewReader(fileContent))
	if err != nil {
		log.Fatal("failed to decode:", err)
	}
	charactersAtlas = ebiten.NewImageFromImage(img)
	characterStanding = charactersAtlas.SubImage(image.Rect(0, 0, spriteWidth, spriteHeight))
	characterRunning1 := charactersAtlas.SubImage(image.Rect(spriteWidth*1, 0, spriteWidth*2, spriteHeight))
	characterRunning2 := charactersAtlas.SubImage(image.Rect(spriteWidth*2, 0, spriteWidth*3, spriteHeight))
	characterRunning3 := charactersAtlas.SubImage(image.Rect(spriteWidth*3, 0, spriteWidth*4, spriteHeight))
	characterRunning4 := charactersAtlas.SubImage(image.Rect(spriteWidth*4, 0, spriteWidth*5, spriteHeight))
	characterRunning = []image.Image{
		characterRunning1,
		characterRunning2,
		characterRunning3,
		characterRunning4,
	}
}

func main() {
	game := &Game{
		levelMap: level_map.NewMapVillage1(),
	}
	// Specify the window size as you like. Here, a doubled size is specified.
	ebiten.SetWindowSize(2*640, 2*480)
	ebiten.SetWindowTitle("Your game's title")
	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
