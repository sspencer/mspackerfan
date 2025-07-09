package main

import rl "github.com/gen2brain/raylib-go/raylib"

// Clyde (Orange)
//    Chase:
//       If Clyde is more than eight tiles away from Player, he targets
//       her directly (like Blinky). If he is closer than eight tiles, he
//       retreats to his Scatter target (bottom-left corner).
//    Scatter:
//       Moves to the bottom-left corner of the maze.

func (b Clyde) Id() GhostId {
	return ClydeId
}

func (b Clyde) Color() rl.Color {
	return rl.Orange
}

func (b Clyde) StartingTile(game *Game) Vec2i {
	//if game.debug {
	return Vec2i{X: 16, Y: 11}
	//}
	//return Vec2i{X: 16, Y: 14}
}

func (b Clyde) StartingDir(game *Game) Direction {
	//if game.debug {
	return Right
	//}
	//return Down
}

func (b Clyde) Sprite() Vec2i {
	return Vec2i{X: 520, Y: 112}
}

func (b Clyde) Chase(game *Game) Vec2i {
	// find self
	for _, ghost := range game.ghosts {
		if ghost.id == b.Id() {
			d := game.player.tile.Distance(ghost.tile)
			if d > 8 {
				return game.player.tile
			}
			break
		}
	}

	return b.Scatter(game)
}

func (b Clyde) Scatter(game *Game) Vec2i {
	return Vec2i{X: 26, Y: 29} // depends on board
}

func (b Clyde) ExitHouse(game *Game) bool {
	if game.debug {
		return game.levelTime > 3.0
	}
	if game.dotsEaten > 60 {
		return game.levelTime > 15
	}

	return false
}
