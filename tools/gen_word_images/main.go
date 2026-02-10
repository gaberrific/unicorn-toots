package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
)

var (
	black   = color.RGBA{0, 0, 0, 255}
	white   = color.RGBA{255, 255, 255, 255}
	orange  = color.RGBA{255, 165, 0, 255}
	yellow  = color.RGBA{255, 255, 0, 255}
	brown   = color.RGBA{139, 69, 19, 255}
	green   = color.RGBA{34, 139, 34, 255}
	dkgreen = color.RGBA{0, 100, 0, 255}
	blue    = color.RGBA{100, 149, 237, 255}
	red     = color.RGBA{220, 20, 60, 255}
	pink    = color.RGBA{255, 182, 193, 255}
	gray    = color.RGBA{169, 169, 169, 255}
	ltblue  = color.RGBA{173, 216, 230, 255}
	dkbrown = color.RGBA{101, 67, 33, 255}
	tan     = color.RGBA{210, 180, 140, 255}
)

func newImg() *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, 32, 32))
}

func set(img *image.RGBA, x, y int, c color.RGBA) {
	if x >= 0 && x < 32 && y >= 0 && y < 32 {
		img.SetRGBA(x, y, c)
	}
}

func fillRect(img *image.RGBA, x0, y0, x1, y1 int, c color.RGBA) {
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			set(img, x, y, c)
		}
	}
}

func fillCircle(img *image.RGBA, cx, cy, r int, c color.RGBA) {
	for y := cy - r; y <= cy+r; y++ {
		for x := cx - r; x <= cx+r; x++ {
			dx, dy := x-cx, y-cy
			if dx*dx+dy*dy <= r*r {
				set(img, x, y, c)
			}
		}
	}
}

func drawCat() *image.RGBA {
	img := newImg()
	// Head
	fillCircle(img, 16, 16, 8, orange)
	// Ears
	fillRect(img, 9, 6, 11, 10, orange)
	fillRect(img, 20, 6, 22, 10, orange)
	// Inner ears
	set(img, 10, 8, pink)
	set(img, 21, 8, pink)
	// Eyes
	set(img, 13, 14, black)
	set(img, 19, 14, black)
	// Nose
	set(img, 16, 17, pink)
	// Mouth
	set(img, 15, 18, black)
	set(img, 17, 18, black)
	// Whiskers
	set(img, 10, 16, white)
	set(img, 11, 17, white)
	set(img, 22, 16, white)
	set(img, 21, 17, white)
	return img
}

func drawDog() *image.RGBA {
	img := newImg()
	// Head
	fillCircle(img, 16, 16, 8, brown)
	// Floppy ears
	fillRect(img, 7, 10, 9, 20, dkbrown)
	fillRect(img, 22, 10, 24, 20, dkbrown)
	// Muzzle
	fillCircle(img, 16, 19, 4, tan)
	// Eyes
	set(img, 13, 14, black)
	set(img, 19, 14, black)
	// Nose
	fillRect(img, 15, 17, 17, 18, black)
	// Tongue
	set(img, 16, 21, red)
	set(img, 16, 22, red)
	return img
}

func drawSun() *image.RGBA {
	img := newImg()
	// Sun body
	fillCircle(img, 16, 16, 7, yellow)
	// Rays
	for _, p := range [][2]int{
		{16, 4}, {16, 5}, {16, 27}, {16, 28},
		{4, 16}, {5, 16}, {27, 16}, {28, 16},
		{8, 8}, {9, 9}, {23, 8}, {22, 9},
		{8, 24}, {9, 23}, {23, 24}, {22, 23},
	} {
		set(img, p[0], p[1], yellow)
	}
	return img
}

func drawMoon() *image.RGBA {
	img := newImg()
	// Full circle
	fillCircle(img, 16, 16, 10, yellow)
	// Cut-out for crescent
	fillCircle(img, 20, 13, 8, color.RGBA{0, 0, 0, 0})
	return img
}

func drawStar() *image.RGBA {
	img := newImg()
	// Simple 5-pointed star shape
	c := yellow
	// Vertical center column
	for y := 4; y <= 10; y++ {
		set(img, 16, y, c)
	}
	// Cross arms
	fillRect(img, 8, 12, 24, 14, c)
	// Lower V
	for i := 0; i < 8; i++ {
		set(img, 10+i, 16+i/2, c)
		set(img, 22-i, 16+i/2, c)
	}
	// Fill center
	fillRect(img, 13, 8, 19, 16, c)
	fillRect(img, 11, 10, 21, 14, c)
	fillRect(img, 12, 16, 20, 18, c)
	fillRect(img, 13, 18, 19, 20, c)
	fillRect(img, 14, 20, 18, 22, c)
	return img
}

