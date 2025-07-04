package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ScreenHeight = 36
	TopPadding   = 3 // make room for score
	GameWidth    = 28
	GameHeight   = 31
	Zoom         = 4
	TileSize     = 8
	TileZoom     = TileSize * Zoom
)

type Game struct {
	font      rl.Texture2D
	texture   rl.Texture2D
	image     *rl.Image
	player    *Entity
	ghosts    []*Entity
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

	dots := 0
	g.mapBoard()
	for y := 0; y < GameHeight; y++ {
		for x := 0; x < GameWidth; x++ {
			tile := g.board[y][x]
			if tile == Dot || tile == Power {
				dots++
			}
		}
	}

	//g.printBoard(true)

	g.player = createPlayer(dots)

	g.ghosts = make([]*Entity, 4)
	g.ghosts[0] = createGhost("Blinky", 9, 14, 64, ShapeUp)
	g.ghosts[1] = createGhost("Inky", 12, 14, 96, ShapeDown)
	g.ghosts[2] = createGhost("Pinky", 14, 14, 80, ShapeUp)
	g.ghosts[3] = createGhost("Clyde", 16, 14, 112, ShapeDown)

	return g
}
