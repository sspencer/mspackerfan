package main

// Pinky (Pink)
//    Chase:
//       Targets a tile four tiles ahead of Packer’s current position
//       in her direction of movement. If Packer is moving up,
//       Pinky’s target is four tiles above her (with an overflow bug in the
//       original game for upward movement).
//    Scatter:
//       Moves to the top-left corner of the maze.

func (h Pinky) String() string {
	return "Pinky"
}

func (h Pinky) StartingTile() Vec2i {
	if trainingMode {
		return Vec2i{x: 24, y: 21}
	}
	return Vec2i{x: 14, y: 14}
}

func (h Pinky) StartingShape() Shape {
	if trainingMode {
		return ShapeDown
	}
	return ShapeUp
}

func (h Pinky) Sprite() Vec2i {
	return Vec2i{x: 520, y: 80}
}
