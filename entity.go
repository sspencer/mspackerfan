package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	GhostSpeed    = 50   // speed through maze in level 1
	PlayerSpeed   = 60   // speed through maze in level 1
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

type Entity struct {
	name          string
	sprite        map[Shape]Vec2i // location in spritesheet
	shape         Shape
	nextShape     Shape
	dir           rl.Vector2
	nextDir       rl.Vector2
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

func (s Shape) Direction() rl.Vector2 {
	switch s {
	case ShapeUp:
		return rl.Vector2{Y: -1}
	case ShapeRight:
		return rl.Vector2{X: 1}
	case ShapeDown:
		return rl.Vector2{Y: 1}
	case ShapeLeft:
		return rl.Vector2{X: -1}
	default:
		panic("unhandled default case")
	}
}