func drawFish() *image.RGBA {
	img := newImg()
	// Body (oval)
	fillCircle(img, 16, 16, 7, blue)
	fillRect(img, 10, 12, 22, 20, blue)
	// Tail
	fillRect(img, 4, 12, 9, 20, blue)
	fillRect(img, 2, 10, 6, 14, blue)
	fillRect(img, 2, 18, 6, 22, blue)
	// Eye
	set(img, 21, 14, white)
	set(img, 22, 14, black)
	// Mouth
	set(img, 24, 16, black)
	return img
}

func drawTree() *image.RGBA {
	img := newImg()
	// Trunk
	fillRect(img, 14, 22, 17, 30, brown)
	// Canopy
	fillCircle(img, 16, 14, 9, green)
	fillCircle(img, 12, 16, 6, green)
	fillCircle(img, 20, 16, 6, green)
	fillCircle(img, 16, 10, 6, dkgreen)
	return img
}

func drawFrog() *image.RGBA {
	img := newImg()
	// Body
	fillCircle(img, 16, 18, 8, green)
	// Head top
	fillCircle(img, 16, 12, 7, green)
	// Eyes (bulging)
	fillCircle(img, 11, 8, 3, green)
	fillCircle(img, 21, 8, 3, green)
	set(img, 11, 8, black)
	set(img, 21, 8, black)
	// Mouth
	fillRect(img, 12, 17, 20, 17, black)
	return img
}

func drawBird() *image.RGBA {
	img := newImg()
	// Body
	fillCircle(img, 16, 18, 6, blue)
	// Head
	fillCircle(img, 20, 12, 4, blue)
	// Eye
	set(img, 21, 11, black)
	// Beak
	set(img, 25, 12, orange)
	set(img, 26, 12, orange)
	// Wing
	fillRect(img, 10, 14, 16, 17, ltblue)
	// Tail
	set(img, 8, 18, blue)
	set(img, 7, 17, blue)
	set(img, 7, 19, blue)
	return img
}

func drawCake() *image.RGBA {
	img := newImg()
	// Base layer
	fillRect(img, 8, 18, 24, 28, pink)
	// Top layer
	fillRect(img, 10, 14, 22, 18, red)
	// Frosting drip
	fillRect(img, 8, 17, 24, 18, white)
	// Candle
	fillRect(img, 15, 8, 17, 14, yellow)
	// Flame
	set(img, 16, 6, orange)
	set(img, 16, 7, orange)
	return img
}

func drawHat() *image.RGBA {
	img := newImg()
	// Brim
	fillRect(img, 4, 22, 28, 24, black)
	// Crown
	fillRect(img, 10, 10, 22, 22, black)
	// Band
	fillRect(img, 10, 19, 22, 20, red)
	return img
}

func drawRun() *image.RGBA {
	img := newImg()
	// Stick figure running
	// Head
	fillCircle(img, 18, 6, 3, white)
	// Body
	fillRect(img, 17, 9, 19, 18, white)
	// Arms (forward/back)
	fillRect(img, 12, 11, 17, 12, white)
	fillRect(img, 19, 13, 24, 14, white)
	// Legs
	fillRect(img, 13, 19, 17, 20, white)
	fillRect(img, 14, 21, 16, 26, white)
	fillRect(img, 19, 19, 23, 20, white)
	fillRect(img, 22, 21, 24, 26, white)
	return img
}

func drawJump() *image.RGBA {
	img := newImg()
	// Stick figure jumping
	// Head
	fillCircle(img, 16, 4, 3, white)
	// Body
	fillRect(img, 15, 7, 17, 16, white)
	// Arms up
	fillRect(img, 10, 6, 15, 7, white)
	fillRect(img, 17, 6, 22, 7, white)
	fillRect(img, 9, 3, 10, 7, white)
	fillRect(img, 22, 3, 23, 7, white)
	// Legs spread
	fillRect(img, 11, 17, 15, 18, white)
	fillRect(img, 10, 19, 12, 22, white)
	fillRect(img, 17, 17, 21, 18, white)
	fillRect(img, 20, 19, 22, 22, white)
	// Ground indicator
	fillRect(img, 6, 28, 26, 29, gray)
	return img
}

func drawPlay() *image.RGBA {
	img := newImg()
	// Play button triangle
	c := white
	for row := 0; row < 20; row++ {
		width := row
		if row > 10 {
			width = 20 - row
		}
		for x := 0; x < width; x++ {
			set(img, 10+x, 6+row, c)
		}
	}
	// Ball
	fillCircle(img, 24, 24, 4, red)
	return img
}

