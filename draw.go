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

	g.drawGhosts() // draw player before ghost when player is eaten
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
	g.drawText("0", 22, y, pixelOffset, rl.White)                              // player 2 score

	pixelOffset = 8
	g.drawText(fmt.Sprintf("%d/%d", g.player.tileX, g.player.tileY), 2, 34, pixelOffset, rl.White) // player 2 score
	g.drawText(fmt.Sprintf("dots %d", g.player.dots), 19, 34, pixelOffset, rl.White)               // player 2 score

}

func (g *Game) drawBoard() {

	x := float32(GameWidth*TileSize) + TileSize/2
	y := float32(g.boardNum * GameHeight * TileSize)
	w := float32(GameWidth * TileSize)  // 28 * 8
	h := float32(GameHeight * TileSize) // 31 * 8
	src := rl.NewRectangle(x, y, w, h)
	dst := rl.NewRectangle(0, 0, w*Zoom, h*Zoom)

	rl.DrawTexturePro(g.texture, src, dst, rl.Vector2{}, 0, rl.White)

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
			if g.board[y][x] == Dot {
				dst = rl.NewRectangle(float32(x*TileSize*Zoom), float32(y*TileSize*Zoom), TileSize*Zoom, TileSize*Zoom)
				rl.DrawTexturePro(g.texture, dot, dst, rl.Vector2{}, 0, rl.White)
			} else if g.board[y][x] == Power {
				dst = rl.NewRectangle(float32(x*TileSize*Zoom), float32(y*TileSize*Zoom), TileSize*Zoom, TileSize*Zoom)
				rl.DrawTexturePro(g.texture, power, dst, rl.Vector2{}, 0, rl.White)
			}
		}
	}

	if g.debug {
		g.drawCheckerBoard()
	}

}

func (g *Game) drawGhosts() {
	// TODO move texture to entity and change receiver from Game to Entity
	for _, s := range g.ghosts {
		loc := s.loc[s.shape]
		sx := float32(loc.x) + float32(s.frame)*s.width
		sy := float32(loc.y)
		src := rl.NewRectangle(sx, sy, s.width, s.height) // sprite

		offsetX, offsetY := float32(-4), float32(-4)
		if s.tileY == 14 && (s.tileX >= 12 && s.tileX <= 16) {
			offsetX = float32(-7)
		}

		dst := rl.NewRectangle(s.pixelX+offsetX*Zoom, s.pixelY+offsetY*Zoom, s.width*Zoom, s.height*Zoom)
		rl.DrawTexturePro(g.texture, src, dst, rl.Vector2{}, 0, rl.White)
	}
}

func (g *Game) drawPlayer() {
	// TODO move texture to entity and change receiver from Game to Entity
	s := g.player
	loc := s.loc[s.shape]
	sx := float32(loc.x) + float32(s.frame)*s.width
	sy := float32(loc.y)
	src := rl.NewRectangle(sx, sy, s.width, s.height) // sprite

	offsetX := float32(-4) //float32(-4) * Zoom
	offsetY := float32(-4) //float32(-4) * Zoom
	dst := rl.NewRectangle(s.pixelX+offsetX*Zoom, s.pixelY+offsetY*Zoom, s.width*Zoom, s.height*Zoom)
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
