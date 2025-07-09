package main

import (
	"fmt"
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

	BounceSpeed           = 0.5
	FrightenedSpeedFactor = 0.5
)

type FrightState int

const (
	FrightBlue FrightState = iota
	FrightWhite
	FrightEyesUp
	FrightEyesRight
	FrightEyesDown
	FrightEyesLeft
)

const (
	BlinkyId GhostId = iota
	PinkyId
	InkyId
	ClydeId
)

type Behavior interface {
	Id() GhostId
	Color() rl.Color
	Sprite() Vec2i
	StartingTile(game *Game) Vec2i
	StartingDir(game *Game) Direction
	Chase(game *Game) Vec2i
	Scatter(game *Game) Vec2i
	ExitHouse(game *Game) bool
}

type Ghost struct {
	Entity
	id               GhostId
	state            GhostState
	frightState      FrightState
	behavior         Behavior
	fright           map[FrightState]Vec2i
	color            rl.Color
	target           Vec2i // temporary for training
	bounce           int
	pixelsMovedInDir float32
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

func NewGhost(game *Game, b Behavior) *Ghost {
	spriteY := b.Sprite().Y
	startX := b.StartingTile(game).X
	startY := b.StartingTile(game).Y
	dir := b.StartingDir(game)

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
			pixel:   rl.Vector2{X: float32(startX * Pixel), Y: float32(startY * Pixel)},
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
		fright: map[FrightState]Vec2i{
			FrightBlue:      {584, 64}, // 2 frames
			FrightWhite:     {616, 64}, // 2 frames
			FrightEyesUp:    {616, 80}, // 1 frame
			FrightEyesRight: {584, 80},
			FrightEyesDown:  {632, 80},
			FrightEyesLeft:  {600, 80},
		},
	}

	if g.InHouse() {
		g.state = InHouse
	} else {
		g.state = Scatter
	}
	fmt.Printf("ghost %s starting at %d,%d in state %s\n", g.id, g.tile.X, g.tile.Y, g.state)
	return &g
}

func (g *Ghost) ChooseDirection(game *Game, target Vec2i) Direction {

	var validDirections []Direction
	//if g.dir == None {
	//	return None
	//}
	for _, dir := range []Direction{Up, Down, Left, Right} {
		if dir == g.dir.Opposite() {
			continue
		}
		nextTile := dir.GetNextTile(g.tile)
		if game.maze.IsValidMove(nextTile) {
			validDirections = append(validDirections, dir)
		}
	}

	if len(validDirections) == 0 {
		fmt.Printf("**** no valid directions for %s at %v\n", g.id, target)
		return None
	}

	// commit to a (previous) decision for at least one tile
	if g.pixelsMovedInDir < Size {
		return g.dir
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

	//fmt.Printf("%s cur: %s: best %s at pos: %s, pixels: %0.2f / %0.2f\n", g.id, g.dir, bestDir, g.tile, g.pixelsMoved, g.pixelsMovedInDir)
	return bestDir
}

func (g *Ghost) Update(game *Game) {
	g.updateFrame()
	g.updateFright(game)
	g.updateState(game)
	if g.pixelsMoved >= Size {
		// Update tile position based on the last move
		g.tile = g.tile.Add(g.vel.X, g.vel.Y)
		// --- Tunnel Teleportation ---
		// if tunnel something

		g.pixelsMoved = 0
	}

	if g.state == Scatter {
		g.target = g.behavior.Scatter(game)
	} else if g.state == Chase {
		g.target = g.behavior.Chase(game) // ghost 0 is Blinky
	}

	curDir := g.dir
	if game.InTunnel(&g.Entity) {
		if g.tile.X < 0 {
			g.tile.X = GameWidth - 1
		} else if g.tile.X >= GameWidth-1 {
			g.tile.X = 0
		}
	} else {
		g.dir = g.ChooseDirection(game, g.target)
		if g.dir == None {
			fmt.Println("no direction")
			return
		}
	}

	g.vel = g.dir.Vector()

	if g.vel.IsNonZero() {

		currentSpeed := g.Speed(game)
		if curDir == g.dir {
			g.pixelsMovedInDir += currentSpeed
		} else {
			g.pixelsMovedInDir = currentSpeed
		}

		g.move(currentSpeed)
	}
}

func (g *Ghost) updateFrame() {
	if g.state == InHouse {
		g.frame = 0
		return
	} else if g.state == Frightened && g.frightState != FrightBlue && g.frightState != FrightWhite {
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

func (g *Ghost) updateFright(game *Game) {
	if g.state != Frightened {
		return
	}
	dt := game.frightTime - rl.GetTime()
	if dt < 0 {
		g.state = Scatter // temporary, set state will determine state
		g.updateState(game)
		return
	}

	if dt > 2.0 {
		g.frightState = FrightBlue
		return
	}

	n := int(math.Round(dt*200)) % 100
	if n > 49 {
		g.frightState = FrightBlue
	} else {
		g.frightState = FrightWhite
	}
}

func (g *Ghost) Speed(game *Game) float32 {
	speed := ghostSpeed(game.level)
	if g.state == Frightened {
		speed = frightSpeed(game.level)
	}

	// Apply tunnel speed reduction
	if game.InTunnel(&g.Entity) {
		speed *= TunnelSpeedFactor
	}

	return speed
}

func (g *Ghost) updateState(game *Game) {
	if game.debug || (g.state != Scatter && g.state != Chase) {
		return
	}

	level := game.level
	if level < 0 {
		level = 1
	}

	// a. Scatter (7s),
	// b. Chase (20s)
	// c. Scatter (7s)
	// d. Chase (20s)
	// e. Scatter (5s),
	// f. Chase (20s),
	// g. Scatter (5s),
	// h. Chase indefinitely.
	state := Chase
	// TODO - change based on levels
	if game.levelTime < 7 { // a
		state = Scatter
	} else if game.levelTime < 27 { // b
		state = Chase
	} else if game.levelTime < 34 { // c
		state = Scatter
	} else if game.levelTime < 54 { // d
		state = Chase
	} else if game.levelTime < 61 { // e
		state = Scatter
	} else if game.levelTime < 81 { // f
		state = Chase
	} else if game.levelTime < 86 { // g
		state = Scatter
	}

	if g.id == PinkyId && g.state != state {
		fmt.Printf("from %s to %s at %0.4f\n", g.state, state, game.levelTime)
	}
	g.state = state
}

func (g *Ghost) InHouse() bool {
	return g.tile.Y == 14 && (g.tile.X >= 12 && g.tile.X <= 16)
}

// GhostSpeed returns ghost speed in pixels per frame based on level
func ghostSpeed(level int) float32 {
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
func frightSpeed(level int) float32 {
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
func cruiseElroySpeed(level int, stage int) float32 {
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
