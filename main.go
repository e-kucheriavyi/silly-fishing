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

	padding = 32
	gap = 32

	frameW = screenW - padding * 2
	frameH = screenH - padding * 2

	barW = (screenW - padding * 2 - gap) / 2

	fishScale = 0.05
	fishW = frameH * fishScale
	fishSpeed = 0.005

	borderWidth = 8
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
	fishStart := g.Fish * frameH
	fishEnd := fishStart + fishW

	cursorStart := g.Cursor * frameH
	cursorEnd := cursorStart + frameH * g.Skill

	return fishStart >= cursorStart && fishEnd <= cursorEnd
}

func (g *Game) Update() error {
	isMousePressed := ebiten.IsMouseButtonPressed(ebiten.MouseButton0)
	isSpacePressed := ebiten.IsKeyPressed(ebiten.KeySpace)
	g.Pressing = isSpacePressed || isMousePressed

	if rand.Float32() > 0.9 {
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
	bg := color.RGBA{48, 98, 48, 255}
	borderColor := color.RGBA{15, 56, 15, 255}

	// bg
	vector.FillRect(
		screen,
		0,
		0,
		screenW,
		screenH,
		bg,
		false,
	)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %.2f", g.Score))

	// fish bar
	vector.StrokeRect(
		screen,
		padding - borderWidth / 2,
		padding - borderWidth / 2,
		barW + borderWidth,
		frameH + borderWidth,
		borderWidth,
		borderColor,
		false,
	)

	var cursorA uint8 = 255

	if g.IsColliding {
		cursorA = 100
	}

	cursorColor := color.RGBA{155, 188, 15, cursorA}

	// cursor
	vector.FillRect(
		screen,
		padding,
		padding + frameH * g.Cursor,
		barW,
		frameH * g.Skill,
		cursorColor,
		false,
	)

	// fish
	vector.FillRect(
		screen,
		padding + barW / 2 - fishW / 2,
		padding + frameH * g.Fish,
		fishW,
		fishW,
		color.RGBA{255, 0, 0, 255},
		false,
	)

	vector.StrokeRect(
		screen,
		padding + barW / 2 - fishW / 2,
		padding + frameH * g.Fish,
		fishW,
		fishW,
		borderWidth,
		borderColor,
		false,
	)

	// progress bar
	vector.StrokeRect(
		screen,
		padding + barW + gap - borderWidth / 2,
		padding - borderWidth / 2,
		barW + borderWidth,
		frameH + borderWidth,
		borderWidth,
		borderColor,
		false,
	)

	progressH := frameH * g.Progress

	// progress filling
	vector.FillRect(
		screen,
		padding + barW + gap,
		padding + frameH - progressH,
		barW,
		progressH,
		cursorColor,
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


