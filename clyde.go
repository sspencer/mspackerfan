package main

import rl "github.com/gen2brain/raylib-go/raylib"

// Clyde (Orange)
//    Chase:
//       If Clyde is more than eight tiles away from Packer, he targets
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

func (b Clyde) StartingTile() Vec2i {
	if trainingMode {
		return Vec2i{X: 15, Y: 17}
	}
	return Vec2i{X: 16, Y: 14}
}

func (b Clyde) StartingDir() Direction {
	if trainingMode {
		return Left
	}

	return Down
}

func (b Clyde) Sprite() Vec2i {
	return Vec2i{X: 520, Y: 112}
}

func (b Clyde) Chase(p *Packer, g *Ghost) Vec2i {
	d := p.tile.Distance(g.tile)
	if d > 8 {
		return p.tile
	}

	return b.Scatter()
}

func (b Clyde) Scatter() Vec2i {
	return Vec2i{X: 26, Y: 29} // depends on board
}
