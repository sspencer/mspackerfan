package main

import (
	"strings"
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
