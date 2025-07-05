package main

// Blinky (Red):
//    Chase:
//       Targets Packer’s current tile directly, making Blinky the
//       most aggressive behavior.
//    Scatter:
//       Moves to the top-right corner of the maze.

func (h Blinky) String() string {
	return "Blinky"
}

func (h Blinky) StartingTile() Vec2i {
	return Vec2i{x: 9, y: 14}
}

func (h Blinky) StartingShape() Shape {
	return ShapeUp
}

func (h Blinky) Sprite() Vec2i {
	return Vec2i{x: 520, y: 64}
}
