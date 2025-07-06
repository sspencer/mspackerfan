package main

import rl "github.com/gen2brain/raylib-go/raylib"

// Pinky (Pink)
//    Chase:
//       Targets a tile four tiles ahead of Packer’s current position
//       in her direction of movement. If Packer is moving up,
//       Pinky’s target is four tiles above her (with an overflow bug in the
//       original game for upward movement).
//    Scatter:
//       Moves to the top-left corner of the maze.

func (b Pinky) Id() GhostId {
	return PinkyId
}

func (b Pinky) Color() rl.Color {
	return rl.Pink
}

func (b Pinky) StartingTile() Vec2i {
	return Vec2i{X: 14, Y: 14}
}

func (b Pinky) StartingDir() Direction {
	return Up
}

func (b Pinky) Sprite() Vec2i {
	return Vec2i{X: 520, Y: 80}
}

func (b Pinky) Chase(game *Game) Vec2i {
	p := game.player
	switch p.dir {
	case Up:
		if ChaseBug {
			// original game had error in logic code
			return p.tile.Add(-4, -4)
		} else {
			return p.tile.Add(0, -4)
		}
	case Down:
		return p.tile.Add(0, 4)
	case Left:
		return p.tile.Add(-4, 0)
	case Right:
		return p.tile.Add(4, 0)
	default:
		panic("unhandled default case")
	}
}

func (b Pinky) Scatter() Vec2i {
	return Vec2i{X: 1, Y: 1} // depends on board
}

func (b Pinky) ExitHouse(game *Game) bool {
	return game.levelTime > 1.0
	//return true
}
