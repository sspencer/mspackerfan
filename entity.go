package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	GhostSpeed    = 50   // speed through maze in level 1
	PlayerSpeed   = 60   // speed through maze in level 1
	SlowSpeed     = 45   // slow down about 15-25% after eating dots
	SlowTime      = 0.25 // slow down for this long
	TeleportSpeed = 10   // slow down about 15-25% after eating dots
	TeleportTime  = 0.75 // slow down for this long
)

type Direction int

const (
	None Direction = iota
	Up
	Right
	Down
	Left
)

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

func (d Direction) Vector() rl.Vector2 {
	switch d {
	case Up:
		return rl.Vector2{Y: -1}
	case Right:
		return rl.Vector2{X: 1}
	case Down:
		return rl.Vector2{Y: 1}
	case Left:
		return rl.Vector2{X: -1}
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

type Entity struct {
	name          string
	sprite        map[Direction]Vec2i // location in spritesheet
	dir           Direction
	nextDir       Direction
	vel           rl.Vector2
	nextVel       rl.Vector2
	tile          Vec2i
	pixel         rl.Vector2
	width         float32
	height        float32
	frameTime     float32
	frameSpeed    float32
	numFrames     int
	frame         int
	speed         float32
	slowTimer     float32
	teleportTimer float32
}
