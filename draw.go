package main

import (
	"fmt"

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
	bottom := int32(ScreenHeight * Pixel)
	msg := fmt.Sprintf("state: %s, dots: %d, time: %0.1f", g.ghosts[0].state, g.dotsEaten, g.levelTime)
	rl.DrawText(msg, 5, bottom-50, 24, rl.Green)
	rl.DrawFPS(10, 10)
	//g.drawText(fmt.Sprintf("dots %d", g.dotsEaten), 19, 34, pixelOffset, rl.White) // player 2 score

}

func (g *Game) drawBoard() {

	x := float32(GameWidth*Size) + Size/2
	y := float32(g.boardNum * GameHeight * Size)
	w := float32(GameWidth * Size)  // 28 * 8
	h := float32(GameHeight * Size) // 31 * 8
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

	dot := rl.NewRectangle(dotX*Size, dotY*Size, Size, Size)
	power := rl.NewRectangle(powerX*Size, powerY*Size, Size, Size)

	for y := 0; y < GameHeight; y++ {
		for x := 0; x < GameWidth; x++ {
			tile := g.maze[y][x]
			if tile == Wall {
				continue
			}

			rec := rl.NewRectangle(float32(x*Pixel), float32(y*Pixel), Pixel, Pixel)

			if tile == Dot {
				rl.DrawTexturePro(g.texture, dot, rec, rl.Vector2{}, 0, rl.White)
			} else if tile == Power {
				rl.DrawTexturePro(g.texture, power, rec, rl.Vector2{}, 0, rl.White)
			} else if tile == Tunnel {
				rl.DrawRectangleRec(rec, rl.Gray)
			}
		}
	}

	if g.debugLayout {
		g.drawCheckerBoard()
	}

}

func (g *Game) drawGhosts() {
	for _, e := range g.ghosts {
		var loc Vec2i
		if e.state == Frightened {
			loc = e.fright[e.frightState]
		} else {
			loc = e.sprite[e.dir]
		}

		sx := float32(loc.X) + float32(e.frame)*e.width
		sy := float32(loc.Y)
		src := rl.NewRectangle(sx, sy, e.width, e.height) // sprite

		dst := rl.NewRectangle(e.pixel.X, e.pixel.Y, e.width*Zoom, e.height*Zoom)
		rl.DrawTexturePro(g.texture, src, dst, rl.Vector2{}, 0, rl.White)

		if g.debug && e.tile.X != 0 && e.tile.Y != 0 {
			f := float32(Pixel)
			f2 := f / 2
			target := e.target.Clamp()
			rl.DrawCircle(int32(target.X*Pixel)+int32(f2), int32(target.Y*Pixel)+int32(f2), Pixel, rl.ColorAlpha(e.color, 0.5))
			v1 := rl.Vector2{X: float32(target.X*Pixel) + f2, Y: float32(target.Y*Pixel) + f2}
			v2 := rl.Vector2{X: e.pixel.X + f, Y: e.pixel.Y + f}
			rl.DrawLineEx(v1, v2, 4, rl.ColorAlpha(e.color, 0.5))
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

			rec := rl.NewRectangle(float32(x*Pixel), float32(y*Pixel), Pixel, Pixel)
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
