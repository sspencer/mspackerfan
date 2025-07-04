package main

import (
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (g *Game) drawText(text string, x, y, pixelOffset int, color rl.Color) {

	str := []byte(strings.ToUpper(text))
	var fx, fy int

	for i, c := range str {
		ch := int(c)
		if ch >= 48 && ch <= 57 {
			// 0 - 9
			fx = ch - 48
			fy = 2
		} else if ch >= 65 && ch <= 79 {
			// A - O
			fx = ch - 65
			fy = 0
		} else if ch >= 80 && ch <= 90 {
			// P - Z
			fx = ch - 80
			fy = 1
		} else if ch == 33 {
			// "!"
			fx = 11
			fy = 1
		} else if ch == 47 {
			// "/"
			fx = 10
			fy = 2
		} else if ch == 45 {
			// "-"
			fx = 11
			fy = 2
		} else if ch == 34 {
			// '"'
			fx = 12
			fy = 2
		} else {
			// space or unmapped char
			continue
		}

		// ! = 33

		src := rl.NewRectangle(float32(fx*TileSize), float32(fy*TileSize), TileSize, TileSize)
		dst := rl.NewRectangle(float32((x+i)*TileSize*Zoom), float32(y*TileSize)*Zoom+float32(pixelOffset), TileSize*Zoom, TileSize*Zoom)
		rl.DrawTexturePro(g.font, src, dst, rl.Vector2{}, 0, color)

	}
}
