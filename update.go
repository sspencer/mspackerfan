package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (g *Game) Update(deltaTime float32) {
	p := g.player
	trainingMove := false
	if rl.IsKeyPressed(rl.KeyRight) {
		trainingMove = true
		p.nextDir = Right
		p.nextVel = rl.Vector2{X: 1}
	}

	if rl.IsKeyPressed(rl.KeyLeft) {
		trainingMove = true
		p.nextDir = Left
		p.nextVel = rl.Vector2{X: -1}
	}

	if rl.IsKeyPressed(rl.KeyUp) {
		trainingMove = true
		p.nextDir = Up
		p.nextVel = rl.Vector2{Y: -1}
	}

	if rl.IsKeyPressed(rl.KeyDown) {
		trainingMove = true
		p.nextDir = Down
		p.nextVel = rl.Vector2{Y: 1}
	}

	if rl.IsKeyPressed(rl.KeyD) {
		g.debug = !g.debug
	}

	if rl.IsKeyPressed(rl.KeyN) {
		g.boardNum = (g.boardNum + 1) % 6
		g.mapBoard()
	}

	if trainingMode {
		if rl.IsKeyPressed(rl.KeyC) {
			g.setGhostMode(Chase)
		}

		if rl.IsKeyPressed(rl.KeyS) {
			g.setGhostMode(Scatter)
		}

		if rl.IsKeyPressed(rl.KeyF) {
			g.setGhostMode(Frightened)
		}

		if rl.IsKeyPressed(rl.KeyF) {
			g.setGhostMode(Frightened)
		}

		if rl.IsKeyPressed(rl.KeySpace) {
			trainingMove = true
		}

		if rl.IsKeyPressed(rl.KeyT) {
			trainingMode = !trainingMode
		}
	}

	if rl.IsKeyPressed(rl.KeyP) {
		g.paused = !g.paused
	}

	score := 0
	if trainingMode {
		score = g.updatePlayerTraining()
		if trainingMove {
			g.updateGhostsTraining(deltaTime)
		}

	} else {
		if g.paused {
			return
		}
		score = p.updatePlayer(deltaTime, g.maze)
		g.updateGhosts(deltaTime)
	}

	if score > 0 {
		p.score += score
		p.dots -= 1
	}
	if p.score > g.highScore {
		g.highScore = p.score
	}
}

func (g *Game) setGhostMode(mode GhostMode) {
	for _, ghost := range g.ghosts {
		ghost.mode = mode
	}
}

func (g *Game) updateGhosts(dt float32) {
	for _, ghost := range g.ghosts {
		ghost.frameTime += dt

		if ghost.frameTime > ghost.frameSpeed {
			ghost.frame = (ghost.frame + 1) % ghost.numFrames
			ghost.frameTime -= ghost.frameSpeed
		}

		ghost.pixel.X = float32(ghost.tile.X * TileZoom)
		ghost.pixel.Y = float32(ghost.tile.Y * TileZoom)
	}
}

func (g *Game) updateGhostsTraining(dt float32) {
	for _, ghost := range g.ghosts {
		ghost.frameTime += dt

		if ghost.frameTime > ghost.frameSpeed {
			ghost.frame = (ghost.frame + 1) % ghost.numFrames
			ghost.frameTime -= ghost.frameSpeed
		}

		target := ghost.Update(g.player, g.ghosts)
		ghost.target = target
		dir := ghost.ChooseDirection(g.maze, target)

		if dir != None {
			vec := dir.Vector()
			ghost.dir = dir
			ghost.tile.X += int(vec.X)
			ghost.tile.Y += int(vec.Y)
		}

		ghost.pixel.X = float32(ghost.tile.X * TileZoom)
		ghost.pixel.Y = float32(ghost.tile.Y * TileZoom)
	}
}

