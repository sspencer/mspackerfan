package main

import rl "github.com/gen2brain/raylib-go/raylib"

// Blinky (Red):
//    Chase:
//       Targets Player’s current tile directly, making Blinky the
//       most aggressive behavior.
//    Scatter:
//       Moves to the top-right corner of the maze.

func (b Blinky) Id() GhostId {
	return BlinkyId
}

func (b Blinky) Color() rl.Color {
	return rl.Red
}

func (b Blinky) StartingTile(game *Game) Vec2i {
	return Vec2i{X: 13, Y: 11}
}

func (b Blinky) StartingDir(game *Game) Direction {
	return Left
}

func (b Blinky) Sprite() Vec2i {
	return Vec2i{X: 520, Y: 64}
}

func (b Blinky) Chase(game *Game) Vec2i {
	return game.player.tile
}

func (b Blinky) Scatter(game *Game) Vec2i {
	return Vec2i{X: 26, Y: 1} // depends on board
}

func (b Blinky) ExitHouse(_ *Game) bool {
	return true
}
