package main

import (
	"math"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type GhostMode int
type GhostId int

const (
	Scatter GhostMode = iota
	Chase
	Frightened
	Eaten
	InHouse
	LeavingHouse
)

const (
	BlinkyId GhostId = iota
	PinkyId
	InkyId
	ClydeId
)

func (g GhostId) String() string {
	switch g {
	case BlinkyId:
		return "Blinky"
	case PinkyId:
		return "Pinky"
	case InkyId:
		return "Inky"
	case ClydeId:
		return "Clyde"
	default:
		panic("unhandled default case")
	}
}

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

type Behavior interface {
	Id() GhostId
	Color() rl.Color
	StartingTile() Vec2i
	StartingDir() Direction
	Sprite() Vec2i
	Chase(packer *Packer, ghost *Ghost) Vec2i
	Scatter() Vec2i
}

type Ghost struct {
	Entity
	id       GhostId
	mode     GhostMode
	behavior Behavior
	color    rl.Color
	target   Vec2i // temporary for training
}

// Blinky is the red behavior
type Blinky struct {
}

// Pinky is the pink behavior
type Pinky struct {
}

// Inky is the blue behavior
type Inky struct {
}

// Clyde is the orange behavior
type Clyde struct {
}

func createGhost(b Behavior) *Ghost {
	spriteY := b.Sprite().Y
	startX := b.StartingTile().X
	startY := b.StartingTile().Y
	dir := b.StartingDir()

	return &Ghost{
		Entity: Entity{
			name: b.Id().String(),
			sprite: map[Direction]Vec2i{
				Up:    {520, spriteY},
				Right: {456, spriteY},
				Down:  {552, spriteY},
				Left:  {488, spriteY},
			},
			tile:       Vec2i{X: startX, Y: startY},
			pixel:      rl.Vector2{X: float32(startX * TileSize * Zoom), Y: float32(startY * TileSize * Zoom)},
			width:      16,
			height:     16,
			dir:        dir,
			nextDir:    dir,
			vel:        dir.Vector(),
			nextVel:    dir.Vector(),
			frameTime:  0.0,
			frameSpeed: 0.15,
			numFrames:  2,
			frame:      0,
			speed:      GhostSpeed * Zoom, // pixels per second
		},

		id:       b.Id(),
		color:    b.Color(),
		mode:     Scatter,
		behavior: b,
	}
}

func (g *Ghost) ChooseDirection(maze Maze, target Vec2i) Direction {
	var validDirections []Direction
	for _, dir := range []Direction{Up, Down, Left, Right} {
		if dir == g.dir.Opposite() {
			continue
		}
		nextTile := dir.GetNextTile(g.tile)
		if maze.IsValidMove(nextTile) {
			validDirections = append(validDirections, dir)
		}
	}

	if len(validDirections) == 0 {
		return None
	}

	if g.mode == Frightened {
		return validDirections[rand.IntN(len(validDirections))]
	}

	bestDir := validDirections[0]
	minDist := float32(math.MaxFloat32)
	for _, dir := range validDirections {
		nextTile := dir.GetNextTile(g.tile)
		dist := target.Distance(nextTile)
		if dist < minDist {
			minDist = dist
			bestDir = dir
		}
	}
	return bestDir
}

func (g *Ghost) Update(p *Packer, ghosts []*Ghost) Vec2i {
	if g.mode == Scatter {
		return g.behavior.Scatter()
	} else if g.mode == Chase {
		if g.id == InkyId {
			return g.behavior.Chase(p, ghosts[0]) // ghost 0 is Blinky
		} else {
			return g.behavior.Chase(p, g)
		}
	}

	return Vec2i{}
}
