package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"github.com/gopxl/pixel/v2/ext/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

// Perlin noise implementation
var permutation = [512]int{}
var debug = true

func initPerlin() {
	p := [256]int{
		151, 160, 137, 91, 90, 15, 131, 13, 201, 95, 96, 53, 194, 233, 7, 225,
		140, 36, 103, 30, 69, 142, 8, 99, 37, 240, 21, 10, 23, 190, 6, 148,
		247, 120, 234, 75, 0, 26, 197, 62, 94, 252, 219, 203, 117, 35, 11, 32,
		57, 177, 33, 88, 237, 149, 56, 87, 174, 20, 125, 136, 171, 168, 68, 175,
		74, 165, 71, 134, 139, 48, 27, 166, 77, 146, 158, 231, 83, 111, 229, 122,
		60, 211, 133, 230, 220, 105, 92, 41, 55, 46, 245, 40, 244, 102, 143, 54,
		65, 25, 63, 161, 1, 216, 80, 73, 209, 76, 132, 187, 208, 89, 18, 169,
		200, 196, 135, 130, 116, 188, 159, 86, 164, 100, 109, 198, 173, 186, 3, 64,
		52, 217, 226, 250, 124, 123, 5, 202, 38, 147, 118, 126, 255, 82, 85, 212,
		207, 206, 59, 227, 47, 16, 58, 17, 182, 189, 28, 42, 223, 183, 170, 213,
		119, 248, 152, 2, 44, 154, 163, 70, 221, 153, 101, 155, 167, 43, 172, 9,
		129, 22, 39, 253, 19, 98, 108, 110, 79, 113, 224, 232, 178, 185, 112, 104,
		218, 246, 97, 228, 251, 34, 242, 193, 238, 210, 144, 12, 191, 179, 162, 241,
		81, 51, 145, 235, 249, 14, 239, 107, 49, 192, 214, 31, 181, 199, 106, 157,
		184, 84, 204, 176, 115, 121, 50, 45, 127, 4, 150, 254, 138, 236, 205, 93,
		222, 114, 67, 29, 24, 72, 243, 141, 128, 195, 78, 66, 215, 61, 156, 180,
	}
	for i := 0; i < 256; i++ {
		permutation[i] = p[i]
		permutation[256+i] = p[i]
	}
}

func fade(t float64) float64 {
	return t * t * t * (t*(t*6-15) + 10)
}

func lerp(t, a, b float64) float64 {
	return a + t*(b-a)
}

func grad(hash int, x, y float64) float64 {
	h := hash & 3
	switch h {
	case 0:
		return x + y
	case 1:
		return -x + y
	case 2:
		return x - y
	default:
		return -x - y
	}
}

func perlin2D(x, y float64) float64 {
	X := int(math.Floor(x)) & 255
	Y := int(math.Floor(y)) & 255

	x -= math.Floor(x)
	y -= math.Floor(y)

	u := fade(x)
	v := fade(y)

	A := permutation[X] + Y
	AA := permutation[A]
	AB := permutation[A+1]
	B := permutation[X+1] + Y
	BA := permutation[B]
	BB := permutation[B+1]

	return lerp(v,
		lerp(u, grad(permutation[AA], x, y), grad(permutation[BA], x-1, y)),
		lerp(u, grad(permutation[AB], x, y-1), grad(permutation[BB], x-1, y-1)),
	)
}

// Fractal Brownian Motion for richer noise
func fbm(x, y float64, octaves int) float64 {
	value := 0.0
	amplitude := 1.0 // amplitude between 0.5 and 1.5
	frequency := 1.0
	maxValue := 0.0

	for i := 0; i < octaves; i++ {
		value += amplitude * perlin2D(x*frequency, y*frequency)
		maxValue += amplitude
		amplitude *= 0.5
		frequency *= 2
	}

	return value / maxValue
}

const (
	bgScale    = 8    // pixels per noise sample (lower = more detail but slower)
	noiseScale = 0.02 // scale of the noise pattern
)

type Background struct {
	canvas *opengl.Canvas
	img    *image.RGBA
	width  int
	height int
	dirX   float64
	dirY   float64
}

func newBackground(w, h int) *Background {
	initPerlin()
	canvas := opengl.NewCanvas(pixel.R(0, 0, float64(w), float64(h)))
	img := image.NewRGBA(image.Rect(0, 0, w/bgScale, h/bgScale))

	// Random direction for animation
	randoX := -1.0 + rand.Float64()*(1.0-(-1.0)) // random float between -1 and 1
	randoY := -1.0 + rand.Float64()*(1.0-(-1.0)) // random float between -1 and 1
	angle := (randoX + randoY) * math.Pi
	speed := 0.2 + rand.Float64()*0.2 // speed between 0.2 and 0.4

	return &Background{
		canvas: canvas,
		img:    img,
		width:  w,
		height: h,
		dirX:   math.Cos(angle) * speed,
		dirY:   math.Sin(angle) * speed,
	}
}

