package main

import (
	"fmt"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (g *Game) Draw() {
	rl.ClearBackground(rl.Black)

	rl.BeginMode2D(g.camera2)
	g.drawBoard()

	// Animate characters
	rl.BeginShaderMode(g.shader) // make background transparent

	g.drawGhosts() // draw player before behavior when player is eaten
	g.drawPlayer()

	rl.EndShaderMode()
	rl.EndMode2D()

	g.drawLayout()
}

func (g *Game) drawLayout() {
	y := 0
	pixelOffset := 8
	g.drawText("1UP", 3, y, pixelOffset, rl.White)
	g.drawText("HIGH SCORE", 9, y, pixelOffset, rl.White)
	g.drawText("2UP", 22, y, pixelOffset, rl.White)

	y += 1
	pixelOffset += 4
	// TODO player 1 vs 2
	g.drawText(fmt.Sprintf("%d", g.player.score), 3, y, pixelOffset, rl.White) // player 1 score
	g.drawText(fmt.Sprintf("%d", g.highScore), 13, y, pixelOffset, rl.White)   // high score
	g.drawText("0", 24, y, pixelOffset, rl.White)                              // player 2 score

	pixelOffset = 8
	if g.debug {
		bottom := int32(ScreenHeight * TileSize * Zoom)
		p := g.player
		msg := fmt.Sprintf("pos=(%d,%d) state=%s", p.tile.X, p.tile.Y, strings.ToUpper(g.ghosts[0].state.String()))
		rl.DrawText(msg, 5, bottom-50, 20, rl.Green)

	} else {
		g.drawText(fmt.Sprintf("%d/%d", g.player.tile.X, g.player.tile.Y), 2, 34, pixelOffset, rl.White) // player 2 score
		g.drawText(fmt.Sprintf("dots %d", g.player.dots), 19, 34, pixelOffset, rl.White)                 // player 2 score
	}

}

func (g *Game) drawBoard() {

	x := float32(GameWidth*TileSize) + TileSize/2
	y := float32(g.boardNum * GameHeight * TileSize)
	w := float32(GameWidth * TileSize)  // 28 * 8
	h := float32(GameHeight * TileSize) // 31 * 8
	src := rl.NewRectangle(x, y, w, h)
	dst := rl.NewRectangle(0, 0, w*Zoom, h*Zoom)

	rl.DrawTexturePro(g.texture, src, dst, rl.Vector2{}, 0, rl.White)

	// source location from original artwork texture of a
	// dot and power up so that the dots are the same
	// color as you see in that image
	var dotX, dotY, powerX, powerY float32 = 1, 1, 1, 2

	if g.boardNum == 1 {
		dotX, dotY = 1, 36
		powerX, powerY = 1, 35
	} else if g.boardNum == 2 {
		dotX, dotY = 1, 63
		powerX, powerY = 1, 65
	} else if g.boardNum == 3 {
		dotX, dotY = 1, 94
		powerX, powerY = 1, 96
	} else if g.boardNum == 4 {
		dotX, dotY = 1, 125
		powerX, powerY = 1, 127
	} else if g.boardNum == 5 {
		dotX, dotY = 1, 157
		powerX, powerY = 1, 158
	}

	dot := rl.NewRectangle(dotX*TileSize, dotY*TileSize, TileSize, TileSize)
	power := rl.NewRectangle(powerX*TileSize, powerY*TileSize, TileSize, TileSize)

	for y := 0; y < GameHeight; y++ {
		for x := 0; x < GameWidth; x++ {
			if g.maze[y][x] == Dot {
				dst = rl.NewRectangle(float32(x*TileSize*Zoom), float32(y*TileSize*Zoom), TileSize*Zoom, TileSize*Zoom)
				rl.DrawTexturePro(g.texture, dot, dst, rl.Vector2{}, 0, rl.White)
			} else if g.maze[y][x] == Power {
				dst = rl.NewRectangle(float32(x*TileSize*Zoom), float32(y*TileSize*Zoom), TileSize*Zoom, TileSize*Zoom)
				rl.DrawTexturePro(g.texture, power, dst, rl.Vector2{}, 0, rl.White)
			}
		}
	}

	if g.debugLayout {
		g.drawCheckerBoard()
	}

}

func (g *Game) drawGhosts() {
	for _, e := range g.ghosts {
		loc := e.sprite[e.dir]
		sx := float32(loc.X) + float32(e.frame)*e.width
		sy := float32(loc.Y)
		src := rl.NewRectangle(sx, sy, e.width, e.height) // sprite

		offsetX, offsetY := float32(0), float32(0)
		if e.InHouse() {
			offsetX = float32(-7)
		}

		dst := rl.NewRectangle(e.pixel.X+offsetX*Zoom, e.pixel.Y+offsetY*Zoom, e.width*Zoom, e.height*Zoom)
		rl.DrawTexturePro(g.texture, src, dst, rl.Vector2{}, 0, rl.White)

		if g.debug && e.tile.X != 0 && e.tile.Y != 0 {
			f := TileSize * Zoom / 2
			target := e.target.Clamp()
			rl.DrawCircle(int32(target.X*TileSize*Zoom+f), int32(target.Y*TileSize*Zoom+f), TileSize*Zoom, rl.ColorAlpha(e.color, 0.5))
			v1 := rl.Vector2{X: float32(target.X*TileSize*Zoom + f), Y: float32(target.Y*TileSize*Zoom + f)}
			v2 := rl.Vector2{X: float32(e.tile.X*TileSize*Zoom + f), Y: float32(e.tile.Y*TileSize*Zoom + f)}
			rl.DrawLineEx(v1, v2, 4, e.color)
		}
	}
}

func (g *Game) drawPlayer() {
	// TODO move texture to entity and change receiver from Game to Entity
	s := g.player
	loc := s.sprite[s.dir]
	sx := float32(loc.X) + float32(s.frame)*s.width
	sy := float32(loc.Y)
	src := rl.NewRectangle(sx, sy, s.width, s.height) // sprite

	dst := rl.NewRectangle(s.pixel.X, s.pixel.Y, s.width*Zoom, s.height*Zoom)
	rl.DrawTexturePro(g.texture, src, dst, rl.Vector2{}, 0, rl.White)
}

func (g *Game) drawCheckerBoard() {
	i := 0
	c1 := rl.Color{255, 255, 255, 120}
	c2 := rl.Color{255, 255, 255, 80}
	for y := 0; y < GameHeight; y++ {
		for x := 0; x < GameWidth; x++ {
			if y%10 == 0 {
				g.drawText(fmt.Sprintf("%d", x%10), x, y, 0, rl.White)
			} else if x == 0 {
				g.drawText(fmt.Sprintf("%d", y%10), 0, y, 0, rl.White)

			}

			rec := rl.NewRectangle(float32(x*TileSize*Zoom), float32(y*TileSize*Zoom), TileSize*Zoom, TileSize*Zoom)
			if i%2 == 1 {
				rl.DrawRectangleRec(rec, c1)
			} else {
				rl.DrawRectangleRec(rec, c2)
			}
			i++

		}
		i++
	}
}

func chromaShader() rl.Shader {
	shader := rl.LoadShader("", "chroma_key.fs") // Empty string for vertex shader (use default)
	// Set shader uniforms
	keyColor := []float32{0.0, 0.0, 0.0} // Black (normalized RGB: 0.0 to 1.0)
	threshold := []float32{0.05}         // Tolerance for slight color variations

	keyColorLoc := rl.GetShaderLocation(shader, "keyColor")
	thresholdLoc := rl.GetShaderLocation(shader, "threshold")
	rl.SetShaderValue(shader, keyColorLoc, keyColor, rl.ShaderUniformVec3)
	rl.SetShaderValue(shader, thresholdLoc, threshold, rl.ShaderUniformFloat)
	return shader
}
