package main

import (
	"math"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type GhostState int
type GhostId int

const (
	Scatter GhostState = iota
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

func (g GhostState) String() string {
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
	state    GhostState
	behavior Behavior
	color    rl.Color
	target   Vec2i // temporary for training
	bounce   int
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

func NewGhost(b Behavior) *Ghost {
	spriteY := b.Sprite().Y
	startX := b.StartingTile().X
	startY := b.StartingTile().Y
	dir := b.StartingDir()

	g := Ghost{
		Entity: Entity{
			name: b.Id().String(),
			sprite: map[Direction]Vec2i{
				Up:    {520, spriteY},
				Right: {456, spriteY},
				Down:  {552, spriteY},
				Left:  {488, spriteY},
			},
			tile:    Vec2i{X: startX, Y: startY},
			pixel:   rl.Vector2{X: float32(startX * TileSize * Zoom), Y: float32(startY * TileSize * Zoom)},
			width:   16,
			height:  16,
			dir:     dir,
			nextDir: dir,
			vel:     dir.Vector(),
			nextVel: dir.Vector(),
			frame:   0,
		},

		id:       b.Id(),
		color:    b.Color(),
		behavior: b,
		bounce:   1,
		state:    Scatter,
	}

	if g.InHouse() {
		g.state = InHouse
	}

	return &g
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

	if g.state == Frightened {
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

func (g *Ghost) SetFrame() {
	if g.state == InHouse {
		g.frame = 0
		return
	}
	framesPerState := 15
	cycleLength := 2 * framesPerState // 30 frames for full cycle
	frameInCycle := g.frameCount % cycleLength
	if frameInCycle < framesPerState {
		g.frame = 0 // First pose
	} else {
		g.frame = 1
	}
	g.frameCount++
	if g.frameCount > cycleLength {
		g.frameCount = 0
	}
}

func (g *Ghost) Update(game *Game) {
	g.SetFrame()

	if g.pixelsMoved >= TileSize {
		// Update tile position based on the last move
		// --- Tunnel Teleportation ---
		// if tunnel something

		g.tile = g.tile.Add(g.vel.X, g.vel.Y)
		g.pixelsMoved = 0
	}

	var currentSpeed float32 = 0.0
	if g.state == InHouse {
		currentSpeed = 0.1
		if g.pixelsMoved >= TileSize/2 {
			g.bounce *= -1
			g.dir = g.dir.Opposite()
		} else if g.pixelsMoved <= -TileSize/2 {
			g.bounce *= -1
			g.dir = g.dir.Opposite()
		}

		currentSpeed = float32(g.bounce) * 0.1

		//fmt.Printf("In house: %s dir=%s pixels: %f\n", g.name, g.dir, g.pixelsMoved)
	} else {
		if g.state == Scatter {
			g.target = g.behavior.Scatter()
		} else if g.state == Chase {
			if g.id == InkyId {
				g.target = g.behavior.Chase(game.player, game.ghosts[0]) // ghost 0 is Blinky
			} else {
				g.target = g.behavior.Chase(game.player, g)
			}
		}
		g.dir = g.ChooseDirection(game.maze, g.target)
		currentSpeed = g.calculateCurrentSpeed()
		g.vel = g.dir.Vector()
	}

	if g.vel.X != 0 || g.vel.Y != 0 {
		g.pixelsMoved += currentSpeed

		clampedPixelsMoved := float32(math.Min(float64(g.pixelsMoved), float64(TileSize)))
		visualOffsetX := float32(g.vel.X) * clampedPixelsMoved
		visualOffsetY := float32(g.vel.Y) * clampedPixelsMoved

		if g.state == InHouse {
			visualOffsetX += TileSize / 2
		}
		g.pixel.X = (float32(g.tile.X*TileSize) + visualOffsetX - TileSize/2) * Zoom
		g.pixel.Y = (float32(g.tile.Y*TileSize) + visualOffsetY - TileSize/2) * Zoom
	}
}

func (g *Ghost) InHouse() bool {
	return g.tile.Y == 14 && (g.tile.X >= 12 && g.tile.X <= 16)
}

func (g *Ghost) calculateCurrentSpeed() float32 {
	speed := g.GhostSpeed(1)

	// Apply tunnel speed reduction
	if g.InTunnel() {
		speed *= TunnelSpeedFactor
	}

	return speed
}

func (g *Ghost) InTunnel() bool {
	return false
}

// GhostSpeed returns ghost speed in pixels per frame based on level
func (g *Ghost) GhostSpeed(level int) float32 {
	// Ghost normal speed progression (slightly slower than Pac-Man)
	var ghostSpeedTable = []float32{
		84.48,   // Level 1: 96% of Pac-Man's speed
		92.928,  // Level 2: 96% of Pac-Man's speed
		92.928,  // Level 3: 96% of Pac-Man's speed
		92.928,  // Level 4: 96% of Pac-Man's speed
		101.376, // Level 5+: 96% of Pac-Man's speed
	}

	if level <= 0 {
		level = 1
	}

	var speed float32
	if level <= len(ghostSpeedTable) {
		speed = ghostSpeedTable[level-1]
	} else {
		speed = ghostSpeedTable[len(ghostSpeedTable)-1]
	}

	return speed / 60.0 // Convert to pixels per frame
}

// GhostFrightSpeed returns ghost speed when they're frightened (blue)
func GhostFrightSpeed(level int) float32 {
	// Frightened ghosts move at 50% of normal speed
	var frightSpeedTable = []float32{
		44.0, // Level 1: 50% of base speed
		48.4, // Level 2: 50% of level 2 speed
		48.4, // Level 3: 50% of level 3 speed
		48.4, // Level 4: 50% of level 4 speed
		52.8, // Level 5+: 50% of level 5+ speed
	}

	if level <= 0 {
		level = 1
	}

	var speed float32
	if level <= len(frightSpeedTable) {
		speed = frightSpeedTable[level-1]
	} else {
		speed = frightSpeedTable[len(frightSpeedTable)-1]
	}

	return speed / 60.0 // Convert to pixels per frame
}

// CruiseElroySpeed returns Blinky's speed when he becomes "Cruise Elroy"
func CruiseElroySpeed(level int, stage int) float32 {
	// Cruise Elroy has two stages of speed increase
	// Stage 1: When dots remaining <= threshold1
	// Stage 2: When dots remaining <= threshold2

	// TBD - find thresholds here
	// https://gamefaqs.gamespot.com/arcade/583976-ms-pac-man/faqs/1298
	var elroySpeedTable = [][]float32{
		{88.0, 96.8},   // Level 1: Stage 1=100%, Stage 2=110%
		{105.6, 114.4}, // Level 2: Stage 1=120%, Stage 2=130%
		{105.6, 114.4}, // Level 3: Stage 1=120%, Stage 2=130%
		{105.6, 114.4}, // Level 4: Stage 1=120%, Stage 2=130%
		{114.4, 123.2}, // Level 5+: Stage 1=130%, Stage 2=140%
	}

	if level <= 0 {
		level = 1
	}
	if stage < 1 {
		stage = 1
	}
	if stage > 2 {
		stage = 2
	}

	var speed float32
	if level <= len(elroySpeedTable) {
		speed = elroySpeedTable[level-1][stage-1]
	} else {
		speed = elroySpeedTable[len(elroySpeedTable)-1][stage-1]
	}

	return speed / 60.0 // Convert to pixels per frame
}
