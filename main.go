package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ScreenHeight  = 36
	TopPadding    = 3 // make room for score
	GameWidth     = 28
	GameHeight    = 31
	Zoom          = 4
	TileSize      = 8
	TileZoom      = TileSize * Zoom
	PlayerSpeed   = 60
	SlowSpeed     = 45   // slow down about 15-25% after eating dots
	SlowTime      = 0.25 // slow down for this long
	TeleportSpeed = 10   // slow down about 15-25% after eating dots
	TeleportTime  = 0.75 // slow down for this long
)

type Shape int

const (
	ShapeUp Shape = iota
	ShapeRight
	ShapeDown
	ShapeLeft
)

func (d Shape) Offset() int {
	switch d {
	case ShapeUp:
		return 2
	case ShapeRight:
		return 0
	case ShapeDown:
		return 3
	case ShapeLeft:
		return 1
	default:
		panic("unhandled default case")
	}
}

type Vec2i struct {
	x, y int
}

type Player struct {
	id            int
	name          string
	loc           map[Shape]Vec2i // location in spritesheet
	shape         Shape
	nextShape     Shape
	dir           rl.Vector2
	nextDir       rl.Vector2
	tileX         int
	tileY         int
	pixelX        float32 // location in maze
	pixelY        float32 // location in maze
	width         float32
	height        float32
	frameTime     float32
	frameSpeed    float32
	numFrames     int
	frame         int
	speed         float32
	slowTimer     float32
	teleportTimer float32
	score         int
	dots          int
	// spriteType bool // need better name
}

type Game struct {
	font      rl.Texture2D
	texture   rl.Texture2D
	image     *rl.Image
	player    *Player
	shader    rl.Shader
	debug     bool
	boardNum  int
	board     [][]Tile
	camera2   rl.Camera2D
	paused    bool
	highScore int
}

func main() {
	rl.SetTraceLogLevel(rl.LogWarning)

	rl.SetConfigFlags(rl.FlagVsyncHint)

	rl.InitWindow(GameWidth*TileSize*Zoom, ScreenHeight*TileSize*Zoom, "Ms. Packer Fan")
	defer rl.CloseWindow()

	//rl.SetTargetFPS(60)

	font := rl.LoadTexture("font.png")
	defer rl.UnloadTexture(font)
	texture := rl.LoadTexture("frozen_tundra.png")
	defer rl.UnloadTexture(texture)
	image := rl.LoadImageFromTexture(texture)
	defer rl.UnloadImage(image)

	g := initGame(font, texture, image)

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()
		g.Update(dt)

		rl.BeginDrawing()

		g.Draw()

		rl.EndDrawing()
	}
}

// newGame := func() *Game {}
func initGame(font, texture rl.Texture2D, image *rl.Image) *Game {
	g := &Game{}
	g.font = font
	g.texture = texture
	g.image = image
	g.shader = chromaShader()
	g.board = make([][]Tile, GameHeight)
	g.highScore = 1640
	for i := 0; i < GameHeight; i++ {
		g.board[i] = make([]Tile, GameWidth)
	}

	g.camera2 = rl.Camera2D{
		Offset:   rl.Vector2{0, TopPadding * TileSize * Zoom},
		Target:   rl.Vector2{0, 0},
		Rotation: 0,
		Zoom:     1,
	}

	startX := 13
	startY := 23

	g.player = &Player{
		id:   1,
		name: "ms. packer",
		loc: map[Shape]Vec2i{
			ShapeUp:    {456, 32},
			ShapeRight: {456, 0},
			ShapeDown:  {456, 48},
			ShapeLeft:  {456, 16},
		},
		pixelX:     float32(startX * TileSize * Zoom),
		pixelY:     float32(startY * TileSize * Zoom),
		width:      16,
		height:     16,
		tileX:      startX * Zoom,
		tileY:      startY * Zoom,
		shape:      ShapeLeft,
		nextShape:  ShapeLeft,
		dir:        rl.Vector2{-1, 0},
		nextDir:    rl.Vector2{-1, 0},
		frameTime:  0.0,
		frameSpeed: 0.1,
		numFrames:  3,
		frame:      0,
		speed:      PlayerSpeed * Zoom, // pixels per second
	}
	g.mapBoard()
	for y := 0; y < GameHeight; y++ {
		for x := 0; x < GameWidth; x++ {
			tile := g.board[y][x]
			if tile == Dot || tile == Power {
				g.player.dots++
			}
		}
	}
	//g.printBoard(true)

	return g
}
