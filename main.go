package main

import (
	"log"
	"fmt"
	"image/color"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenW = 640
	screenH = 480
	screenPadding = 32
	framePadding = 32
	gap = 32
	frameW = screenW - screenPadding * 2
	innerFrameW = frameW - framePadding * 2
	frameH = screenH - screenPadding * 2
	innerFrameH = frameH - framePadding * 2
	barW = (screenW - screenPadding * 2 - framePadding * 2 - gap) / 2
	fishScale = 0.05
	fishW = innerFrameH * fishScale
	fishSpeed = 0.005
)

type Game struct {
	Cursor float32
	Fish float32
	Skill float32
	Score float32
	Progress float32
	IsColliding bool
	Pressing bool
	D float32
}

func (g *Game) Collide() bool {
	fishStart := g.Fish * innerFrameH
	fishEnd := fishStart + fishW

	cursorStart := g.Cursor * innerFrameH
	cursorEnd := cursorStart + innerFrameH * g.Skill

	return fishStart >= cursorStart && fishEnd <= cursorEnd
}

func (g *Game) Update() error {
	g.Pressing = ebiten.IsKeyPressed(ebiten.KeySpace)

	if rand.Float32() > 0.8 {
		g.D *= -1.0
	}

	g.Fish = g.Fish + fishSpeed * g.D

	if g.Fish > 1.0 - fishScale {
		g.Fish = 1.0 - fishScale
	}

	if g.Fish < 0 {
		g.Fish = 0.0
	}

	if g.Pressing {
		g.Cursor -= 0.01
	} else {
		g.Cursor += 0.01
	}

	if g.Cursor > 1.0 - g.Skill {
		g.Cursor = 1.0 - g.Skill
	}
	if g.Cursor < 0 {
		g.Cursor = 0.0
	}

	g.IsColliding = g.Collide()

	if g.IsColliding {
		g.Progress += 0.01

		if g.Progress >= 1.0 {
			g.Progress = 0
			g.Skill -= g.Skill * 0.05
			g.Score += 100
			g.Fish = rand.Float32()
		}

		return nil
	}

	g.Progress -= 0.01

	if g.Progress <= 0.0 {
		g.Progress = 0.0
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %.2f", g.Score))

	// frame
	vector.FillRect(
		screen,
		screenPadding,
		screenPadding,
		frameW,
		frameH,
		color.RGBA{100, 100, 100, 255},
		false,
	)

	// fish bar
	vector.FillRect(
		screen,
		screenPadding + framePadding,
		screenPadding + framePadding,
		barW,
		innerFrameH,
		color.RGBA{20, 20, 20, 255},
		false,
	)

	var cursorA uint8 = 100

	if g.IsColliding {
		cursorA = 255
	}

	// cursor
	vector.FillRect(
		screen,
		screenPadding + framePadding,
		screenPadding + framePadding + innerFrameH * g.Cursor,
		barW,
		innerFrameH * g.Skill,
		color.RGBA{0, cursorA, 0, cursorA},
		false,
	)

	// fish
	vector.FillRect(
		screen,
		screenPadding + framePadding + barW / 2 - fishW / 2,
		screenPadding + framePadding + innerFrameH * g.Fish,
		fishW,
		fishW,
		color.RGBA{255, 0, 0, 255},
		false,
	)

	// progress bar
	vector.FillRect(
		screen,
		screenPadding + framePadding + barW + gap,
		screenPadding + framePadding,
		barW,
		innerFrameH,
		color.RGBA{20, 20, 20, 255},
		false,
	)

	progressH := innerFrameH * g.Progress

	// progress filling
	vector.FillRect(
		screen,
		screenPadding + framePadding + barW + gap,
		screenPadding + framePadding + innerFrameH - progressH,
		barW,
		progressH,
		color.RGBA{0, 0, 255, 255},
		false,
	)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenW, screenH
}

func main() {
	game := Game{
		Cursor: 0.65,
		Skill: 0.7,
		Score: 0.0,
		Progress: 0.7,
		Fish: 0.7,
		D: 1.0,
	}

	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Silly Fishing")
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}


