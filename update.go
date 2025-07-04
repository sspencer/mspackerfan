package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (g *Game) Update(deltaTime float32) {
	p := g.player
	if rl.IsKeyPressed(rl.KeyRight) {
		p.nextShape = ShapeRight
		p.nextDir = rl.Vector2{1, 0}
	}

	if rl.IsKeyPressed(rl.KeyLeft) {
		p.nextShape = ShapeLeft
		p.nextDir = rl.Vector2{-1, 0}
	}

	if rl.IsKeyPressed(rl.KeyUp) {
		p.nextShape = ShapeUp
		p.nextDir = rl.Vector2{0, -1}
	}

	if rl.IsKeyPressed(rl.KeyDown) {
		p.nextShape = ShapeDown
		p.nextDir = rl.Vector2{0, 1}
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
			g.ghostMode = Chase
		}

		if rl.IsKeyPressed(rl.KeyS) {
			g.ghostMode = Scatter
		}

		if rl.IsKeyPressed(rl.KeyF) {
			g.ghostMode = Frightened
		}
	}

	if rl.IsKeyPressed(rl.KeyP) {
		g.paused = !g.paused
	}

	score := 0
	if trainingMode {
		score = g.updateTrainingMode()
	} else {
		if g.paused {
			return
		}
		for _, e := range g.ghosts {
			e.updateGhost(deltaTime)
		}

		score = p.updatePlayer(deltaTime, g.board)
	}

	if score > 0 {
		p.score += score
		p.dots -= 1
	}
	if p.score > g.highScore {
		g.highScore = p.score
	}
}

func (g *Game) updateTrainingMode() int {
	p := g.player
	score := 0
	if p.nextDir.X != 0 || p.nextDir.Y != 0 {
		p.shape = p.nextShape
		p.dir = p.nextDir
		p.tile.x = p.tile.x + int(p.nextDir.X)
		p.tile.y = p.tile.y + int(p.nextDir.Y)
		tile := g.board[p.tile.y][p.tile.x]
		if tile == Dot {
			score = 10
			g.board[p.tile.y][p.tile.x] = Empty
		} else if tile == Power {
			score = 50
			g.board[p.tile.y][p.tile.x] = Empty
		}
		p.pixel.X = float32(p.tile.x * TileZoom)
		p.pixel.Y = float32(p.tile.y * TileZoom)
		p.nextDir = rl.Vector2{0, 0}
	}

	return score
}

func (p *Ghost) updateGhost(dt float32) {
	p.frameTime += dt

	if p.frameTime > p.frameSpeed {
		p.frame = (p.frame + 1) % p.numFrames
		p.frameTime -= p.frameSpeed
	}

	//p.pixel.X = float32(p.tile.x * TileZoom)
	//p.pixel.Y = float32(p.tile.y * TileZoom)

	p.pixel.X = float32(p.tile.x * TileZoom)
	p.pixel.Y = float32(p.tile.y * TileZoom)
}