func (bg *Background) update(t float64) {
	// if int(t*10)%2 == 0 {
	// 	// Random direction for animation
	// 	randoX := -1.0 + rand.Float64()*(1.0-(-1.0)) // random float between -1 and 1
	// 	randoY := -1.0 + rand.Float64()*(1.0-(-1.0)) // random float between -1 and 1
	// 	angle := (randoX + randoY) * math.Pi
	// 	speed := 0.2 + rand.Float64()*0.2 // speed between 0.2 and 0.4

	// 	bg.dirX = math.Cos(angle) * speed
	// 	bg.dirY = math.Sin(angle) * speed
	// }

	w := bg.width / bgScale
	h := bg.height / bgScale

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			nx := float64(x) * noiseScale
			ny := float64(y) * noiseScale

			// Animate by moving through noise space using random direction
			n := fbm(nx+t*bg.dirX, ny+t*bg.dirY, 4)

			// Map noise from [-1, 1] to [0, 1]
			n = (n + 1) / 2

			// Create a purple-ish gradient based on noise
			r := uint8(40 + n*60)
			g := uint8(20 + n*40)
			b := uint8(80 + n*100)

			bg.img.SetRGBA(x, y, color.RGBA{r, g, b, 255})
		}
	}

	// Convert to pixel picture and draw to canvas
	pic := pixel.PictureDataFromImage(bg.img)
	sprite := pixel.NewSprite(pic, pic.Bounds())

	bg.canvas.Clear(colornames.Black)
	sprite.Draw(bg.canvas, pixel.IM.Scaled(pixel.ZV, float64(bgScale)).Moved(pixel.V(float64(bg.width)/2, float64(bg.height)/2)))
}

const (
	winWidth  = 800
	winHeight = 600
	shapeSize = 64  // 32px sprite scaled 2x
	moveSpeed = 200 // pixels per second
)

type gameMode int

const (
	modeMenu gameMode = iota
	modeSpelling
	modeGem
)

type gameState int

const (
	statePlaying gameState = iota
	stateTryAgain
	stateWordComplete
)

type Letter struct {
	char      rune
	pos       pixel.Vec
	collected bool
}

type Gem struct {
	pos       pixel.Vec
	collected bool
}

func loadSpritesheet(path string) (*pixel.PictureData, []pixel.Rect) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	pic := pixel.PictureDataFromImage(img)
	bounds := pic.Bounds()
	numberOfFrames := bounds.W() / 32

	frames := []pixel.Rect{}
	for i := 0; i < int(numberOfFrames); i++ {
		frames = append(frames, pixel.R(float64(i*32), 0, float64((i+1)*32), 32))
	}

	return pic, frames
}

func loadSprite(path string) *pixel.Sprite {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	pic := pixel.PictureDataFromImage(img)
	return pixel.NewSprite(pic, pic.Bounds())
}

func loadWords(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var words []string
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		w := strings.TrimSpace(line)
		if w != "" {
			words = append(words, strings.ToUpper(w))
		}
	}
	return words
}

func loadWordImages(dir string) map[string]*pixel.Sprite {
	images := make(map[string]*pixel.Sprite)
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Warning: could not load word images:", err)
		return images
	}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".png" {
			continue
		}
		name := strings.TrimSuffix(entry.Name(), ".png")
		f, err := os.Open(filepath.Join(dir, entry.Name()))
		if err != nil {
			continue
		}
		img, _, err := image.Decode(f)
		f.Close()
		if err != nil {
			continue
		}
		pic := pixel.PictureDataFromImage(img)
		spr := pixel.NewSprite(pic, pic.Bounds())
		images[strings.ToUpper(name)] = spr
	}
	return images
}

const letterSize = 30.0 // approximate collision radius for a letter
const gemSize = 30.0    // collision radius for a gem

func randomLetterPositions(word string) []Letter {
	letters := make([]Letter, utf8.RuneCountInString(word))
	margin := 60.0
	minDist := 70.0

	i := 0
	for _, ch := range word {
		var pos pixel.Vec
		for attempts := 0; attempts < 100; attempts++ {
			pos = pixel.V(
				margin+rand.Float64()*(winWidth-2*margin),
				margin+rand.Float64()*(winHeight-2*margin-60), // leave room for HUD at top
			)
			ok := true
			for j := 0; j < i; j++ {
				if pos.Sub(letters[j].pos).Len() < minDist {
					ok = false
					break
				}
			}
			if ok {
				break
			}
		}
		letters[i] = Letter{char: ch, pos: pos, collected: false}
		i++
	}
	return letters
}