func (g *Game) updatePlayerTraining() int {
	p := g.player
	score := 0
	if p.nextVel.X != 0 || p.nextVel.Y != 0 {
		newPos := p.tile.Add(int(p.nextVel.X), int(p.nextVel.Y))
		if newPos.X < 0 || newPos.X >= GameWidth || newPos.Y < 0 || newPos.Y >= GameHeight {
			return 0

		}
		p.dir = p.nextDir
		p.vel = p.nextVel
		p.tile.X = p.tile.X + int(p.nextVel.X)
		p.tile.Y = p.tile.Y + int(p.nextVel.Y)
		tile := g.maze[p.tile.Y][p.tile.X]
		if tile == Dot {
			score = 10
			g.maze[p.tile.Y][p.tile.X] = Empty
		} else if tile == Power {
			score = 50
			g.maze[p.tile.Y][p.tile.X] = Empty
		}
		p.pixel.X = float32(p.tile.X * TileZoom)
		p.pixel.Y = float32(p.tile.Y * TileZoom)
		p.nextVel = rl.Vector2{}
	}

	return score
}

func (p *Entity) updatePlayer(dt float32, maze Maze) int {
	// === ANIMATION ALWAYS ON ===
	if p.teleportTimer > 0 {
		p.speed = TeleportSpeed * Zoom
	} else if p.slowTimer > 0 {
		p.speed = SlowSpeed * Zoom
	} else {
		p.slowTimer = 0
		p.teleportTimer = 0
		p.speed = PlayerSpeed * Zoom
	}

	p.frameTime += dt
	p.slowTimer -= dt
	p.teleportTimer -= dt

	if p.frameTime > p.frameSpeed {
		p.frame = (p.frame + 1) % p.numFrames
		p.frameTime -= p.frameSpeed
	}
	// === END ANIMATION ALWAYS ON ===

	didMove := false
	// === CRITICAL FIX: Ensure tile coordinates are always up-to-date ===
	// This ensures that all subsequent calculations within this frame (especially turn logic)
	// use the player'p current tile, based on its precise pixel position.
	// Steve: fudge factor required, otherwise player can get stuck as a lot of
	// pixel coordinates come in a about 1/10 away from the desired number,
	// numbers like 927.919006 instead of 928
	// 	Y: 28 instead of 29 for 927.919006
	//	Y: 24 instead of 25 for 799.995911
	//	Y: 23 instead of 24 for 767.920776
	ff := float32(0.15)
	p.tile.X = int(math.Floor(float64(p.pixel.X+ff) / float64(TileZoom)))
	p.tile.Y = int(math.Floor(float64(p.pixel.Y+ff) / float64(TileZoom)))
	// =================================================================

	// Check for direction changes
	if p.nextVel.X != 0 || p.nextVel.Y != 0 {
		canChangeDir := false

		// Calculate the center of the current tile in pixel coordinates (top-left of tile + half tile size)
		currentTileCenterX := float32(p.tile.X*TileZoom) + float32(TileZoom/2)
		currentTileCenterY := float32(p.tile.Y*TileZoom) + float32(TileZoom/2)

		// Calculate player'p *actual* center (assuming pixelX/Y is top-left of player sprite)
		playerCenterX := p.pixel.X + float32(TileZoom/2)
		playerCenterY := p.pixel.Y + float32(TileZoom/2)

		// Define a tolerance for being "at" or "near" the center for turning.
		turnTolerance := 2.0

		// Condition 1: Entity is currently stopped. Try the next direction.
		if p.vel.X == 0 && p.vel.Y == 0 {
			if p.canMove(p.nextVel, maze) {
				// If stopped and can move, snap to the top-left of the current tile (pixel origin)
				p.pixel.X = float32(p.tile.X * TileZoom)
				p.pixel.Y = float32(p.tile.Y * TileZoom)
				canChangeDir = true
			} else {
				// Cannot move in that direction, clear nextVel.
				p.nextVel = rl.Vector2{0, 0}
			}
		} else if (p.vel.X != 0 && p.nextVel.Y != 0) || (p.vel.Y != 0 && p.nextVel.X != 0) {
			// Condition 2: Entity is moving and attempting a 90-degree turn.
			// Check if player is aligned enough on the perpendicular axis for a turn
			isAlignedForTurn := false
			if p.vel.X != 0 { // Currently moving horizontally, attempting vertical turn
				if math.Abs(float64(playerCenterY-currentTileCenterY)) < float64(turnTolerance) {
					isAlignedForTurn = true
				}
			} else if p.vel.Y != 0 { // Currently moving vertically, attempting horizontal turn
				if math.Abs(float64(playerCenterX-currentTileCenterX)) < float64(turnTolerance) {
					isAlignedForTurn = true
				}
			}

			if isAlignedForTurn {
				if p.canMove(p.nextVel, maze) {
					// === FIX: Snap BOTH coordinates to the current tile'p top-left corner ===
					// This ensures pixel-perfect alignment with the grid for the new direction.
					p.pixel.X = float32(p.tile.X * TileZoom)
					p.pixel.Y = float32(p.tile.Y * TileZoom)
					canChangeDir = true
				}
			}
		} else if (p.vel.X != 0 && p.nextVel.X == -p.vel.X) || (p.vel.Y != 0 && p.nextVel.Y == -p.vel.Y) {
			// Condition 3: Attempting a 180-degree turn (reverse direction).
			if p.canMove(p.nextVel, maze) {
				canChangeDir = true
			}
		} else if p.nextVel.X == p.vel.X && p.nextVel.Y == p.vel.Y {
			// Condition 4: Entity is trying to reinforce current direction.
			p.nextVel = rl.Vector2{0, 0}
		}

		// Apply the direction change if allowed
		if canChangeDir {
			p.dir = p.nextDir
			p.vel = p.nextVel
			p.nextVel = rl.Vector2{0, 0}
		}
	}

	// === REFINED MOVEMENT AND COLLISION HANDLING ===
	if p.vel.X != 0 || p.vel.Y != 0 {
		// Calculate the distance player would attempt to move this frame
		moveDistanceX := p.vel.X * p.speed * dt
		moveDistanceY := p.vel.Y * p.speed * dt

		// Entity'p current tile (based on top-left corner, which is updated at the top of function)
		currentTileX := p.tile.X
		currentTileY := p.tile.Y

		// Assume player sprite occupies one full tile for collision purposes
		playerSize := float32(TileZoom)

		// --- Handle Horizontal Movement ---
		if p.vel.X != 0 {
			// Calculate the tile the *leading edge* of the player would move into
			var leadingEdgePixelX float32
			var nextCheckTileX int

			if p.vel.X > 0 { // Moving Right
				leadingEdgePixelX = p.pixel.X + playerSize + moveDistanceX
				// Check the tile that the right edge of the player would cross into
				nextCheckTileX = int(math.Floor(float64(leadingEdgePixelX-0.001) / float64(TileZoom))) // -0.001 for float safety
			} else { // Moving Left
				leadingEdgePixelX = p.pixel.X + moveDistanceX
				// Check the tile that the left edge of the player would cross into
				nextCheckTileX = int(math.Floor(float64(leadingEdgePixelX) / float64(TileZoom)))
			}

			if nextCheckTileX == 0 && maze[p.tile.Y][0] == Tunnel {
				p.tile.X = GameWidth - 1
				p.pixel.X = float32(p.tile.X * TileZoom)
				p.teleportTimer = TeleportTime

				return 0
			} else if nextCheckTileX == GameWidth && maze[p.tile.Y][GameWidth-1] == Tunnel {
				p.tile.X = 0
				p.pixel.X = float32(p.tile.X * TileZoom)
				p.teleportTimer = TeleportTime
				return 0
			}

			// Check if the next tile (determined by leading edge) is a wall
			// We only check against the current row (currentTileY) as we are moving horizontally.
			if nextCheckTileX < 0 || nextCheckTileX >= GameWidth || currentTileY < 0 || currentTileY >= GameHeight || maze[currentTileY][nextCheckTileX] == Wall {
				// Collision in X direction
				if p.vel.X > 0 { // Moving Right
					// Calculate max allowed movement: up to the left edge of the wall tile
					// This means player'p right edge (p.pixel.X + playerSize) aligns with wall'p left edge.
					maxMoveX := (float32(nextCheckTileX * TileZoom)) - (p.pixel.X + playerSize)
					p.pixel.X += maxMoveX // Move only up to the collision point
					didMove = true

				} else { // Moving Left
					// Calculate max allowed movement: up to the right edge of the wall tile
					// This means player'p left edge (p.pixel.X) aligns with wall'p right edge.
					maxMoveX := (float32((nextCheckTileX + 1) * TileZoom)) - p.pixel.X
					p.pixel.X += maxMoveX // Move only up to the collision point
					didMove = true
				}
				p.vel.X = 0 // Stop horizontal movement
			} else {
				// No collision horizontally, apply full horizontal movement
				p.pixel.X += moveDistanceX
				didMove = true
			}
		}

		// --- Handle Vertical Movement ---
		if p.vel.Y != 0 {
			// Calculate the tile the *leading edge* of the player would move into
			var leadingEdgePixelY float32
			var nextCheckTileY int

			if p.vel.Y > 0 { // Moving Down
				leadingEdgePixelY = p.pixel.Y + playerSize + moveDistanceY
				// Check the tile that the bottom edge of the player would cross into
				nextCheckTileY = int(math.Floor(float64(leadingEdgePixelY-0.001) / float64(TileZoom))) // -0.001 for float safety
			} else { // Moving Up
				leadingEdgePixelY = p.pixel.Y + moveDistanceY
				// Check the tile that the top edge of the player would cross into
				nextCheckTileY = int(math.Floor(float64(leadingEdgePixelY) / float64(TileZoom)))
			}

			// Check if the next tile (determined by leading edge) is a wall
			// We only check against the current column (currentTileX) as we are moving vertically.
			if nextCheckTileY < 0 || nextCheckTileY >= GameHeight || currentTileX < 0 || currentTileX >= GameWidth || maze[nextCheckTileY][currentTileX] == Wall {
				// Collision in Y direction
				if p.vel.Y > 0 { // Moving Down
					// Calculate max allowed movement: up to the top edge of the wall tile
					// Entity'p bottom edge (p.pixel.Y + playerSize) aligns with wall'p top edge.
					maxMoveY := (float32(nextCheckTileY * TileZoom)) - (p.pixel.Y + playerSize)
					p.pixel.Y += maxMoveY // Move only up to the collision point
					didMove = true
				} else { // Moving Up
					// Calculate max allowed movement: up to the bottom edge of the wall tile
					// Entity'p top edge (p.pixel.Y) aligns with wall'p bottom edge.
					maxMoveY := (float32((nextCheckTileY + 1) * TileZoom)) - p.pixel.Y
					p.pixel.Y += maxMoveY // Move only up to the collision point
					didMove = true
				}
				p.vel.Y = 0 // Stop vertical movement
			} else {
				// No collision vertically, apply full vertical movement
				p.pixel.Y += moveDistanceY
				didMove = true
			}
		}
		// p.tile.X and p.tile.Y are updated at the top of the function based on current pixel position,
		// so no need to update them here again after pixel movement.
	}

	// === END REFINED MOVEMENT AND COLLISION HANDLING ===

	if didMove && p.tile.X >= 0 && p.tile.X < GameWidth {
		tile := maze[p.tile.Y][p.tile.X]
		if tile == Dot {
			maze[p.tile.Y][p.tile.X] = Empty
			p.slowTimer = SlowTime
			return 10
		} else if tile == Power {
			maze[p.tile.Y][p.tile.X] = Empty
			p.slowTimer = SlowTime
			return 50
		}
	}

	return 0
}

// Check if movement in a direction is possible
func (p *Entity) canMove(dir rl.Vector2, maze Maze) bool {
	if int(dir.X) == 0 && int(dir.Y) == 0 {
		return true // Standing still is always possible
	}

	nextTileX := p.tile.X + int(dir.X)
	nextTileY := p.tile.Y + int(dir.Y)

	// Check boundary conditions for the maze
	if nextTileX < 0 || nextTileX >= GameWidth || nextTileY < 0 || nextTileY >= GameHeight {
		return false // Cannot move out of bounds
	}

	// Check if the next tile is a wall
	return maze[nextTileY][nextTileX] != Wall
}
