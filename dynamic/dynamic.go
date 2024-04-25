package dynamic

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type DynamicObj struct {
	PX, PY     float64
	VX, VY     float64
	solidVsMap bool
	solidVsDyn bool
	friendly   bool
	name       string
}

func NewDynamicObj(name string) *DynamicObj {
	return &DynamicObj{
		name: name,
		PX:   0,
		PY:   0,
		VX:   0,
		VY:   0,
		// upfront assumption
		solidVsMap: true,
		solidVsDyn: true,
		friendly:   true,
	}
}

type Dynamic interface {
	Draw(screen *ebiten.Image, ox, oy float64)
	Update()
	GetPX() float64
	SetPX(float64)
	GetPY() float64
	SetPY(float64)
	GetVX() float64
	SetVX(float64)
	GetVY() float64
	SetVY(float64)
}

func NewDynamicCreate(name string, sprite *ebiten.Image) *DynamicCreature {
	return &DynamicCreature{
		DynamicObj:      *NewDynamicObj(name),
		sprite:          sprite,
		Health:          10,
		HealthMax:       10,
		FacingDirection: South,
		GraphicState:    Standing,
		graphicCounter:  0,
		time:            0,
	}
}

type Direction int

const (
	North Direction = iota
	East
	South
	West
)

type State int

const (
	Standing State = iota
	Walking
	Celebrating
	Dead
)

type DynamicCreature struct {
	DynamicObj
	sprite          *ebiten.Image
	Health          int
	HealthMax       int
	FacingDirection Direction
	GraphicState    State
	time            int
	graphicCounter  int
}

func (c *DynamicCreature) Draw(screen *ebiten.Image, ox, oy float64) {
	sheetOffsetX := 0
	sheetOffsetY := 0

	switch c.GraphicState {
	case Standing:
		sheetOffsetX = int(c.FacingDirection) * 16
	case Walking:
		sheetOffsetX = int(c.FacingDirection) * 16
		sheetOffsetY = c.graphicCounter * 16
	case Dead:
		sheetOffsetX = 4 * 16
		sheetOffsetY = 1 * 16

	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(c.PX-ox, c.PY-oy)
	img := c.sprite.SubImage(image.Rect(sheetOffsetX, sheetOffsetY, sheetOffsetX+16, sheetOffsetY+16))
	screen.DrawImage(ebiten.NewImageFromImage(img), op)
}

func (c *DynamicCreature) Update() {
	c.time++
	if c.time > 10 {
		c.time -= 10
		c.graphicCounter++
		c.graphicCounter %= 2
	}

	if math.Abs(c.VX) > 0 || math.Abs(c.VY) > 0 {
		c.GraphicState = Walking
	} else {
		c.GraphicState = Standing
	}

	if c.VX < -0.1 {
		c.FacingDirection = West
	}
	if c.VX > 0.1 {
		c.FacingDirection = East
	}
	if c.VY < -0.1 {
		c.FacingDirection = North
	}
	if c.VY > 0.1 {
		c.FacingDirection = South
	}

	if c.Health <= 0 {
		c.GraphicState = Dead
	}
}

func (c *DynamicCreature) SetPX(px float64) { c.PX = px }
func (c *DynamicCreature) SetPY(py float64) { c.PY = py }
func (c *DynamicCreature) GetPX() float64   { return c.PX }
func (c *DynamicCreature) GetPY() float64   { return c.PY }
func (c *DynamicCreature) SetVX(vx float64) { c.VX = vx }
func (c *DynamicCreature) SetVY(vy float64) { c.VY = vy }
func (c *DynamicCreature) GetVX() float64   { return c.VX }
func (c *DynamicCreature) GetVY() float64   { return c.VY }