func randomGemPositions(count int) []Gem {
	gems := make([]Gem, count)
	margin := 60.0
	minDist := 70.0

	for i := 0; i < count; i++ {
		var pos pixel.Vec
		for attempts := 0; attempts < 100; attempts++ {
			pos = pixel.V(
				margin+rand.Float64()*(winWidth-2*margin),
				margin+rand.Float64()*(winHeight-2*margin-60),
			)
			ok := true
			for j := 0; j < i; j++ {
				if pos.Sub(gems[j].pos).Len() < minDist {
					ok = false
					break
				}
			}
			if ok {
				break
			}
		}
		gems[i] = Gem{pos: pos, collected: false}
	}
	return gems
}

func hsvToRGB(h, s, v float64) color.RGBA {
	h = math.Mod(h, 360)
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := v - c

	var r1, g1, b1 float64
	switch {
	case h < 60:
		r1, g1, b1 = c, x, 0
	case h < 120:
		r1, g1, b1 = x, c, 0
	case h < 180:
		r1, g1, b1 = 0, c, x
	case h < 240:
		r1, g1, b1 = 0, x, c
	case h < 300:
		r1, g1, b1 = x, 0, c
	default:
		r1, g1, b1 = c, 0, x
	}
	return color.RGBA{
		R: uint8((r1 + m) * 255),
		G: uint8((g1 + m) * 255),
		B: uint8((b1 + m) * 255),
		A: 255,
	}
}

