package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	TunnelSpeedFactor = 0.5 // Pac-Man moves at 50% speed in tunnels
	DotEatPause       = 1   // 1 frame pause when eating regular dots
	PowerPelletPause  = 3   // 3 frames pause when eating power pellets
)

type Packer struct {
	Entity
	score       int
	pauseFrames int
	isEatingDot bool
}

func NewPacker() *Packer {
	startX := 13
	startY := 23
	shape := Left

	return &Packer{
		Entity: Entity{
			name: "ms. packer",
			sprite: map[Direction]Vec2i{
				Up:    {456, 32},
				Right: {456, 0},
				Down:  {456, 48},
				Left:  {456, 16},
			},
			pixel:   rl.Vector2{X: float32(startX * TileSize * Zoom), Y: float32(startY * TileSize * Zoom)},
			width:   16,
			height:  16,
			tile:    Vec2i{X: startX, Y: startY},
			dir:     shape,
			nextDir: shape,
			vel:     shape.Vector(),
			nextVel: shape.Vector(),
		},

		score: 0,
	}
}

// SetFrame sets the frame ID for Pac-Man's 4-frame animation.
// Frame IDs: 0 (fully open), 1 (half-open), 2 (closed).
// Animation cycle: closed (2) -> half-open (1) -> fully open (0) -> half-open (1).
// Each animation state lasts 5 frames at 60 FPS (20 frames per full cycle).
func (p *Packer) SetFrame() {
	framesPerState := 5
	cycleLength := 4 * framesPerState // 20 frames for full cycle
	frameInCycle := p.frameCount % cycleLength
	switch {
	case frameInCycle < framesPerState:
		p.frame = 2 // Closed
	case frameInCycle < 2*framesPerState:
		p.frame = 1 // Half-open
	case frameInCycle < 3*framesPerState:
		p.frame = 0 // Fully open
	default:
		p.frame = 1 // Half-open (returning)
	}

	p.frameCount++
	if p.frameCount > cycleLength {
		p.frameCount = 0
	}
}

func (p *Packer) Update(game *Game) {
	p.SetFrame() // animate mouth opening / closing

	if p.pauseFrames > 0 {
		p.pauseFrames--
		return
	}

	// --- Intersection and Direction Change Logic ---
	if p.pixelsMoved >= TileSize {
		// Update tile position based on the last move
		p.tile = p.tile.Add(p.vel.X, p.vel.Y)
		// --- Tunnel Teleportation ---
		//if p.TileY == 17 {
		//	if p.TileX <= -1 {
		//		p.TileX = TileCols - 1
		//	} else if p.TileX >= TileCols {
		//		p.TileX = 0
		//	}
		//}

		p.pixelsMoved = 0
		p.isEatingDot = false // Reset eating flag when reaching new tile

		// At the new intersection, decide the next move
		if p.canMove(game.maze, p.nextVel) {
			p.dir = p.nextDir
			p.vel = p.nextVel
		} else if !p.canMove(game.maze, p.vel) {
			p.vel = Vec2i{}
		}
	}

	// --- Handle starting from a standstill ---
	if p.vel.IsZero() {
		if p.canMove(game.maze, p.nextVel) {
			p.dir = p.nextDir
			p.vel = p.nextVel
		}
	}

	// --- Calculate current speed based on context ---
	currentSpeed := p.calculateCurrentSpeed()

	// --- Advance Pac-Man ---
	if p.vel.X != 0 || p.vel.Y != 0 {
		p.pixelsMoved += currentSpeed

		clampedPixelsMoved := float32(math.Min(float64(p.pixelsMoved), float64(TileSize)))
		visualOffsetX := float32(p.vel.X) * clampedPixelsMoved
		visualOffsetY := float32(p.vel.Y) * clampedPixelsMoved

		p.pixel.X = (float32(p.tile.X*TileSize) + visualOffsetX - TileSize/2) * Zoom
		p.pixel.Y = (float32(p.tile.Y*TileSize) + visualOffsetY - TileSize/2) * Zoom
	}

	// --- Eat Dots (only check once per tile entry) ---
	if p.tile.InMaze() && !p.isEatingDot {
		tile := game.maze[p.tile.Y][p.tile.X]
		// Check for power pellet first
		if tile == Power {
			game.maze[p.tile.Y][p.tile.X] = Empty
			p.pauseFrames = PowerPelletPause
			p.isEatingDot = true
		} else if tile == Dot {
			game.maze[p.tile.Y][p.tile.X] = Empty
			game.dotsEaten++
			p.pauseFrames = DotEatPause
			p.isEatingDot = true
		}
	}
}

func (p *Packer) calculateCurrentSpeed() float32 {
	speed := p.PackerSpeed(1)

	// Apply tunnel speed reduction
	if p.InTunnel() {
		speed *= TunnelSpeedFactor
	}

	return speed
}

func (p *Packer) InTunnel() bool {
	return false
}

func (p *Packer) canMove(maze Maze, dir Vec2i) bool {
	if dir.X == 0 && dir.Y == 0 {
		return false
	}

	nextTile := p.tile.Add(dir.X, dir.Y)

	// Special case: Allow movement into the tunnel off-screen
	//if nextTileY == 17 && (nextTileX < 0 || nextTileX >= TileCols) {
	//	return true
	//}

	// Check for moving off the map boundaries (non-tunnel)
	if nextTile.X < 0 || nextTile.X >= GameWidth || nextTile.Y < 0 || nextTile.Y >= GameHeight {
		return false
	}

	// Check for collision with a wall
	return maze[nextTile.Y][nextTile.X] != Wall
}

// PackerSpeed returns players's speed in pixels per frame based on level
// Level 1 returns 88.0 pixels/second (base speed)
func (p *Packer) PackerSpeed(level int) float32 {

	// Arcade-accurate speed progression based on research
	// Speeds are in pixels per second, converted to pixels per frame at 60 FPS
	var speedTable = []float32{
		88.0,  // Level 1
		96.8,  // Level 2 (110% of base)
		96.8,  // Level 3 (110% of base)
		96.8,  // Level 4 (110% of base)
		105.6, // Level 5+ (120% of base)
	}

	if level <= 0 {
		level = 1
	}

	var speed float32
	if level <= len(speedTable) {
		speed = speedTable[level-1]
	} else {
		// Level 5+ all use the same speed
		speed = speedTable[len(speedTable)-1]
	}

	return speed / 60.0 // Convert to pixels per frame
}

// PackerFrightSpeed returns player's speed when power pellet is active
func PackerFrightSpeed(level int) float32 {
	// During fright state (after eating power pellet)
	var frightSpeedTable = []float32{
		105.6, // Level 1-4: 120% of base speed
		105.6, // Level 2-4: 120% of base speed
		105.6, // Level 3-4: 120% of base speed
		105.6, // Level 4: 120% of base speed
		88.0,  // Level 5+: Back to base speed (100%)
	}

	if level <= 0 {
		level = 1
	}

	var speed float32
	if level <= len(frightSpeedTable) {
		speed = frightSpeedTable[level-1]
	} else {
		// Level 5+ all use base speed during fright
		speed = frightSpeedTable[len(frightSpeedTable)-1]
	}

	return speed / 60.0 // Convert to pixels per frame
}
