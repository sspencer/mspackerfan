package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Inky (Blue)
//    Chase:
//       Uses a complex strategy. Inky calculates a target by:
//       1. Taking a point two tiles ahead of Packer in her direction.
//       2. Drawing a line from Blinky’s current position to that point.
//       3. Doubling the distance from Blinky to that point to select a target tile.
//          This makes Inky’s behavior dependent on Blinky’s position, often resulting
//          in flanking maneuvers.
//    Scatter:
//       Moves to the bottom-right corner of the maze.

func (b Inky) Id() GhostId {
	return InkyId
}

func (b Inky) Color() rl.Color {
	return rl.SkyBlue
}

func (b Inky) StartingTile() Vec2i {
	if trainingMode {
		return Vec2i{X: 5, Y: 4}
	}
	return Vec2i{X: 12, Y: 14}
}

func (b Inky) StartingDir() Direction {
	if trainingMode {
		return Left
	}

	return Down
}

func (b Inky) Sprite() Vec2i {
	return Vec2i{X: 520, Y: 96}
}

func (b Inky) Chase(p *Packer, g *Ghost) Vec2i {
	var offset Vec2i
	switch p.dir {
	case Up:
		if ChaseBug {
			offset.X = -2
			offset.Y = -2
		} else {
			offset.Y = -2
		}
	case Down:
		offset.Y = 2
	case Left:
		offset.X = -2
	case Right:
		offset.X = 2
	default:
		panic("unhandled default case")
	}

	pivot := p.tile.Add(offset.X, offset.Y)
	dx := pivot.X - g.tile.X
	dy := pivot.Y - g.tile.Y
	fmt.Printf("dx: %d, dy: %d\n", dx, dy)
	return Vec2i{X: g.tile.X + 2*dx, Y: g.tile.Y + 2*dy}
}

func (b Inky) Scatter() Vec2i {
	return Vec2i{X: 1, Y: 29} // depends on board
}
