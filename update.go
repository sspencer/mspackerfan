package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func (g *Game) Update() {
	g.levelTime = float32(rl.GetTime() - g.startTime)
	p := g.player
	if rl.IsKeyPressed(rl.KeyRight) {
		p.nextDir = Right
		p.nextVel = Vec2i{X: 1}
	}

	if rl.IsKeyPressed(rl.KeyLeft) {
		p.nextDir = Left
		p.nextVel = Vec2i{X: -1}
	}

	if rl.IsKeyPressed(rl.KeyUp) {
		p.nextDir = Up
		p.nextVel = Vec2i{Y: -1}
	}

	if rl.IsKeyPressed(rl.KeyDown) {
		p.nextDir = Down
		p.nextVel = Vec2i{Y: 1}
	}

	if rl.IsKeyPressed(rl.KeyD) {
		g.debug = !g.debug
	}

	if rl.IsKeyPressed(rl.KeyL) {
		g.debugLayout = !g.debugLayout
	}

	if rl.IsKeyPressed(rl.KeyN) {
		g.boardNum = (g.boardNum + 1) % 6
		g.mapBoard()
	}

	if rl.IsKeyPressed(rl.KeyP) {
		g.paused = !g.paused
	}

	if g.debug {
		if rl.IsKeyPressed(rl.KeyC) {
			g.setGhostMode(Chase)
		}

		if rl.IsKeyPressed(rl.KeyS) {
			g.setGhostMode(Scatter)
		}

		if rl.IsKeyPressed(rl.KeyF) {
			g.setGhostMode(Frightened)
		}
	}

	if g.paused {
		return
	}

	p.Update(g)
	for _, ghost := range g.ghosts {
		ghost.Update(g)
	}
}

func (g *Game) setGhostMode(mode GhostState) {
	for _, ghost := range g.ghosts {
		ghost.state = mode
	}
}
