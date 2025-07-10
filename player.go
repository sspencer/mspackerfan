package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	TunnelSpeedFactor = 0.5 // Pac-Man moves at 50% speed in tunnels
	DotEatPause       = 1   // 1 frame pause when eating regular dots
	PowerPelletPause  = 3   // 3 frames pause when eating power pellets
)

type Player struct {
	Entity
	score       int
	pauseFrames int
	isEatingDot bool
}

func NewPlayer() *Player {
	startX := 13
	startY := 23
	shape := Left

	return &Player{
		Entity: Entity{
			name: "ms. packer",
			sprite: map[Direction]Vec2i{
				Up:    {456, 32},
				Right: {456, 0},
				Down:  {456, 48},
				Left:  {456, 16},
			},
			pixel:     rl.Vector2{X: float32(startX * Pixel), Y: float32(startY * Pixel)},
			width:     16,
			height:    16,
			tile:      Vec2i{X: startX, Y: startY},
			dir:       shape,
			nextDir:   shape,
			vel:       shape.Vector(),
			nextVel:   shape.Vector(),
			speedTime: rl.GetTime() + SpeedTime,
		},

		score: 0,
	}
}

func (p *Player) Update(game *Game) {
	p.updateFrame() // animate mouth opening / closing

	if p.pauseFrames > 0 {
		p.pauseFrames--
		return
	}

	// --- Intersection and Direction Change Logic ---
	if p.pixelsMoved >= Size {
		// Update tile position based on the last move
		p.tile = p.tile.Add(p.vel.X, p.vel.Y)
		p.pixelsMoved = 0
		p.isEatingDot = false // Reset eating flag when reaching new tile

		// At the new intersection, decide the next move
		if p.canMove(game.maze, p.nextVel) {
			p.dir = p.nextDir
			p.vel = p.nextVel
		} else if game.InTunnel(&p.Entity) {
			if p.tile.X < 0 {
				p.tile.X = GameWidth - 1
			} else if p.tile.X >= GameWidth-1 {
				p.tile.X = 0
			}
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
	speed := p.calculateSpeed(game)
	p.move(speed)
	p.speedPixels += speed
	if p.speedTime-game.levelTime <= 0 {
		p.speedTime = game.levelTime + SpeedTime
		fmt.Printf("player moving %0.3f/s\n", p.speedPixels)
		p.speedPixels = 0
	}

	// --- Eat Dots (only check once per tile entry) ---
	if p.tile.InMaze() && !p.isEatingDot {
		tile := game.maze[p.tile.Y][p.tile.X]
		// Check for power pellet first
		if tile == Power {
			game.maze[p.tile.Y][p.tile.X] = Empty
			p.pauseFrames = PowerPelletPause
			p.isEatingDot = true
			game.setGhostMode(Frightened)
		} else if tile == Dot {
			game.maze[p.tile.Y][p.tile.X] = Empty
			game.dotsEaten++
			p.pauseFrames = DotEatPause
			p.isEatingDot = true
		}
	}
}

// updateFrame sets the frame ID for Pac-Man's 4-frame animation.
// Frame IDs: 0 (fully open), 1 (half-open), 2 (closed).
// Animation cycle: closed (2) -> half-open (1) -> fully open (0) -> half-open (1).
// Each animation state lasts 5 frames at 60 FPS (20 frames per full cycle).
func (p *Player) updateFrame() {
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

func (p *Player) calculateSpeed(game *Game) float32 {
	var speed float32

	if game.frightTime-rl.GetTime() > 0 {
		speed = playerFrightSpeed(game.level)
	} else {
		speed = playerSpeed(game.level)
	}

	// Apply tunnel speed reduction
	if game.InTunnel(&p.Entity) {
		speed *= TunnelSpeedFactor
	}

	return speed
}

func (p *Player) canMove(maze Maze, dir Vec2i) bool {
	if dir.X == 0 && dir.Y == 0 {
		return false
	}

	nextTile := p.tile.Add(dir.X, dir.Y)

	// Check for moving off the map boundaries (non-tunnel)
	if nextTile.X < 0 || nextTile.X >= GameWidth || nextTile.Y < 0 || nextTile.Y >= GameHeight {
		return false
	}

	// Check for collision with a wall
	return maze[nextTile.Y][nextTile.X] != Wall
}

// playerSpeed returns players's speed in pixels per frame based on level
// Level 1 returns 88.0 pixels/second (base speed)
func playerSpeed(level int) float32 {

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

	// TODO maybe remove cheat boost here, but player should be able to outrun
	// ghosts by a little -- garbage collection problem (but wouldn't that effect
	// all sprites)
	return (1.023 * speed) / 60.0 // Convert to pixels per frame
}

// playerFrightSpeed returns player's speed when power pellet is active
func playerFrightSpeed(level int) float32 {
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
