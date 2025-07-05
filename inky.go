package main

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

func (h Inky) String() string {
	return "Inky"
}

func (h Inky) StartingTile() Vec2i {
	if trainingMode {
		return Vec2i{x: 5, y: 4}
	}
	return Vec2i{x: 12, y: 14}
}

func (h Inky) StartingShape() Shape {
	if trainingMode {
		return ShapeLeft
	}

	return ShapeDown
}

func (h Inky) Sprite() Vec2i {
	return Vec2i{x: 520, y: 96}
}
