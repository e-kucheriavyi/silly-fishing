package main

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenW = 640
	screenH = 480

	padding = 32
	gap     = 32

	frameW = screenW - padding*2
	frameH = screenH - padding*2

	barW = (screenW - padding*2 - gap) / 2

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
	UPGRADES
	SHOP
	SECRET_ANIMATION
)

type Game struct {
	Cursor   *Cursor
	Target   *Target
	Progress *Progress

	Score  int
	Record int
	Total  int

	Prestige int

	IsPressing bool

	Stage   GameStage
	EndTime time.Time
}

func NewGame() *Game {
	target := NewTarget()
	cursor := NewCursor()
	progress := NewProgress()

	cursor.Color = paleColor
	cursor.ActiveColor = limeColor

	return &Game{
		Target:   target,
		Cursor:   cursor,
		Progress: progress,
	}
}

func (g *Game) GetScoreGain() int {
	s := 2 * g.Prestige

	if s == 0 {
		return 1
	}

	return s
}

func (g *Game) Start() {
	g.Score = 0
	g.Cursor.Level = 1
	g.Target.Level = 1

	g.Progress.Reset()

	g.Stage = GAME
}

func (g *Game) Catch() {
	g.Progress.Reset()

	if g.Cursor.IsMinSize {
		g.Target.Level += 1
	} else {
		g.Cursor.Level += 1
	}

	g.Score += g.GetScoreGain()
	g.Total += g.GetScoreGain()

	if g.Score > g.Record {
		g.Record = g.Score
	}

	g.Cursor.X = 0.3
	g.Target.X = 0.35
}

func (g *Game) Finish() {
	g.Stage = SCORE
	g.EndTime = time.Now()
}

func (g *Game) ListenPressing() {
	isMousePressed := ebiten.IsMouseButtonPressed(ebiten.MouseButton0)
	isSpacePressed := ebiten.IsKeyPressed(ebiten.KeySpace)
	// TODO: update to non deprecated version
	isTouching := len(ebiten.TouchIDs()) > 0

	g.IsPressing = isSpacePressed || isMousePressed || isTouching
}

func (g *Game) OpenUpgrades() {
	g.Stage = UPGRADES
}

func (g *Game) OpenShop() {
	g.Stage = SHOP
}

func (g *Game) Update() error {
	g.ListenPressing()

	if g.Stage == INTRO && g.IsPressing {
		g.Start()
		return nil
	}

	if g.Stage == SCORE && g.IsPressing {
		if time.Since(g.EndTime) < 1*time.Second {
			return nil
		}

		g.Start()
		return nil
	}

	g.Target.Vobble()

	g.Cursor.Move(g)

	g.Cursor.Collide(g.Target)

	if g.Cursor.IsColliding {
		g.Progress.Increase()

		if g.Progress.IsMax() {
			g.Catch()
		}

		return nil
	}

	g.Progress.Decrease()

	if g.Progress.IsMin() {
		g.Finish()
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

	vector.StrokeRect(
		screen,
		padding-borderWidth/2,
		padding-borderWidth/2,
		barW+borderWidth,
		frameH+borderWidth,
		borderWidth,
		borderColor,
		false,
	)

	g.Cursor.Draw(screen, g)

	g.Target.Draw(screen, g)

	g.Progress.Draw(screen, g)
}

func (g *Game) DrawScore(screen *ebiten.Image) {
	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf("Game Over\nScore: %d\nRecord: %d", g.Score, g.Record),
	)
}

func (g *Game) DrawUpgrades(screen *ebiten.Image) {
	// TODO
	w := float32(frameW/3 - gap*2)
	h := float32(frameH/3 - gap*2)

	vector.FillRect(
		screen,
		padding,
		padding,
		w,
		h,
		limeColor,
		false,
	)
}

func (g *Game) DrawShop(screen *ebiten.Image) {
	// TODO
}

func (g *Game) Draw(screen *ebiten.Image) {
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

	if g.Stage == UPGRADES {
		g.DrawUpgrades(screen)
		return
	}

	if g.Stage == SHOP {
		g.DrawShop(screen)
		return
	}

	g.DrawGame(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenW, screenH
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Silly Fishing")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
