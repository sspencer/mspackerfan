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
	loc           map[Shape]Vec2i // location in spritesheet
	shape         Shape
	nextShape     Shape
	dir           rl.Vector2
	nextDir       rl.Vector2
	tileX         int
	tileY         int
	pixelX        float32
	pixelY        float32
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
}

// red 280, pink 296, blue 312, orange: 328
func createGhost(name string, startX, startY, spriteY int, shape Shape) *Entity {
	return &Entity{
		name: name,
		loc: map[Shape]Vec2i{
			ShapeUp:    {520, spriteY},
			ShapeRight: {456, spriteY},
			ShapeDown:  {552, spriteY},
			ShapeLeft:  {488, spriteY},
		},
		pixelX:     float32(startX * TileSize * Zoom),
		pixelY:     float32(startY * TileSize * Zoom),
		width:      16,
		height:     16,
		tileX:      startX, // * Zoom,
		tileY:      startY, // * Zoom,
		shape:      shape,
		nextShape:  shape,
		dir:        shape.Direction(),
		nextDir:    shape.Direction(),
		frameTime:  0.0,
		frameSpeed: 0.1,
		numFrames:  2,
		frame:      0,
		speed:      GhostSpeed * Zoom, // pixels per second
	}

}

func createPlayer(dots int) *Entity {

	startX := 13
	startY := 23
	shape := ShapeLeft

	return &Entity{
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
		tileX:      startX, // * Zoom,
		tileY:      startY, // * Zoom,
		shape:      shape,
		nextShape:  shape,
		dir:        shape.Direction(),
		nextDir:    shape.Direction(),
		frameTime:  0.0,
		frameSpeed: 0.1,
		numFrames:  3,
		frame:      0,
		speed:      PlayerSpeed * Zoom, // pixels per second
		dots:       dots,
	}
}

func (s Shape) Direction() rl.Vector2 {
	switch s {
	case ShapeUp:
		return rl.Vector2{0, -1}
	case ShapeRight:
		return rl.Vector2{1, 0}
	case ShapeDown:
		return rl.Vector2{0, 1}
	case ShapeLeft:
		return rl.Vector2{-1, 0}
	default:
		panic("unhandled default case")
	}
}
