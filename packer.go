package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Packer struct {
	Entity
	score int
	dots  int
}

func createPacker(dots int) *Packer {
	startX := 13
	startY := 23
	shape := ShapeLeft

	return &Packer{
		Entity: Entity{
			name: "ms. packer",
			sprite: map[Shape]Vec2i{
				ShapeUp:    {456, 32},
				ShapeRight: {456, 0},
				ShapeDown:  {456, 48},
				ShapeLeft:  {456, 16},
			},
			pixel:      rl.Vector2{X: float32(startX * TileSize * Zoom), Y: float32(startY * TileSize * Zoom)},
			width:      16,
			height:     16,
			tile:       Vec2i{x: startX, y: startY},
			shape:      shape,
			nextShape:  shape,
			dir:        shape.Direction(),
			nextDir:    shape.Direction(),
			frameTime:  0.0,
			frameSpeed: 0.1,
			numFrames:  3,
			frame:      0,
			speed:      PlayerSpeed * Zoom, // pixels per second
		},

		score: 0,
		dots:  dots,
	}
}
