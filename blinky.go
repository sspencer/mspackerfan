package main

import rl "github.com/gen2brain/raylib-go/raylib"

// Blinky (Red):
//    Chase:
//       Targets Packer’s current tile directly, making Blinky the
//       most aggressive behavior.
//    Scatter:
//       Moves to the top-right corner of the maze.

func (b Blinky) Id() GhostId {
	return BlinkyId
}

func (b Blinky) Color() rl.Color {
	return rl.Red
}

func (b Blinky) StartingTile() Vec2i {
	return Vec2i{X: 13, Y: 11}
}

func (b Blinky) StartingDir() Direction {
	return Left
}

func (b Blinky) Sprite() Vec2i {
	return Vec2i{X: 520, Y: 64}
}

func (b Blinky) Chase(packer *Packer, _ *Ghost) Vec2i {
	return packer.tile
}

func (b Blinky) Scatter() Vec2i {
	return Vec2i{X: 26, Y: 1} // depends on board
}
