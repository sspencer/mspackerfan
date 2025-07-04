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

	if rl.IsKeyPressed(rl.KeyP) {
		g.paused = !g.paused
	}

	if g.paused {
		return
	}

	for _, e := range g.ghosts {
		e.updateGhost(deltaTime)
	}

	score := p.updatePlayer(deltaTime, g.board)
	if score > 0 {
		p.score += score
		p.dots -= 1
	}
	if p.score > g.highScore {
		g.highScore = p.score
	}
}

func (p *Entity) updateGhost(dt float32) {
	p.frameTime += dt

	if p.frameTime > p.frameSpeed {
		p.frame = (p.frame + 1) % p.numFrames
		p.frameTime -= p.frameSpeed
	}

	p.pixelX = float32(p.tileX * TileZoom)
	p.pixelY = float32(p.tileY * TileZoom)

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
	p.tileX = int(math.Floor(float64(p.pixelX+ff) / float64(TileZoom)))
	p.tileY = int(math.Floor(float64(p.pixelY+ff) / float64(TileZoom)))
	// =================================================================

	// Check for direction changes
	if p.nextDir.X != 0 || p.nextDir.Y != 0 {
		canChangeDir := false

		// Calculate the center of the current tile in pixel coordinates (top-left of tile + half tile size)
		currentTileCenterX := float32(p.tileX*TileZoom) + float32(TileZoom/2)
		currentTileCenterY := float32(p.tileY*TileZoom) + float32(TileZoom/2)

		// Calculate player'p *actual* center (assuming pixelX/Y is top-left of player sprite)
		playerCenterX := p.pixelX + float32(TileZoom/2)
		playerCenterY := p.pixelY + float32(TileZoom/2)

		// Define a tolerance for being "at" or "near" the center for turning.
		turnTolerance := 2.0

		// Condition 1: Entity is currently stopped. Try the next direction.
		if p.dir.X == 0 && p.dir.Y == 0 {
			if p.canMove(p.nextDir, board) {
				// If stopped and can move, snap to the top-left of the current tile (pixel origin)
				p.pixelX = float32(p.tileX * TileZoom)
				p.pixelY = float32(p.tileY * TileZoom)
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
					p.pixelX = float32(p.tileX * TileZoom)
					p.pixelY = float32(p.tileY * TileZoom)
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
		currentTileX := p.tileX
		currentTileY := p.tileY

		// Assume player sprite occupies one full tile for collision purposes
		playerSize := float32(TileZoom)

		// --- Handle Horizontal Movement ---
		if p.dir.X != 0 {
			// Calculate the tile the *leading edge* of the player would move into
			var leadingEdgePixelX float32
			var nextCheckTileX int

			if p.dir.X > 0 { // Moving Right
				leadingEdgePixelX = p.pixelX + playerSize + moveDistanceX
				// Check the tile that the right edge of the player would cross into
				nextCheckTileX = int(math.Floor(float64(leadingEdgePixelX-0.001) / float64(TileZoom))) // -0.001 for float safety
			} else { // Moving Left
				leadingEdgePixelX = p.pixelX + moveDistanceX
				// Check the tile that the left edge of the player would cross into
				nextCheckTileX = int(math.Floor(float64(leadingEdgePixelX) / float64(TileZoom)))
			}

			if nextCheckTileX == 0 && board[p.tileY][0] == Tunnel {
				p.tileX = GameWidth - 1
				p.pixelX = float32(p.tileX * TileZoom)
				p.teleportTimer = TeleportTime

				return 0
			} else if nextCheckTileX == GameWidth && board[p.tileY][GameWidth-1] == Tunnel {
				p.tileX = 0
				p.pixelX = float32(p.tileX * TileZoom)
				p.teleportTimer = TeleportTime
				return 0
			}

			// Check if the next tile (determined by leading edge) is a wall
			// We only check against the current row (currentTileY) as we are moving horizontally.
			if nextCheckTileX < 0 || nextCheckTileX >= GameWidth || currentTileY < 0 || currentTileY >= GameHeight || board[currentTileY][nextCheckTileX] == Wall {
				// Collision in X direction
				if p.dir.X > 0 { // Moving Right
					// Calculate max allowed movement: up to the left edge of the wall tile
					// This means player'p right edge (p.pixelX + playerSize) aligns with wall'p left edge.
					maxMoveX := (float32(nextCheckTileX * TileZoom)) - (p.pixelX + playerSize)
					p.pixelX += maxMoveX // Move only up to the collision point
					didMove = true

				} else { // Moving Left
					// Calculate max allowed movement: up to the right edge of the wall tile
					// This means player'p left edge (p.pixelX) aligns with wall'p right edge.
					maxMoveX := (float32((nextCheckTileX + 1) * TileZoom)) - p.pixelX
					p.pixelX += maxMoveX // Move only up to the collision point
					didMove = true
				}
				p.dir.X = 0 // Stop horizontal movement
			} else {
				// No collision horizontally, apply full horizontal movement
				p.pixelX += moveDistanceX
				didMove = true
			}
		}

		// --- Handle Vertical Movement ---
		if p.dir.Y != 0 {
			// Calculate the tile the *leading edge* of the player would move into
			var leadingEdgePixelY float32
			var nextCheckTileY int

			if p.dir.Y > 0 { // Moving Down
				leadingEdgePixelY = p.pixelY + playerSize + moveDistanceY
				// Check the tile that the bottom edge of the player would cross into
				nextCheckTileY = int(math.Floor(float64(leadingEdgePixelY-0.001) / float64(TileZoom))) // -0.001 for float safety
			} else { // Moving Up
				leadingEdgePixelY = p.pixelY + moveDistanceY
				// Check the tile that the top edge of the player would cross into
				nextCheckTileY = int(math.Floor(float64(leadingEdgePixelY) / float64(TileZoom)))
			}

			// Check if the next tile (determined by leading edge) is a wall
			// We only check against the current column (currentTileX) as we are moving vertically.
			if nextCheckTileY < 0 || nextCheckTileY >= GameHeight || currentTileX < 0 || currentTileX >= GameWidth || board[nextCheckTileY][currentTileX] == Wall {
				// Collision in Y direction
				if p.dir.Y > 0 { // Moving Down
					// Calculate max allowed movement: up to the top edge of the wall tile
					// Entity'p bottom edge (p.pixelY + playerSize) aligns with wall'p top edge.
					maxMoveY := (float32(nextCheckTileY * TileZoom)) - (p.pixelY + playerSize)
					p.pixelY += maxMoveY // Move only up to the collision point
					didMove = true
				} else { // Moving Up
					// Calculate max allowed movement: up to the bottom edge of the wall tile
					// Entity'p top edge (p.pixelY) aligns with wall'p bottom edge.
					maxMoveY := (float32((nextCheckTileY + 1) * TileZoom)) - p.pixelY
					p.pixelY += maxMoveY // Move only up to the collision point
					didMove = true
				}
				p.dir.Y = 0 // Stop vertical movement
			} else {
				// No collision vertically, apply full vertical movement
				p.pixelY += moveDistanceY
				didMove = true
			}
		}
		// p.tileX and p.tileY are updated at the top of the function based on current pixel position,
		// so no need to update them here again after pixel movement.
	}

	// === END REFINED MOVEMENT AND COLLISION HANDLING ===

	if didMove && p.tileX >= 0 && p.tileX < GameWidth {
		tile := board[p.tileY][p.tileX]
		if tile == Dot {
			board[p.tileY][p.tileX] = Empty
			p.slowTimer = SlowTime
			return 10
		} else if tile == Power {
			board[p.tileY][p.tileX] = Empty
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

	nextTileX := p.tileX + int(dir.X)
	nextTileY := p.tileY + int(dir.Y)

	// Check boundary conditions for the board
	if nextTileX < 0 || nextTileX >= GameWidth || nextTileY < 0 || nextTileY >= GameHeight {
		return false // Cannot move out of bounds
	}

	// Check if the next tile is a wall
	return board[nextTileY][nextTileX] != Wall
}
