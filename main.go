package main

import (
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ScreenHeight = 36
	TopPadding   = 3 // make room for score
	GameWidth    = 28
	GameHeight   = 31
	Zoom         = 4
	TileSize     = 8
	ChaseBug     = true // error in chase state in original game
)

type Game struct {
	font        rl.Texture2D
	texture     rl.Texture2D
	image       *rl.Image
	player      *Packer
	ghosts      []*Ghost
	shader      rl.Shader
	boardNum    int
	maze        Maze
	camera2     rl.Camera2D
	paused      bool
	debug       bool
	debugLayout bool
	highScore   int
}

func main() {
	rl.SetTraceLogLevel(rl.LogWarning)

	rl.InitWindow(GameWidth*TileSize*Zoom, ScreenHeight*TileSize*Zoom, "Ms. Packer Fan")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	font := rl.LoadTexture("font.png")
	defer rl.UnloadTexture(font)
	texture := rl.LoadTexture("frozen_tundra.png")
	defer rl.UnloadTexture(texture)
	image := rl.LoadImageFromTexture(texture)
	defer rl.UnloadImage(image)

	g := initGame(font, texture, image)

	for !rl.WindowShouldClose() {
		g.Update()

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
	g.maze = make(Maze, GameHeight)
	g.highScore = 0

	for i := 0; i < GameHeight; i++ {
		g.maze[i] = make([]Tile, GameWidth)
	}

	g.camera2 = rl.Camera2D{
		Offset:   rl.Vector2{Y: TopPadding * TileSize * Zoom},
		Target:   rl.Vector2{},
		Rotation: 0,
		Zoom:     1,
	}

	dots := 0
	g.mapBoard()
	for y := 0; y < GameHeight; y++ {
		for x := 0; x < GameWidth; x++ {
			tile := g.maze[y][x]
			if tile == Dot {
				dots++
			}
		}
	}

	// g.maze.String()

	g.player = NewPacker(dots)
	g.ghosts = make([]*Ghost, 0, 1)
	//g.ghosts = append(g.ghosts, NewGhost(Blinky{})) // do not reorder ghosts
	g.ghosts = append(g.ghosts, NewGhost(Pinky{}))
	//g.ghosts = append(g.ghosts, NewGhost(Inky{}))
	//g.ghosts = append(g.ghosts, NewGhost(Clyde{}))

	return g
}

// always get starting "dot" color at boardNum pos 1,1
// then go row by row to get either dot (.), power (o), border (pixelX) or empty (<sp>)

func (g *Game) readPixelArea(startX, startY, width, height int32) (uint64, string) {

	var result uint64
	tile := strings.Builder{}

	// Read each pixel in the specified area
	for y := startY; y < startY+height; y++ {
		for x := startX; x < startX+width; x++ {
			// Get the color of the pixel at position (pixelX, pixelY)
			result <<= 1

			color := rl.GetImageColor(*g.image, x, y)
			if color.R > 0 || color.G > 0 || color.B > 0 {
				result |= 1
				tile.WriteString("1")
			} else {
				tile.WriteString("0")
			}
		}
		tile.WriteString("\n")
	}

	return result, tile.String()
}

func (g *Game) mapBoard() {

	offset := g.boardNum * GameHeight * TileSize
	var piece Tile
	for y := 0; y < GameHeight; y++ {
		for x := 0; x < GameWidth; x++ {
			p, _ := g.readPixelArea(int32(x*TileSize), int32(y*TileSize+offset), TileSize, TileSize)
			if p == 0 {
				piece = Empty
			} else if p == DotMask { // dot
				piece = Dot
			} else if p == PowerMask {
				piece = Power
			} else {
				piece = Wall
			}

			g.maze[y][x] = piece
		}
	}

	// find tunnels
	for y := 0; y < GameHeight; y++ {
		if g.maze[y][0] == Empty &&
			(g.maze[y][1] == Empty || g.maze[y][1] == Dot) &&
			(g.maze[y][2] == Empty || g.maze[y][2] == Dot) {
			g.maze[y][0] = Tunnel
		}

		if (g.maze[y][GameWidth-3] == Empty || g.maze[y][GameWidth-3] == Dot) &&
			(g.maze[y][GameWidth-2] == Empty || g.maze[y][GameWidth-2] == Dot) &&
			g.maze[y][GameWidth-1] == Empty {
			g.maze[y][GameWidth-1] = Tunnel
		}
	}
}
