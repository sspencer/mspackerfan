package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func (g *Game) Update() {
	g.levelTime = rl.GetTime() - g.startTime
	p := g.player
	moved := false
	if rl.IsKeyPressed(rl.KeyRight) {
		moved = true
		p.nextDir = Right
		p.nextVel = Vec2i{X: 1}
	}

	if rl.IsKeyPressed(rl.KeyLeft) {
		moved = true
		p.nextDir = Left
		p.nextVel = Vec2i{X: -1}
	}

	if rl.IsKeyPressed(rl.KeyUp) {
		moved = true
		p.nextDir = Up
		p.nextVel = Vec2i{Y: -1}
	}

	if rl.IsKeyPressed(rl.KeyDown) {
		moved = true
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

	if rl.IsKeyPressed(rl.KeyP) || rl.IsKeyPressed(rl.KeySpace) {
		moved = true
		g.paused = !g.paused
	}

	g.moved = moved

	if g.debug {
		if rl.IsKeyPressed(rl.KeyC) {
			g.setGhostMode(Chase)
		}

		if rl.IsKeyPressed(rl.KeyS) {
			g.setGhostMode(Scatter)
		}

		if rl.IsKeyPressed(rl.KeyF) {
			g.setGhostMode(Frightened)
			g.frightTime = rl.GetTime() + FrightDuration
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
		if mode == Frightened {
			g.frightTime = rl.GetTime() + FrightDuration
		}
	}
}