func drawRain() *image.RGBA {
	img := newImg()
	// Cloud
	fillCircle(img, 12, 8, 5, gray)
	fillCircle(img, 20, 8, 5, gray)
	fillCircle(img, 16, 6, 5, gray)
	fillRect(img, 8, 8, 24, 12, gray)
	// Raindrops
	for _, p := range [][2]int{
		{10, 16}, {10, 17},
		{15, 18}, {15, 19},
		{20, 15}, {20, 16},
		{12, 22}, {12, 23},
		{18, 24}, {18, 25},
		{22, 21}, {22, 22},
	} {
		set(img, p[0], p[1], blue)
	}
	return img
}

func drawSnow() *image.RGBA {
	img := newImg()
	// Snowflake pattern - central
	c := white
	// Vertical
	fillRect(img, 15, 4, 16, 28, c)
	// Horizontal
	fillRect(img, 4, 15, 28, 16, c)
	// Diagonals
	for i := 0; i < 12; i++ {
		set(img, 4+i, 4+i, c)
		set(img, 5+i, 4+i, c)
		set(img, 27-i, 4+i, c)
		set(img, 26-i, 4+i, c)
	}
	return img
}

func drawLeaf() *image.RGBA {
	img := newImg()
	// Leaf shape
	fillCircle(img, 16, 14, 8, green)
	fillCircle(img, 18, 12, 6, green)
	// Tip
	set(img, 24, 6, green)
	set(img, 23, 7, green)
	set(img, 22, 8, green)
	// Stem
	fillRect(img, 10, 22, 11, 28, brown)
	fillRect(img, 12, 20, 13, 23, brown)
	// Vein
	for i := 0; i < 8; i++ {
		set(img, 14+i, 12+i/2, dkgreen)
	}
	return img
}

func drawBear() *image.RGBA {
	img := newImg()
	// Head
	fillCircle(img, 16, 16, 9, brown)
	// Ears
	fillCircle(img, 8, 8, 3, brown)
	fillCircle(img, 24, 8, 3, brown)
	fillCircle(img, 8, 8, 1, tan)
	fillCircle(img, 24, 8, 1, tan)
	// Muzzle
	fillCircle(img, 16, 19, 4, tan)
	// Eyes
	set(img, 12, 14, black)
	set(img, 20, 14, black)
	// Nose
	fillRect(img, 15, 17, 17, 18, black)
	return img
}

func drawDuck() *image.RGBA {
	img := newImg()
	// Body
	fillCircle(img, 14, 20, 8, yellow)
	// Head
	fillCircle(img, 22, 12, 5, yellow)
	// Eye
	set(img, 24, 11, black)
	// Beak
	fillRect(img, 27, 13, 30, 14, orange)
	// Wing
	fillCircle(img, 12, 18, 4, color.RGBA{255, 220, 0, 255})
	// Water line
	fillRect(img, 2, 26, 30, 27, blue)
	return img
}

func drawShip() *image.RGBA {
	img := newImg()
	// Hull
	fillRect(img, 4, 20, 28, 26, brown)
	// Make hull bottom narrower
	fillRect(img, 6, 26, 26, 28, dkbrown)
	// Cabin
	fillRect(img, 12, 14, 20, 20, white)
	// Window
	fillRect(img, 14, 16, 16, 18, blue)
	// Mast
	fillRect(img, 15, 4, 16, 14, black)
	// Flag
	fillRect(img, 17, 4, 22, 8, red)
	// Water
	fillRect(img, 0, 28, 31, 30, blue)
	return img
}

func save(img *image.RGBA, dir, name string) {
	path := filepath.Join(dir, name+".png")
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
	fmt.Println("Created", path)
}

func main() {
	dir := "assets/words"
	if err := os.MkdirAll(dir, 0o755); err != nil {
		panic(err)
	}

	generators := map[string]func() *image.RGBA{
		"cat":  drawCat,
		"dog":  drawDog,
		"sun":  drawSun,
		"moon": drawMoon,
		"star": drawStar,
		"fish": drawFish,
		"tree": drawTree,
		"frog": drawFrog,
		"bird": drawBird,
		"cake": drawCake,
		"hat":  drawHat,
		"run":  drawRun,
		"jump": drawJump,
		"play": drawPlay,
		"rain": drawRain,
		"snow": drawSnow,
		"leaf": drawLeaf,
		"bear": drawBear,
		"duck": drawDuck,
		"ship": drawShip,
	}

	for name, gen := range generators {
		save(gen(), dir, name)
	}
	fmt.Printf("Generated %d word images in %s/\n", len(generators), dir)
}
