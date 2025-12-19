package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	BaseStart         = 0.5
	BaseProgressGain  = 0.01
	BaseProgressLoose = 0.01
)

type Progress struct {
	Value      float32
	StartLevel int
	GainLevel  int
	LooseLevel int
}

func NewProgress() *Progress {
	return &Progress{
		Value:      BaseStart,
		StartLevel: 1,
		GainLevel:  1,
		LooseLevel: 1,
	}
}

func (p *Progress) UpgradeStart() {
	p.StartLevel += 1
}

func (p *Progress) UpgradeGain() {
	p.GainLevel += 1
}

func (p *Progress) UpgradeLoose() {
	p.LooseLevel += 1
}

func (p *Progress) Reset() {
	p.Value = p.GetStart()
}

func (p *Progress) GetGain() float32 {
	return BaseProgressGain * float32(math.Pow(1.05, float64(p.GainLevel)))
}

func (p *Progress) GetLoose() float32 {
	return BaseProgressLoose * float32(math.Pow(0.95, float64(p.LooseLevel)))
}

func (p *Progress) GetStart() float32 {
	start := BaseStart * float32(math.Pow(1.05, float64(p.StartLevel)))
	if start > 1 {
		start = 1
	}
	return start
}

func (p *Progress) Increase() {
	p.Value += p.GetGain()
	if p.Value >= 1 {
		p.Value = 1
	}
}

func (p *Progress) Decrease() {
	p.Value -= p.GetLoose()
	if p.Value < 0 {
		p.Value = 0
	}
}

func (p *Progress) IsMax() bool {
	return p.Value >= 1.0
}

func (p *Progress) IsMin() bool {
	return p.Value <= 0
}

func (p *Progress) Draw(screen *ebiten.Image, g *Game) {
	vector.StrokeRect(
		screen,
		padding+barW+gap-borderWidth/2,
		padding-borderWidth/2,
		barW+borderWidth,
		frameH+borderWidth,
		borderWidth,
		borderColor,
		false,
	)

	progressH := frameH * p.Value

	vector.FillRect(
		screen,
		padding+barW+gap,
		padding+frameH-progressH,
		barW,
		progressH,
		g.Cursor.GetColor(),
		false,
	)
}
