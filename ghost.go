package main

type GhostMode int

const (
	Scatter GhostMode = iota
	Chase
	Frightened
	Eaten
	InHouse
	LeavingHouse
)

func (g GhostMode) String() string {
	switch g {
	case Scatter:
		return "scatter"
	case Chase:
		return "chase"
	case Frightened:
		return "frightened"
	case Eaten:
		return "eaten"
	case InHouse:
		return "in house"
	case LeavingHouse:
		return "leaving house"
	default:
		panic("unhandled default case")
	}
}

// Blinky is the red ghost
type Blinky struct {
}

// Pinky is the pink ghost
type Pinky struct {
}

// Inky is the blue ghost
type Inky struct {
}

// Clyde is the orange ghost
type Clyde struct {
}

type Ghost interface {
	String() string
	StartingTile() Vec2i
	StartingShape() Shape
	Sprite() Vec2i
}
