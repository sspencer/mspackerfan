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
	shape := Left

	return &Packer{
		Entity: Entity{
			name: "ms. packer",
			sprite: map[Direction]Vec2i{
				Up:    {456, 32},
				Right: {456, 0},
				Down:  {456, 48},
				Left:  {456, 16},
			},
			pixel:      rl.Vector2{X: float32(startX * TileSize * Zoom), Y: float32(startY * TileSize * Zoom)},
			width:      16,
			height:     16,
			tile:       Vec2i{X: startX, Y: startY},
			dir:        shape,
			nextDir:    shape,
			vel:        shape.Vector(),
			nextVel:    shape.Vector(),
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
