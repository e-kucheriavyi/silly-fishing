package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	BaseCursorSize   = 0.7
	BaseReduceAmount = 0.05
	BaseMinSize      = 0.1
)

type Cursor struct {
	X                 float32
	IsColliding       bool
	Color             color.Color
	ActiveColor       color.Color
	MinSizeLevel      int
	ReduceAmountLevel int
	IsMinSize         bool
	Level             int
}

func NewCursor() *Cursor {
	return &Cursor{
		X:                 0.3,
		MinSizeLevel:      1,
		ReduceAmountLevel: 1,
	}
}

func (c *Cursor) UpgradeMinSize() {
	c.MinSizeLevel += 1
}

func (c *Cursor) UpgradeReduceAmount() {
	c.ReduceAmountLevel += 1
}

func (c *Cursor) Collide(t *Target) bool {
	tStart := t.X
	tEnd := t.X + t.GetSize()

	cStart := c.X
	cEnd := c.X + c.GetSize()

	c.IsColliding = tStart >= cStart && tEnd <= cEnd

	return c.IsColliding
}

func (c *Cursor) GetColor() color.Color {
	if c.IsColliding {
		return c.ActiveColor
	}
	return c.Color
}

func (c *Cursor) GetSize() float32 {
	minSize := c.GetMinSize()

	if c.IsMinSize {
		return minSize
	}

	r := math.Pow(1-BaseReduceAmount, float64(c.ReduceAmountLevel))

	size := float32(BaseCursorSize * math.Pow(r, float64(c.Level)))

	if size > minSize {
		return size
	}

	c.IsMinSize = true

	return minSize
}

func (c *Cursor) GetMinSize() float32 {
	return BaseMinSize
}

func (c *Cursor) Move(g *Game) {
	if g.IsPressing {
		c.X -= 0.01
	} else {
		c.X += 0.01
	}

	size := c.GetSize()

	if c.X > 1.0-size {
		c.X = 1.0 - size
	} else if c.X < 0 {
		c.X = 0.0
	}
}

func (c *Cursor) Draw(screen *ebiten.Image, g *Game) {
	vector.FillRect(
		screen,
		padding,
		padding+frameH*c.X,
		barW,
		frameH*c.GetSize(),
		c.GetColor(),
		false,
	)
}
