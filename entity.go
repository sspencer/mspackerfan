package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Direction int

const (
	None Direction = iota
	Up
	Right
	Down
	Left
)

type Entity struct {
	name        string
	sprite      map[Direction]Vec2i // location in spritesheet
	dir         Direction
	nextDir     Direction
	vel         Vec2i
	nextVel     Vec2i
	tile        Vec2i
	pixel       rl.Vector2
	pixelsMoved float32
	width       float32
	height      float32
	frameCount  int
	frame       int
}

func (d Direction) String() string {
	switch d {
	case Up:
		return "Up"
	case Right:
		return "Right"
	case Down:
		return "Down"
	case Left:
		return "Left"
	default:
		panic("unhandled default case")
	}
}

func (d Direction) Opposite() Direction {
	switch d {
	case Up:
		return Down
	case Right:
		return Left
	case Down:
		return Up
	case Left:
		return Right
	default:
		panic("unhandled default case")
	}
}

func (d Direction) Vector() Vec2i {
	switch d {
	case Up:
		return Vec2i{Y: -1}
	case Right:
		return Vec2i{X: 1}
	case Down:
		return Vec2i{Y: 1}
	case Left:
		return Vec2i{X: -1}
	default:
		panic("unhandled default case")
	}
}

func (d Direction) GetNextTile(vec Vec2i) Vec2i {
	return Vec2i{
		X: vec.X + int(d.Vector().X),
		Y: vec.Y + int(d.Vector().Y),
	}
}

type Vec2i struct {
	X, Y int
}

var ZeroVec = Vec2i{X: 0, Y: 0}

func (v Vec2i) String() string {
	return fmt.Sprintf("(%d, %d)", v.X, v.Y)
}

func (v Vec2i) Add(x, y int) Vec2i {
	return Vec2i{
		X: v.X + x,
		Y: v.Y + y,
	}
}

func (v Vec2i) Distance(b Vec2i) float32 {
	dx := v.X - b.X
	dy := v.Y - b.Y
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}

func (v Vec2i) IsZero() bool {
	return v.X == 0 && v.Y == 0
}

func (v Vec2i) InMaze() bool {
	return v.X >= 0 && v.X < GameWidth && v.Y >= 0 && v.Y < GameHeight
}

func (v Vec2i) Clamp() Vec2i {
	if v.X < 0 {
		v.X = 0
	} else if v.X > GameWidth-1 {
		v.X = GameWidth - 1
	}

	if v.Y < 0 {
		v.Y = 0
	} else if v.Y > GameHeight-1 {
		v.Y = GameHeight - 1
	}

	return v
}

func (e *Entity) move(speed float32) {
	if e.vel.X != 0 || e.vel.Y != 0 {
		e.pixelsMoved += speed

		clampedPixelsMoved := float32(math.Min(float64(e.pixelsMoved), float64(TileSize)))
		visualOffsetX := float32(e.vel.X) * clampedPixelsMoved
		visualOffsetY := float32(e.vel.Y) * clampedPixelsMoved

		e.pixel.X = (float32(e.tile.X*TileSize) + visualOffsetX - TileSize/2) * Zoom
		e.pixel.Y = (float32(e.tile.Y*TileSize) + visualOffsetY - TileSize/2) * Zoom
	}
}
