package main

import (
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Tile byte

type Maze [][]Tile

const (
	Wall Tile = iota
	Dot
	Power
	Empty
	Tunnel

	DotMask   = 103481868288
	PowerMask = 4359202964317896252
	//DoorMask  = 16776960
)

func (t Tile) String() string {
	switch t {
	case Wall:
		return "X"
	case Dot:
		return "."
	case Power:
		return "*"
	case Empty:
		return " "
	case Tunnel:
		return "@"
	default:
		panic("unhandled default case")
	}
}

func (t Tile) Name() string {
	switch t {
	case Wall:
		return "wall"
	case Dot:
		return "dot"
	case Power:
		return "power"
	case Empty:
		return "empty"
	case Tunnel:
		return "tunnel"
	default:
		panic("unhandled default case")
	}
}

func (t Tile) Pretty() string {
	switch t {
	case Wall:
		return "XXX"
	case Dot:
		return " + "
	case Power:
		return "(*)"
	case Empty:
		return "   "
	case Tunnel:
		return "<@>"
	default:
		panic("unhandled default case")
	}
}

func (m Maze) String() string {
	sb := strings.Builder{}
	for y := 0; y < len(m); y++ {
		for x := 0; x < len(m[y]); x++ {
			sb.WriteString(m[y][x].String())
		}
		sb.WriteString("\n")
	}

	return sb.String()

}

func (m Maze) Pretty() string {
	sb := strings.Builder{}
	for y := 0; y < len(m); y++ {
		for x := 0; x < len(m[y]); x++ {
			sb.WriteString(m[y][x].Pretty())
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (m Maze) IsValidMove(tile Vec2i) bool {
	if tile.X < 0 || tile.X >= len(m[0]) || tile.Y < 0 || tile.Y >= len(m) {
		return false
	}
	return m[tile.Y][tile.X] != Wall
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

	offset := g.boardNum * GameHeight * Size
	var piece Tile
	for y := 0; y < GameHeight; y++ {
		for x := 0; x < GameWidth; x++ {
			p, _ := g.readPixelArea(int32(x*Size), int32(y*Size+offset), Size, Size)
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
	g.tunnels = make([]Vec2i, 0, 4)
	for y := 0; y < GameHeight; y++ {
		if g.maze[y][0] == Empty &&
			(g.maze[y][1] == Empty || g.maze[y][1] == Dot) &&
			(g.maze[y][2] == Empty || g.maze[y][2] == Dot) {
			g.maze[y][0] = Tunnel
			g.tunnels = append(g.tunnels, Vec2i{X: 0, Y: y})
		}

		if (g.maze[y][GameWidth-3] == Empty || g.maze[y][GameWidth-3] == Dot) &&
			(g.maze[y][GameWidth-2] == Empty || g.maze[y][GameWidth-2] == Dot) &&
			g.maze[y][GameWidth-1] == Empty {
			g.maze[y][GameWidth-1] = Tunnel
			g.tunnels = append(g.tunnels, Vec2i{X: 0, Y: y})
		}
	}
}

func (g *Game) InTunnel(e *Entity) bool {
	x, y := e.tile.X, e.tile.Y
	for _, t := range g.tunnels {
		if y == t.Y && (x <= 0 || x >= GameWidth-1) {
			return true
		}
	}

	return false
}