func run() {
	cfg := opengl.WindowConfig{
		Title:  "Unicorn Toots",
		Bounds: pixel.R(0, 0, winWidth, winHeight),
		VSync:  true,
	}
	win, err := opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Load unicorn
	pic, frames := loadSpritesheet("assets/unicorn-v2.png")
	sprite := pixel.NewSprite(pic, frames[0])

	// Load gem sprite
	gemSprite := loadSprite("assets/gem.png")

	// Load words
	words := loadWords("assets/words.txt")

	// Load word images
	wordImages := loadWordImages("assets/words")

	// Text atlas
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	// Mode & state
	mode := modeMenu
	state := statePlaying
	stateTimer := 0.0

	// Spelling mode state
	pickWord := func() (string, []Letter) {
		w := words[rand.Intn(len(words))]
		return w, randomLetterPositions(w)
	}

	currentWord := ""
	var letters []Letter
	nextLetterIdx := 0

	// Gem mode state
	const gemsPerBatch = 5
	var gems []Gem
	gemScore := 0

	pos := pixel.V(winWidth/2, winHeight/2)
	last := time.Now()
	frameTime := 0.0
	frameIdx := 0
	const frameDuration = 0.15

	imd := imdraw.New(nil)
	hue := 0.0
	noiseTime := 0.0

	// Create animated background
	bg := newBackground(winWidth, winHeight)

	// Menu button rects
	spellingBtnRect := pixel.R(winWidth/2-150, winHeight/2-10, winWidth/2+150, winHeight/2+50)
	gemBtnRect := pixel.R(winWidth/2-150, winHeight/2-80, winWidth/2+150, winHeight/2-20)

	startSpellingMode := func() {
		mode = modeSpelling
		state = statePlaying
		currentWord, letters = pickWord()
		nextLetterIdx = 0
		pos = pixel.V(winWidth/2, winHeight/2)
	}

	startGemMode := func() {
		mode = modeGem
		state = statePlaying
		gems = randomGemPositions(gemsPerBatch)
		gemScore = 0
		pos = pixel.V(winWidth/2, winHeight/2)
	}

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		switch mode {
		case modeMenu:
			// Check for button clicks
			if win.JustPressed(pixel.MouseButtonLeft) {
				mpos := win.MousePosition()
				if spellingBtnRect.Contains(mpos) {
					startSpellingMode()
				} else if gemBtnRect.Contains(mpos) {
					startGemMode()
				}
			}

			// Draw menu with animated background
			noiseTime += dt
			bg.update(noiseTime)
			bg.canvas.Draw(win, pixel.IM.Moved(pixel.V(winWidth/2, winHeight/2)))

			// Title
			titleTxt := text.New(pixel.ZV, atlas)
			titleTxt.Color = colornames.Yellow
			titleTxt.WriteString("UNICORN TOOTS")
			tb := titleTxt.Bounds()
			titleCenter := pixel.V(winWidth/2, winHeight/2+120).Sub(tb.Center().Scaled(4))
			titleTxt.Draw(win, pixel.IM.Scaled(pixel.ZV, 4).Moved(titleCenter))

			// Spelling Mode button
			imd.Clear()
			imd.Color = colornames.Darkgreen
			imd.Push(spellingBtnRect.Min, spellingBtnRect.Max)
			imd.Rectangle(0)
			imd.Draw(win)

			sTxt := text.New(pixel.ZV, atlas)
			sTxt.Color = colornames.White
			sTxt.WriteString("SPELLING MODE")
			sb := sTxt.Bounds()
			sCenter := spellingBtnRect.Center().Sub(sb.Center().Scaled(2))
			sTxt.Draw(win, pixel.IM.Scaled(pixel.ZV, 2).Moved(sCenter))

			// Gem Mode button
			imd.Clear()
			imd.Color = colornames.Darkblue
			imd.Push(gemBtnRect.Min, gemBtnRect.Max)
			imd.Rectangle(0)
			imd.Draw(win)

			gTxt := text.New(pixel.ZV, atlas)
			gTxt.Color = colornames.White
			gTxt.WriteString("GEM MODE")
			gb := gTxt.Bounds()
			gCenter := gemBtnRect.Center().Sub(gb.Center().Scaled(2))
			gTxt.Draw(win, pixel.IM.Scaled(pixel.ZV, 2).Moved(gCenter))

			if debug {
				debugText := text.New(pixel.V(10, 30), atlas)
				debugText.Color = colornames.White
				fmt.Fprintf(debugText, "DirX: %.2f, DirY: %.2f", bg.dirX, bg.dirY)
				debugText.Draw(win, pixel.IM.Scaled(debugText.Orig, 2))
			}

			win.Update()
			continue

		case modeSpelling, modeGem:
			// shared movement and animation below
		}

		// Handle input
		delta := pixel.ZV
		if win.Pressed(pixel.KeyLeft) || win.Pressed(pixel.KeyA) {
			delta.X -= moveSpeed * dt
		}
		if win.Pressed(pixel.KeyRight) || win.Pressed(pixel.KeyD) {
			delta.X += moveSpeed * dt
		}
		if win.Pressed(pixel.KeyDown) || win.Pressed(pixel.KeyS) {
			delta.Y -= moveSpeed * dt
		}
		if win.Pressed(pixel.KeyUp) || win.Pressed(pixel.KeyW) {
			delta.Y += moveSpeed * dt
		}
		pos = pos.Add(delta)

		// Clamp to window bounds
		half := shapeSize / 2.0
		if pos.X < half {
			pos.X = half
		}
		if pos.X > winWidth-half {
			pos.X = winWidth - half
		}
		if pos.Y < half {
			pos.Y = half
		}
		if pos.Y > winHeight-half {
			pos.Y = winHeight - half
		}

		// Advance animation frame
		frameTime += dt
		if frameTime >= frameDuration {
			frameTime -= frameDuration
			frameIdx = (frameIdx + 1) % len(frames)
			sprite.Set(pic, frames[frameIdx])
		}

		// Back to menu with Escape
		if win.JustPressed(pixel.KeyEscape) {
			mode = modeMenu
			win.Update()
			continue
		}

		// Game logic per mode
		if mode == modeSpelling {
			switch state {
			case statePlaying:
				unicornRect := pixel.R(pos.X-half, pos.Y-half, pos.X+half, pos.Y+half)
				for i := range letters {
					if letters[i].collected {
						continue
					}
					lh := letterSize / 2
					letterRect := pixel.R(
						letters[i].pos.X-lh, letters[i].pos.Y-lh,
						letters[i].pos.X+lh, letters[i].pos.Y+lh,
					)
					if unicornRect.Intersects(letterRect) {
						if i == nextLetterIdx {
							letters[i].collected = true
							nextLetterIdx++
							if nextLetterIdx >= len(letters) {
								state = stateWordComplete
								stateTimer = 0
							}
						} else {
							state = stateTryAgain
							stateTimer = 0
						}
						break
					}
				}

			case stateTryAgain:
				stateTimer += dt
				if stateTimer >= 2.0 {
					letters = randomLetterPositions(currentWord)
					nextLetterIdx = 0
					state = statePlaying
				}

			case stateWordComplete:
				stateTimer += dt
				hue = math.Mod(hue+dt*180, 360)
				if stateTimer >= 3.0 {
					currentWord, letters = pickWord()
					nextLetterIdx = 0
					state = statePlaying
					hue = 0
				}
			}
		}

		if mode == modeGem {
			unicornRect := pixel.R(pos.X-half, pos.Y-half, pos.X+half, pos.Y+half)
			for i := range gems {
				if gems[i].collected {
					continue
				}
				gh := gemSize / 2
				gemRect := pixel.R(
					gems[i].pos.X-gh, gems[i].pos.Y-gh,
					gems[i].pos.X+gh, gems[i].pos.Y+gh,
				)
				if unicornRect.Intersects(gemRect) {
					gems[i].collected = true
					gemScore++
				}
			}

			// Check if all gems collected, spawn new batch
			allCollected := true
			for _, g := range gems {
				if !g.collected {
					allCollected = false
					break
				}
			}
			if allCollected {
				gems = randomGemPositions(gemsPerBatch)
			}
		}

		// Draw
		noiseTime += dt
		bg.update(noiseTime)
		bg.canvas.Draw(win, pixel.IM.Moved(pixel.V(winWidth/2, winHeight/2)))

		if mode == modeSpelling {
			// Draw letters on field
			if state == statePlaying || state == stateTryAgain {
				for _, l := range letters {
					if l.collected {
						continue
					}
					txt := text.New(pixel.ZV, atlas)
					txt.Color = colornames.Yellow
					txt.WriteString(string(l.char))
					bounds := txt.Bounds()
					offset := bounds.Center().Scaled(-1)
					txt.Draw(win, pixel.IM.Scaled(pixel.ZV, 3).Moved(l.pos.Add(offset.Scaled(3))))
				}
			}

			// Draw unicorn
			sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 2).Moved(pos))

			// Draw HUD - spelling progress at top
			hudTxt := text.New(pixel.V(10, winHeight-30), atlas)
			hudTxt.Color = colornames.White
			hudTxt.WriteString("Spell: ")
			for i, ch := range currentWord {
				if i < nextLetterIdx {
					hudTxt.Color = colornames.Lime
					hudTxt.WriteRune(ch)
				} else {
					hudTxt.Color = colornames.White
					hudTxt.WriteRune('_')
				}
				hudTxt.WriteRune(' ')
			}
			hudTxt.Draw(win, pixel.IM.Scaled(hudTxt.Orig, 2))

			// Draw word image prompt
			if spr, ok := wordImages[currentWord]; ok {
				imgX := hudTxt.Orig.X + hudTxt.Bounds().W()*2 + 40
				imgY := float64(winHeight - 30)
				spr.Draw(win, pixel.IM.Scaled(pixel.ZV, 2).Moved(pixel.V(imgX, imgY)))
			}

			// Draw state overlays
			switch state {
			case stateTryAgain:
				imd.Clear()
				imd.Color = color.RGBA{0, 0, 0, 150}
				imd.Push(pixel.V(0, 0), pixel.V(winWidth, winHeight))
				imd.Rectangle(0)
				imd.Draw(win)

				tryTxt := text.New(pixel.ZV, atlas)
				tryTxt.Color = colornames.Red
				tryTxt.WriteString("Try Again!")
				bounds := tryTxt.Bounds()
				center := pixel.V(winWidth/2, winHeight/2).Sub(bounds.Center().Scaled(3))
				tryTxt.Draw(win, pixel.IM.Scaled(pixel.ZV, 3).Moved(center))

			case stateWordComplete:
				imd.Clear()
				imd.Color = color.RGBA{0, 0, 0, 150}
				imd.Push(pixel.V(0, 0), pixel.V(winWidth, winHeight))
				imd.Rectangle(0)
				imd.Draw(win)

				completeTxt := text.New(pixel.ZV, atlas)
				completeTxt.Color = hsvToRGB(hue, 1, 1)
				completeTxt.WriteString(currentWord)
				bounds := completeTxt.Bounds()
				center := pixel.V(winWidth/2, winHeight/2).Sub(bounds.Center().Scaled(4))
				completeTxt.Draw(win, pixel.IM.Scaled(pixel.ZV, 4).Moved(center))
			}
		}

		if mode == modeGem {
			// Draw gems
			for _, g := range gems {
				if g.collected {
					continue
				}
				gemSprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 2).Moved(g.pos))
			}

			// Draw unicorn
			sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 2).Moved(pos))

			// Draw HUD - gem count
			hudTxt := text.New(pixel.V(10, winHeight-30), atlas)
			hudTxt.Color = colornames.Yellow
			fmt.Fprintf(hudTxt, "Gems: %d", gemScore)
			hudTxt.Draw(win, pixel.IM.Scaled(hudTxt.Orig, 2))
		}
		win.Update()
	}

}

func main() {
	opengl.Run(run)
}
