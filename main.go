package main

import (
	"flag"
	"fmt"
	"runtime/debug"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ScreenHeight   = 36
	TopPadding     = 3 // make room for score
	GameWidth      = 28
	GameHeight     = 31
	Zoom           = 4
	Size           = 8
	Pixel          = Size * Zoom
	ChaseBug       = true // error in chase state in original game
	FrightDuration = 6.0
)

type Game struct {
	font     rl.Texture2D
	texture  rl.Texture2D
	image    *rl.Image
	player   *Player
	ghosts   []*Ghost
	shader   rl.Shader
	boardNum int
	level    int
	maze     Maze
	//maze        [31][28]Tile
	tunnels     []Vec2i
	camera2     rl.Camera2D
	paused      bool
	debug       bool
	moved       bool
	debugLayout bool
	highScore   int
	startTime   float64
	levelTime   float64
	frightTime  float64
	dotsEaten   int
}

func main() {
	debugMode := false
	flag.BoolVar(&debugMode, "d", false, "enable debug mode")
	flag.Parse()

	// less GC
	debug.SetGCPercent(200)

	rl.SetTraceLogLevel(rl.LogWarning)

	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(GameWidth*Pixel, ScreenHeight*Pixel, "Ms. Player Fan")
	defer rl.CloseWindow()

	//rl.SetTargetFPS(60)

	font := rl.LoadTexture("font.png")
	defer rl.UnloadTexture(font)
	texture := rl.LoadTexture("frozen_tundra.png")
	defer rl.UnloadTexture(texture)
	image := rl.LoadImageFromTexture(texture)
	defer rl.UnloadImage(image)

	g := initGame(font, texture, image, debugMode)

	for !rl.WindowShouldClose() {
		g.Update()

		rl.BeginDrawing()

		g.Draw()

		rl.EndDrawing()
	}
}

// newGame := func() *Game {}
func initGame(font, texture rl.Texture2D, image *rl.Image, debugMode bool) *Game {
	g := &Game{}
	g.font = font
	g.texture = texture
	g.image = image
	g.shader = chromaShader()
	g.maze = make(Maze, GameHeight)
	g.highScore = 0
	g.startTime = rl.GetTime()
	g.level = 1
	g.debug = debugMode
	for i := 0; i < GameHeight; i++ {
		g.maze[i] = make([]Tile, GameWidth)
	}

	g.camera2 = rl.Camera2D{
		Offset:   rl.Vector2{Y: TopPadding * Pixel},
		Target:   rl.Vector2{},
		Rotation: 0,
		Zoom:     1,
	}

	g.mapBoard()
	for i, t := range g.tunnels {
		fmt.Printf("Tunnel %d: %v\n", i, t)
	}

	//dots := 0
	//for y := 0; y < GameHeight; y++ {
	//	for x := 0; x < GameWidth; x++ {
	//		tile := g.maze[y][x]
	//		if tile == Dot {
	//			dots++
	//		}
	//	}
	//}

	// g.maze.String()

	g.player = NewPlayer()
	//if g.debug {
	//	g.ghosts = make([]*Ghost, 1)
	//	g.ghosts[0] = NewGhost(g, Blinky{})
	//
	//} else {
	g.ghosts = make([]*Ghost, 4)
	g.ghosts[0] = NewGhost(g, Blinky{})
	g.ghosts[1] = NewGhost(g, Pinky{})
	g.ghosts[2] = NewGhost(g, Inky{})
	g.ghosts[3] = NewGhost(g, Clyde{})
	//}
	return g
}
