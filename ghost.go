package main

import rl "github.com/gen2brain/raylib-go/raylib"

type GhostMode int

const (
	Scatter GhostMode = iota
	Chase
	Frightened
	Eaten
	InHouse
	LeavingHouse
)

func (g GhostMode) String() string {
	switch g {
	case Scatter:
		return "scatter"
	case Chase:
		return "chase"
	case Frightened:
		return "frightened"
	case Eaten:
		return "eaten"
	case InHouse:
		return "in house"
	case LeavingHouse:
		return "leaving house"
	default:
		panic("unhandled default case")
	}
}

// Blinky is the red behavior
type Blinky struct {
}

// Pinky is the pink behavior
type Pinky struct {
}

// Inky is the blue behavior
type Inky struct {
}

// Clyde is the orange behavior
type Clyde struct {
}

type Behavior interface {
	String() string
	StartingTile() Vec2i
	StartingShape() Shape
	Sprite() Vec2i
}

type Ghost struct {
	Entity
	ghostMode GhostMode
	behavior  Behavior
}

func createGhost(h Behavior) *Ghost {
	spriteY := h.Sprite().y
	startX := h.StartingTile().x
	startY := h.StartingTile().y
	shape := h.StartingShape()

	return &Ghost{
		Entity: Entity{
			name: h.String(),
			sprite: map[Shape]Vec2i{
				ShapeUp:    {520, spriteY},
				ShapeRight: {456, spriteY},
				ShapeDown:  {552, spriteY},
				ShapeLeft:  {488, spriteY},
			},
			tile:       Vec2i{x: startX, y: startY},
			pixel:      rl.Vector2{X: float32(startX * TileSize * Zoom), Y: float32(startY * TileSize * Zoom)},
			width:      16,
			height:     16,
			shape:      shape,
			nextShape:  shape,
			dir:        shape.Direction(),
			nextDir:    shape.Direction(),
			frameTime:  0.0,
			frameSpeed: 0.15,
			numFrames:  2,
			frame:      0,
			speed:      GhostSpeed * Zoom, // pixels per second
		},

		ghostMode: Scatter,
		behavior:  h,
	}
}
