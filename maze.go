package main

import (
	"strings"
)

type Maze [][]Tile

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
