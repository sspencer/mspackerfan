package main

import (
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
	return Vec2i{X: 12, Y: 14}
}

func (b Inky) StartingDir() Direction {
	return Down
}

func (b Inky) Sprite() Vec2i {
	return Vec2i{X: 520, Y: 96}
}

func (b Inky) Chase(game *Game) Vec2i {
	p := game.player
	g := game.ghosts[0] // Blinky
	var pivot Vec2i
	switch p.dir {
	case Up:
		if ChaseBug {
			pivot = p.tile.Add(-2, -2)
		} else {
			pivot = p.tile.Add(0, -2)
		}
	case Down:
		pivot = p.tile.Add(0, 2)
	case Left:
		pivot = p.tile.Add(-2, 0)
	case Right:
		pivot = p.tile.Add(2, 0)
	default:
		panic("unhandled default case")
	}

	dx := pivot.X - g.tile.X
	dy := pivot.Y - g.tile.Y

	return Vec2i{X: g.tile.X + 2*dx, Y: g.tile.Y + 2*dy}
}

func (b Inky) Scatter() Vec2i {
	return Vec2i{X: 1, Y: 29} // depends on board
}

func (b Inky) ExitHouse(game *Game) bool {
	if game.dotsEaten > 30 {
		return game.levelTime > 7
	}

	return false
}
