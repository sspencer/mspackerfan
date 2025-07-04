package main

import (
	"fmt"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

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

			g.board[y][x] = piece
		}
	}

	// find tunnels
	for y := 0; y < GameHeight; y++ {
		if g.board[y][0] == Empty &&
			(g.board[y][1] == Empty || g.board[y][1] == Dot) &&
			(g.board[y][2] == Empty || g.board[y][2] == Dot) {
			g.board[y][0] = Tunnel
		}

		if (g.board[y][GameWidth-3] == Empty || g.board[y][GameWidth-3] == Dot) &&
			(g.board[y][GameWidth-2] == Empty || g.board[y][GameWidth-2] == Dot) &&
			g.board[y][GameWidth-1] == Empty {
			g.board[y][GameWidth-1] = Tunnel
		}
	}
}

func (g *Game) printBoard(pretty bool) {
	fmt.Printf("--- Board %d ---\n", g.boardNum)
	for y := 0; y < GameHeight; y++ {
		for x := 0; x < GameWidth; x++ {
			fmt.Printf("%s", g.board[y][x].Pretty())
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")
}
