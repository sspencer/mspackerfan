package main

// Clyde (Orange)
//    Chase:
//       If Clyde is more than eight tiles away from Packer, he targets
//       her directly (like Blinky). If he is closer than eight tiles, he
//       retreats to his Scatter target (bottom-left corner).
//    Scatter:
//       Moves to the bottom-left corner of the maze.

func (h Clyde) String() string {
	return "Clyde"
}

func (h Clyde) StartingTile() Vec2i {
	if trainingMode {
		return Vec2i{x: 15, y: 17}
	}
	return Vec2i{x: 16, y: 14}
}

func (h Clyde) StartingShape() Shape {
	if trainingMode {
		return ShapeLeft
	}

	return ShapeDown
}

func (h Clyde) Sprite() Vec2i {
	return Vec2i{x: 520, y: 112}
}
