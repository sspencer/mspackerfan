package main

type GhostMode int

const (
	Scatter GhostMode = iota
	Chase
	Frightened
)

func (g GhostMode) String() string {
	switch g {
	case Scatter:
		return "scatter"
	case Chase:
		return "chase"
	case Frightened:
		return "frightened"
	default:
		panic("unhandled default case")
	}
}
