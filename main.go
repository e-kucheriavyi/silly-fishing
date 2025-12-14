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

var bg = color.RGBA{48, 98, 48, 255}
var borderColor = color.RGBA{15, 56, 15, 255}
var limeColor = color.RGBA{155, 188, 15, 100}
var paleColor = color.RGBA{155, 188, 15, 255}

type GameStage byte

const (
	INTRO GameStage = iota
	GAME
	SCORE
)

type Game struct {
	Cursor float32
	Fish float32
	Skill float32
	Score int
	Record int
	Progress float32
	IsColliding bool
	Pressing bool
	D float32
	IsCatching bool
	Stage GameStage
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

	if g.Stage == INTRO || g.Stage == SCORE {
		if g.IsCatching && g.Pressing {
			return nil
		}

		g.IsCatching = false

		if g.Pressing {
			g.Score = 0
			g.Stage = GAME
		}

		return nil
	}

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
		g.IsCatching = true
		g.Progress += 0.01

		if g.Progress >= 1.0 {
			g.Progress = 0
			g.Skill -= g.Skill * 0.05
			g.Score += 1
			if g.Score > g.Record {
				g.Record = g.Score
			}
			g.Fish = rand.Float32()
			g.IsCatching = false
		}

		return nil
	}

	g.Progress -= 0.01

	if g.Progress <= 0.0 {
		g.Progress = 0.0

		if g.IsCatching {
			g.Stage = SCORE
		}
	}

	return nil
}

func (g *Game) DrawIntro(screen *ebiten.Image) {
	ebitenutil.DebugPrint(
		screen,
		`Rules:
		- Press Space or Left Mouse Button to move cursor higher
		- The cursor drops when nothing is pressed
		- Keep the target inside your cursor to fill the progress bar
		- When the bar is filled you get one point
		- With each point cursor shinks by 5%
		- You lose if the bar drops to zero

		Press Space or Left Mouse Button to start`,
	)

}

func (g *Game) DrawGame(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %d", g.Score))

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

	cursorColor := paleColor

	if g.IsColliding {
		cursorColor = limeColor
	}

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

func (g *Game) DrawScore(screen *ebiten.Image) {
	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf("Game Over\nScore: %d\nRecord: %d", g.Score, g.Record),
	)
}

func (g *Game) Draw(screen *ebiten.Image) {
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

	if g.Stage == INTRO {
		g.DrawIntro(screen)
		return
	}

	if g.Stage == SCORE {
		g.DrawScore(screen)
		return
	}

	g.DrawGame(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenW, screenH
}

func main() {
	game := Game{
		Cursor: 0.65,
		Skill: 0.7,
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


