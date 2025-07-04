package main

type Tile byte

const (
	Wall Tile = iota
	Dot
	Power
	Empty
	Tunnel

	DotMask   = 103481868288
	PowerMask = 4359202964317896252
	//DoorMask  = 16776960
)

func (t Tile) String() string {
	switch t {
	case Wall:
		return "wall"
	case Dot:
		return "dot"
	case Power:
		return "power"
	case Empty:
		return "empty"
	case Tunnel:
		return "tunnel"
	default:
		panic("unhandled default case")
	}
}

func (t Tile) Pretty() string {
	switch t {
	case Wall:
		return "XXX"
	case Dot:
		return " + "
	case Power:
		return "(*)"
	case Empty:
		return "   "
	case Tunnel:
		return "<@>"
	default:
		panic("unhandled default case")
	}
}
