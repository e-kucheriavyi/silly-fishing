package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	BaseTargetSpeed        = 0.005
	BaseMaxSpeed           = 0.05
	BaseAccelerationAmount = 1.0005
)

type Target struct {
	X                 float32
	D                 float32
	MoveD             float32
	Level             int
	AccelerationLevel int
	MaxSpeedLevel     int
}

func NewTarget() *Target {
	return &Target{
		X:                 0.35,
		D:                 1,
		AccelerationLevel: 1,
		MaxSpeedLevel:     1,
	}
}

func (t *Target) UpgradeMaxSpeed() {
	t.MaxSpeedLevel += 1
}

func (t *Target) UpgradeAcceleration() {
	t.AccelerationLevel += 1
}

func (t *Target) GetSize() float32 {
	return BaseMinSize * 0.7
}

func (t *Target) GetSpeed() float32 {
	a := math.Pow(BaseAccelerationAmount, float64(t.Level))
	speed := BaseTargetSpeed * float32(a)

	maxSpeed := t.GetMaxSpeed()

	if speed < maxSpeed {
		return speed
	}

	return maxSpeed
}

func (t *Target) GetMaxSpeed() float32 {
	a := math.Pow(0.95, float64(t.MaxSpeedLevel))
	return BaseMaxSpeed * float32(a)
}

func (t *Target) Vobble() {
	if t.MoveD <= 0 {
		t.D *= -1.0
		t.MoveD = rand.Float32()
		return
	}

	speed := t.GetSpeed()
	size := t.GetSize()

	t.X = t.X + speed*t.D
	t.MoveD -= speed

	if t.X > 1.0-size {
		t.X = 1.0 - size
		t.MoveD = 0.0
	}

	if t.X < 0 {
		t.X = 0.0
		t.MoveD = 0.0
	}
}

func (t *Target) Draw(screen *ebiten.Image, g *Game) {
	size := t.GetSize()
	w := frameH * size

	vector.FillRect(
		screen,
		padding+barW/2-w/2,
		padding+frameH*t.X,
		w,
		w,
		color.RGBA{255, 0, 0, 255},
		false,
	)

	vector.StrokeRect(
		screen,
		padding+barW/2-w/2,
		padding+frameH*t.X,
		w,
		w,
		borderWidth,
		borderColor,
		false,
	)
}