func (p *Entity) updatePlayer(dt float32, board [][]Tile) int {
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
	p.tile.x = int(math.Floor(float64(p.pixel.X+ff) / float64(TileZoom)))
	p.tile.y = int(math.Floor(float64(p.pixel.Y+ff) / float64(TileZoom)))
	// =================================================================

	// Check for direction changes
	if p.nextDir.X != 0 || p.nextDir.Y != 0 {
		canChangeDir := false

		// Calculate the center of the current tile in pixel coordinates (top-left of tile + half tile size)
		currentTileCenterX := float32(p.tile.x*TileZoom) + float32(TileZoom/2)
		currentTileCenterY := float32(p.tile.y*TileZoom) + float32(TileZoom/2)

		// Calculate player'p *actual* center (assuming pixelX/Y is top-left of player sprite)
		playerCenterX := p.pixel.X + float32(TileZoom/2)
		playerCenterY := p.pixel.Y + float32(TileZoom/2)

		// Define a tolerance for being "at" or "near" the center for turning.
		turnTolerance := 2.0

		// Condition 1: Entity is currently stopped. Try the next direction.
		if p.dir.X == 0 && p.dir.Y == 0 {
			if p.canMove(p.nextDir, board) {
				// If stopped and can move, snap to the top-left of the current tile (pixel origin)
				p.pixel.X = float32(p.tile.x * TileZoom)
				p.pixel.Y = float32(p.tile.y * TileZoom)
				canChangeDir = true
			} else {
				// Cannot move in that direction, clear nextDir.
				p.nextDir = rl.Vector2{0, 0}
			}
		} else if (p.dir.X != 0 && p.nextDir.Y != 0) || (p.dir.Y != 0 && p.nextDir.X != 0) {
			// Condition 2: Entity is moving and attempting a 90-degree turn.
			// Check if player is aligned enough on the perpendicular axis for a turn
			isAlignedForTurn := false
			if p.dir.X != 0 { // Currently moving horizontally, attempting vertical turn
				if math.Abs(float64(playerCenterY-currentTileCenterY)) < float64(turnTolerance) {
					isAlignedForTurn = true
				}
			} else if p.dir.Y != 0 { // Currently moving vertically, attempting horizontal turn
				if math.Abs(float64(playerCenterX-currentTileCenterX)) < float64(turnTolerance) {
					isAlignedForTurn = true
				}
			}

			if isAlignedForTurn {
				if p.canMove(p.nextDir, board) {
					// === FIX: Snap BOTH coordinates to the current tile'p top-left corner ===
					// This ensures pixel-perfect alignment with the grid for the new direction.
					p.pixel.X = float32(p.tile.x * TileZoom)
					p.pixel.Y = float32(p.tile.y * TileZoom)
					canChangeDir = true
				}
			}
		} else if (p.dir.X != 0 && p.nextDir.X == -p.dir.X) || (p.dir.Y != 0 && p.nextDir.Y == -p.dir.Y) {
			// Condition 3: Attempting a 180-degree turn (reverse direction).
			if p.canMove(p.nextDir, board) {
				canChangeDir = true
			}
		} else if p.nextDir.X == p.dir.X && p.nextDir.Y == p.dir.Y {
			// Condition 4: Entity is trying to reinforce current direction.
			p.nextDir = rl.Vector2{0, 0}
		}

		// Apply the direction change if allowed
		if canChangeDir {
			p.shape = p.nextShape
			p.dir = p.nextDir
			p.nextDir = rl.Vector2{0, 0}
		}
	}

	// === REFINED MOVEMENT AND COLLISION HANDLING ===
	if p.dir.X != 0 || p.dir.Y != 0 {
		// Calculate the distance player would attempt to move this frame
		moveDistanceX := p.dir.X * p.speed * dt
		moveDistanceY := p.dir.Y * p.speed * dt

		// Entity'p current tile (based on top-left corner, which is updated at the top of function)
		currentTileX := p.tile.x
		currentTileY := p.tile.y

		// Assume player sprite occupies one full tile for collision purposes
		playerSize := float32(TileZoom)

		// --- Handle Horizontal Movement ---
		if p.dir.X != 0 {
			// Calculate the tile the *leading edge* of the player would move into
			var leadingEdgePixelX float32
			var nextCheckTileX int

			if p.dir.X > 0 { // Moving Right
				leadingEdgePixelX = p.pixel.X + playerSize + moveDistanceX
				// Check the tile that the right edge of the player would cross into
				nextCheckTileX = int(math.Floor(float64(leadingEdgePixelX-0.001) / float64(TileZoom))) // -0.001 for float safety
			} else { // Moving Left
				leadingEdgePixelX = p.pixel.X + moveDistanceX
				// Check the tile that the left edge of the player would cross into
				nextCheckTileX = int(math.Floor(float64(leadingEdgePixelX) / float64(TileZoom)))
			}

			if nextCheckTileX == 0 && board[p.tile.y][0] == Tunnel {
				p.tile.x = GameWidth - 1
				p.pixel.X = float32(p.tile.x * TileZoom)
				p.teleportTimer = TeleportTime

				return 0
			} else if nextCheckTileX == GameWidth && board[p.tile.y][GameWidth-1] == Tunnel {
				p.tile.x = 0
				p.pixel.X = float32(p.tile.x * TileZoom)
				p.teleportTimer = TeleportTime
				return 0
			}

			// Check if the next tile (determined by leading edge) is a wall
			// We only check against the current row (currentTileY) as we are moving horizontally.
			if nextCheckTileX < 0 || nextCheckTileX >= GameWidth || currentTileY < 0 || currentTileY >= GameHeight || board[currentTileY][nextCheckTileX] == Wall {
				// Collision in X direction
				if p.dir.X > 0 { // Moving Right
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
				p.dir.X = 0 // Stop horizontal movement
			} else {
				// No collision horizontally, apply full horizontal movement
				p.pixel.X += moveDistanceX
				didMove = true
			}
		}

		// --- Handle Vertical Movement ---
		if p.dir.Y != 0 {
			// Calculate the tile the *leading edge* of the player would move into
			var leadingEdgePixelY float32
			var nextCheckTileY int

			if p.dir.Y > 0 { // Moving Down
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
			if nextCheckTileY < 0 || nextCheckTileY >= GameHeight || currentTileX < 0 || currentTileX >= GameWidth || board[nextCheckTileY][currentTileX] == Wall {
				// Collision in Y direction
				if p.dir.Y > 0 { // Moving Down
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
				p.dir.Y = 0 // Stop vertical movement
			} else {
				// No collision vertically, apply full vertical movement
				p.pixel.Y += moveDistanceY
				didMove = true
			}
		}
		// p.tile.x and p.tile.y are updated at the top of the function based on current pixel position,
		// so no need to update them here again after pixel movement.
	}

	// === END REFINED MOVEMENT AND COLLISION HANDLING ===

	if didMove && p.tile.x >= 0 && p.tile.x < GameWidth {
		tile := board[p.tile.y][p.tile.x]
		if tile == Dot {
			board[p.tile.y][p.tile.x] = Empty
			p.slowTimer = SlowTime
			return 10
		} else if tile == Power {
			board[p.tile.y][p.tile.x] = Empty
			p.slowTimer = SlowTime
			return 50
		}
	}

	return 0
}

// Check if movement in a direction is possible
func (p *Entity) canMove(dir rl.Vector2, board [][]Tile) bool {
	if int(dir.X) == 0 && int(dir.Y) == 0 {
		return true // Standing still is always possible
	}

	nextTileX := p.tile.x + int(dir.X)
	nextTileY := p.tile.y + int(dir.Y)

	// Check boundary conditions for the board
	if nextTileX < 0 || nextTileX >= GameWidth || nextTileY < 0 || nextTileY >= GameHeight {
		return false // Cannot move out of bounds
	}

	// Check if the next tile is a wall
	return board[nextTileY][nextTileX] != Wall
}
