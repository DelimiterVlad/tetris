// tetris.go
package main

import (
	//	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"time"
	//	"unsafe"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	/// Window is the global SDL window.
	///
	Window *sdl.Window

	/// Renderer is the global SDL renderer.
	///
	Renderer *sdl.Renderer

	/// Screen is the global SDL render target for the VM's video memory.
	///
	Screen *sdl.Texture

	Bmp *sdl.Surface

	rect sdl.Rect
	/// Font is a fixed-width, bitmap font.
	///
	Font *sdl.Texture

	Gopher *sdl.Texture

	goph_y *sdl.Texture

	ox     int = 8
	oy     int = 8
	w      int = 286
	h      int = 480
	countx int = 11
	county int = 20
	wx     int
	hy     int

	KeyMap = map[sdl.Scancode]uint{
		sdl.SCANCODE_X: 0x0,
		sdl.SCANCODE_1: 0x1,
		sdl.SCANCODE_2: 0x2,
		sdl.SCANCODE_3: 0x3,
		sdl.SCANCODE_Q: 0x4,
		sdl.SCANCODE_W: 0x5,
		sdl.SCANCODE_E: 0x6,
		sdl.SCANCODE_A: 0x7,
		sdl.SCANCODE_S: 0x8,
		sdl.SCANCODE_D: 0x9,
		sdl.SCANCODE_Z: 0xA,
		sdl.SCANCODE_C: 0xB,
		sdl.SCANCODE_4: 0xC,
		sdl.SCANCODE_R: 0xD,
		sdl.SCANCODE_F: 0xE,
		sdl.SCANCODE_V: 0xF,
	}
)

func minit() {
	runtime.LockOSThread()
}

func createWindow() {
	var err error

	// window attributes
	flags := sdl.WINDOW_OPENGL // | sdl.WINDOWPOS_CENTERED

	// create the window and renderer
	Window, Renderer, err = sdl.CreateWindowAndRenderer(300, 520, uint32(flags))
	if err != nil {
		panic(err)
	}

	// set the title
	Window.SetTitle("Tetris on GOlang!")

	// load the icon and use it if found
	//setIcon()

	// desired screen format and access
	format := sdl.PIXELFORMAT_RGB888
	access := sdl.TEXTUREACCESS_TARGET

	// create a render target for the display
	Screen, err = Renderer.CreateTexture(uint32(format), access, 128, 64)
	if err != nil {
		panic(err)
	}
}

/// loadFont loads the bitmap surface with font on it.
///
func loadFont() {
	var surface *sdl.Surface
	var err error

	if surface, err = sdl.LoadBMP("font.bmp"); err != nil {
		panic(err)
	}

	// get the magenta color
	mask := sdl.MapRGB(surface.Format, 255, 0, 255)

	// set the mask color key
	surface.SetColorKey(1, mask)

	// create the texture
	if Font, err = Renderer.CreateTextureFromSurface(surface); err != nil {
		panic(err)
	}

}

func loadGopher() {
	var usurface *sdl.Surface
	var err error

	if usurface, err = sdl.LoadBMP("gopher.bmp"); err != nil {
		panic(err)
	}

	// get the magenta color
	mask := sdl.MapRGB(usurface.Format, 255, 0, 255)

	// set the mask color key
	usurface.SetColorKey(1, mask)

	// create the texture
	if Gopher, err = Renderer.CreateTextureFromSurface(usurface); err != nil {
		panic(err)
	}

	if usurface, err = sdl.LoadBMP("y.bmp"); err != nil {
		panic(err)
	}
	// get the magenta color
	mask2 := sdl.MapRGB(usurface.Format, 255, 0, 255)
	usurface.SetColorKey(1, mask2)
	if goph_y, err = Renderer.CreateTextureFromSurface(usurface); err != nil {
		panic(err)
	}
}
func calculate_aim(offsetX int, wx int, offsetY int, hy int) (int, int, bool) {
	rez := false
	var cx, cy int
	cx = 0
	cy = 0
	for i := 0; i < 20; i++ {
		for j := 0; j < 11; j++ {
			if myfield.field[i][j] == 2 {
				rez = true
				cx = (offsetX + j*wx + 2)
				cy = (offsetY + i*hy + 2)
			}
		}
	}
	return cx, cy, rez
}
func draw_gopher(x, y int, offsetX int, wX int, offsetY int, hy int) {
	var new1_x, new1_y, new2_x, new2_y int
	var a_x, a_y int
	var rez bool
	src := sdl.Rect{
		W: 157,
		H: 196,
		X: 0,
		Y: 0,
	}
	dst := sdl.Rect{
		X: int32(x),
		Y: int32(y),
		W: 157,
		H: 196,
	}

	Renderer.Copy(Gopher, &src, &dst)

	coord1_x := 45
	coord1_y := 40
	coord2_x := 93
	coord2_y := 36
	src1 := sdl.Rect{
		W: 17,
		H: 16,
		X: 0,
		Y: 0,
	}
	a_x, a_y, rez = calculate_aim(offsetX, wX, offsetY, hy)
	if rez {
		m1x := x + coord1_x
		m1y := y + coord1_y
		m2x := x + coord2_x
		m2y := y + coord2_y
		d1x := a_x - m1x
		d1y := a_y - m1y
		d2x := a_x - m2x
		d2y := a_y - m2y
		if d1x != 0 && d1y != 0 && d2x != 0 && d2y != 0 {
			new1_x = m1x + int(d1x/20)
			new1_y = m1y + int(d1y/25)
			new2_x = m2x + int(d2x/17)
			new2_y = m2y + int(d2y/22)
			dst1 := sdl.Rect{
				X: int32(new1_x),
				Y: int32(new1_y),
				W: 17,
				H: 16,
			}
			Renderer.Copy(goph_y, &src1, &dst1)
			dst1.X = int32(new2_x)
			dst1.Y = int32(new2_y)
			Renderer.Copy(goph_y, &src1, &dst1)
		} else {
			dst1 := sdl.Rect{
				X: int32(x + coord1_x),
				Y: int32(y + coord1_y),
				W: 17,
				H: 16,
			}
			Renderer.Copy(goph_y, &src1, &dst1)
			dst1.X = int32(x + coord2_x)
			dst1.Y = int32(y + coord2_y)
			Renderer.Copy(goph_y, &src1, &dst1)
		}
	} else {
		dst1 := sdl.Rect{
			X: int32(x + coord1_x),
			Y: int32(y + coord1_y),
			W: 17,
			H: 16,
		}
		Renderer.Copy(goph_y, &src1, &dst1)
		dst1.X = int32(x + coord2_x)
		dst1.Y = int32(y + coord2_y)
		Renderer.Copy(goph_y, &src1, &dst1)
	}

}

/// drawText using the bitmap font a string at a given location.
///
func drawText(s string, x, y int) {
	src := sdl.Rect{W: 5, H: 7}
	dst := sdl.Rect{
		X: int32(x),
		Y: int32(y),
		W: 5,
		H: 7,
	}

	// loop over all the characters in the string
	for _, c := range s {
		if c > 32 && c < 127 {
			src.X = (c - 33) * 6

			// draw the character to the renderer
			Renderer.Copy(Font, &src, &dst)
		}

		// advance
		dst.X += 7
	}
}

type mcoord struct {
	x         int
	y         int
	direction int
}

var my_coord mcoord

/// clear the renderer, redraw everything, and present.
///
func redraw() {
	//updateScreen()

	// clear the renderer
	Renderer.SetDrawColor(32, 42, 53, 255)
	rect.X = 0
	rect.Y = 0
	rect.W = 300
	rect.H = 520
	Renderer.FillRect(&rect)

	// frame the screen, instructions, log, and registers
	frame(ox, oy, w, h)
	//frame(8, 208, 386, 164)
	//frame(402, 8, 204, 194)
	//frame(402, 208, 204, 164)
	drawText("<DELIMITER-Lab>", 190, 500)
	drawText("Smolentsev Vladimir", 30, 22)
	drawText("telegram: @DelimiterVlad", 20, 43)

	drawScreen(ox, oy, ox+w, oy+h)

	// show it
	Renderer.Present()
	Renderer.Clear()
	//Renderer.SetRenderTarget(nil)

}

type mfld struct {
	field     [20][11]int
	Timestamp time.Time
	lock      bool
	counter   int
}

var myfield mfld
var copyfield mfld

func initField() {
	for i := 0; i < 20; i++ {
		for j := 0; j < 11; j++ {
			myfield.field[i][j] = 0
			copyfield.field[i][j] = 0
		}
	}
	myfield.counter = 0
	myfield.lock = false
}

type smn struct {
	field [4][4]int8
}

var my_elements [10]smn

func generate(mfield *mfld) bool {
	var rez bool = true
	var i, j, l int
	l = rand.Intn(100) % 10
	for i = 0; i < 4 && rez; i++ {
		for j = 0; j < 4 && rez; j++ {
			if my_elements[l].field[i][j] != 0 {
				if mfield.field[i][j+3] == 0 {
					mfield.field[i][j+3] = int(my_elements[l].field[i][j])
				} else {
					rez = false
				}

			}
		}
	}
	return rez
}

func init_generic() {

	my_elements[0].field = [4][4]int8{{0, 0, 1, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}
	my_elements[1].field = [4][4]int8{{0, 1, 1, 0},
		{0, 1, 1, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}
	my_elements[2].field = [4][4]int8{{0, 1, 2, 1},
		{0, 0, 1, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}
	my_elements[3].field = [4][4]int8{{0, 1, 2, 1},
		{0, 1, 0, 1},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}
	my_elements[4].field = [4][4]int8{{1, 1, 2, 1},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}
	my_elements[5].field = [4][4]int8{{1, 1, 2, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}
	my_elements[6].field = [4][4]int8{{0, 2, 1, 1},
		{0, 1, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}
	my_elements[7].field = [4][4]int8{{0, 0, 0, 1},
		{0, 0, 2, 0},
		{0, 1, 0, 0},
		{0, 0, 0, 0}}
	my_elements[8].field = [4][4]int8{{0, 0, 1, 1},
		{0, 1, 2, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}
	my_elements[9].field = [4][4]int8{{1, 1, 0, 0},
		{0, 2, 1, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0}}

}

func check_full_rows(mfield *mfld) {
	var i, j, k, l, flag int

	for i = 19; i >= 0; i-- {
		flag = 0
		for j = 0; j < 11; j++ {
			flag += mfield.field[i][j]
		}
		if flag == -11 {
			mfield.counter += 11
			for k = i; k > 0; k-- {
				for l = 0; l < 11; l++ {
					mfield.field[k][l] = mfield.field[k-1][l]
				}
			}
			for l = 0; l < 11; l++ {
				mfield.field[0][l] = 0
			}
			i++
		}
	}
}

func fallStep(mfield *mfld) bool {
	var flag int = 0
	var i, j int
	for i = 19; i >= 0; i-- {
		for j = 0; j < 11; j++ {
			if mfield.field[i][j] == 1 || mfield.field[i][j] == 2 {
				if i == 19 {
					//zamorozka
					flag = 1
				} else {
					if mfield.field[i+1][j] < 0 {
						flag = 1
					}
				}
			}
		}
	}
	if flag == 1 {
		for i = 0; i < 20; i++ {
			for j = 0; j < 11; j++ {
				if mfield.field[i][j] == 1 || mfield.field[i][j] == 2 {
					mfield.field[i][j] = -1
				}
			}
		}
	} else {
		//padaem
		for i = 19; i >= 0; i-- {
			for j = 0; j < 11; j++ {
				if mfield.field[i][j] == 1 || mfield.field[i][j] == 2 {
					mfield.field[i+1][j] = mfield.field[i][j]
					mfield.field[i][j] = 0
				}
			}
		}
	}
	if flag == 1 {
		check_full_rows(mfield)
		return true
	} else {
		return false
	}

}

func stepLeft(mfield *mfld) {
	var i, j int
	var flag int = 0
	for i = 0; i < 11; i++ {
		for j = 0; j < 20; j++ {
			if mfield.field[j][i] == 1 || mfield.field[j][i] == 2 {
				if i == 0 {
					flag = 1
				} else {
					if mfield.field[j][i-1] < 0 {
						flag = 1
					}
				}
			}
		}
	}
	if flag != 1 {
		for i = 0; i < 11; i++ {
			for j = 0; j < 20; j++ {
				if mfield.field[j][i] == 1 || mfield.field[j][i] == 2 {
					mfield.field[j][i-1] = mfield.field[j][i]
					mfield.field[j][i] = 0
				}
			}
		}
	}
}
func stepRight(mfield *mfld) {
	var i, j int
	var flag int = 0
	for i = 10; i >= 0; i-- {
		for j = 0; j < 20; j++ {
			if mfield.field[j][i] == 1 || mfield.field[j][i] == 2 {
				if i == 10 {
					flag = 1
				} else {
					if mfield.field[j][i+1] < 0 {
						flag = 1
					}
				}
			}
		}
	}
	if flag != 1 {
		for i = 10; i >= 0; i-- {
			for j = 0; j < 20; j++ {
				if mfield.field[j][i] == 1 || mfield.field[j][i] == 2 {
					mfield.field[j][i+1] = mfield.field[j][i]
					mfield.field[j][i] = 0
				}
			}
		}
	}
}
func rotateL(mfield *mfld) {
	var i, j, k, l int
	var dX, dY int
	var new_dX, new_dY int
	var new_x, new_y int
	var flag int = 0
	for i = 0; i < 20; i++ {
		for j = 0; j < 11; j++ {
			copyfield.field[i][j] = 0
		}
	}
	for i = 0; i < 11; i++ {
		for j = 0; j < 20; j++ {
			if mfield.field[j][i] == 2 {
				k = i
				l = j
			}
		}
	}
	for i = 0; i < 11 && flag == 0; i++ {
		for j = 0; j < 20 && flag == 0; j++ {
			if mfield.field[j][i] == 1 {
				dX = k - i
				dY = l - j
				new_dY = dX
				new_dX = dY * (-1)
				new_y = l - new_dY
				new_x = k - new_dX
				if new_y >= 0 && new_y < 20 {
					if new_x >= 0 && new_x < 11 {
						if mfield.field[new_y][new_x] < 0 {
							flag = 1
						}

					} else {
						flag = 1
					}
				} else {
					flag = 1
				}

			}
		}
	}
	if flag != 1 {
		for i = 0; i < 11; i++ {
			for j = 0; j < 20; j++ {
				if mfield.field[j][i] == 1 {
					dX = k - i
					dY = l - j
					new_dY = dX
					new_dX = dY * (-1)
					new_y = l - new_dY
					new_x = k - new_dX
					copyfield.field[new_y][new_x] = 1
					mfield.field[j][i] = 0
				}
			}
		}

		for i = 0; i < 11; i++ {
			for j = 0; j < 20; j++ {
				if copyfield.field[j][i] != 0 {
					mfield.field[j][i] = copyfield.field[j][i]
				}
			}
		}

	}
}
func rotateR(mfield *mfld) {
	var i, j, k, l int
	var dX, dY int
	var new_dX, new_dY int
	var new_x, new_y int
	var flag int = 0
	for i = 0; i < 20; i++ {
		for j = 0; j < 11; j++ {
			copyfield.field[i][j] = 0
		}
	}
	for i = 0; i < 11; i++ {
		for j = 0; j < 20; j++ {
			if mfield.field[j][i] == 2 {
				k = i
				l = j
			}
		}
	}
	for i = 0; i < 11 && flag == 0; i++ {
		for j = 0; j < 20 && flag == 0; j++ {
			if mfield.field[j][i] == 1 {
				dX = k - i
				dY = l - j
				new_dY = dX * (-1)
				new_dX = dY
				new_y = l - new_dY
				new_x = k - new_dX
				if new_y >= 0 && new_y < 20 {
					if new_x >= 0 && new_x < 11 {
						if mfield.field[new_y][new_x] < 0 {
							flag = 1
						}

					} else {
						flag = 1
					}
				} else {
					flag = 1
				}

			}
		}
	}
	if flag != 1 {
		for i = 0; i < 11; i++ {
			for j = 0; j < 20; j++ {
				if mfield.field[j][i] == 1 {
					dX = k - i
					dY = l - j
					new_dY = dX * (-1)
					new_dX = dY
					new_y = l - new_dY
					new_x = k - new_dX
					copyfield.field[new_y][new_x] = 1
					mfield.field[j][i] = 0
				}
			}
		}

		for i = 0; i < 11; i++ {
			for j = 0; j < 20; j++ {
				if copyfield.field[j][i] != 0 {
					mfield.field[j][i] = copyfield.field[j][i]
				}
			}
		}

	}
}
func manager(mfield *mfld, ch chan int) {
	ex := 0
	var pause bool = false
	var delta int64
	delta = 0
	generate(mfield)
	for ex != 1 {
		//tmpdur = time.Since(myfield.Timestamp)
		if time.Since(myfield.Timestamp).Nanoseconds() > 900000000-delta && !mfield.lock {
			if !pause {
				if fallStep(mfield) {
					generate(mfield)
				}
				delta += 10000
			}
			myfield.Timestamp = time.Now().UTC()
		} else {
			select {
			case c := <-ch:
				switch c {
				case 1:
					rotateL(mfield)
				case 2:
					rotateR(mfield)
				case 3:
					if pause == true {
						pause = false
					} else {
						pause = true
					}
				case 4:
					stepLeft(mfield)
				case 5:
					stepRight(mfield)
				}
			default:
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

/// copyScreen to the render target at a given location.
///
func drawScreen(offsetX int, offsetY int, right int, bottom int) {

	myfield.lock = true

	wx = int((right - offsetX) / countx)
	hy = int((bottom - offsetY) / county)
	//32, 42, 53,
	draw_gopher(80, 294, offsetX, wx, offsetY, hy)
	for i := 0; i < 20; i++ {
		for j := 0; j < 11; j++ {
			rect.X = int32(offsetX + j*wx + 2)
			rect.Y = int32(offsetY + i*hy + 2)
			rect.W = int32(wx - 4)
			rect.H = int32(hy - 4)
			switch myfield.field[i][j] {
			case 0:
				//Renderer.SetDrawColor(32, 42, 53, 255)
				//Renderer.FillRect(&rect)
			case 1:
				Renderer.SetDrawColor(132, 142, 153, 255)
				Renderer.DrawRect(&rect)
				rect.X += 1
				rect.Y += 1
				rect.W -= 2
				rect.H -= 2
				Renderer.SetDrawColor(0, 0, 0, 255)
				Renderer.DrawRect(&rect)
				rect.X += 1
				rect.Y += 1
				rect.W -= 2
				rect.H -= 2
				Renderer.SetDrawColor(90, 90, 190, 255)
				Renderer.FillRect(&rect)
			case 2:
				Renderer.SetDrawColor(132, 142, 153, 255)
				Renderer.DrawRect(&rect)
				rect.X += 1
				rect.Y += 1
				rect.W -= 2
				rect.H -= 2
				Renderer.SetDrawColor(0, 0, 0, 255)
				Renderer.DrawRect(&rect)
				rect.X += 1
				rect.Y += 1
				rect.W -= 2
				rect.H -= 2
				Renderer.SetDrawColor(90, 90, 190, 255)
				Renderer.FillRect(&rect)
			case -1:
				Renderer.SetDrawColor(190, 190, 190, 255)
				Renderer.DrawRect(&rect)
				rect.X += 1
				rect.Y += 1
				rect.W -= 2
				rect.H -= 2
				Renderer.SetDrawColor(0, 0, 0, 255)
				Renderer.DrawRect(&rect)
				rect.X += 1
				rect.Y += 1
				rect.W -= 2
				rect.H -= 2
				Renderer.SetDrawColor(90, 90, 90, 255)
				Renderer.FillRect(&rect)

			}
		}
	}
	drawText("Score:", 40, 500)
	drawText(strconv.Itoa(myfield.counter), 100, 500)
	myfield.lock = false
}

func frame(x, y, w, h int) {
	Renderer.SetDrawColor(0, 0, 0, 255)
	Renderer.DrawLine(x, y, x+w, y)
	Renderer.DrawLine(x, y, x, y+h)

	// highlight
	Renderer.SetDrawColor(95, 112, 120, 255)
	Renderer.DrawLine(x+w, y, x+w, y+h)
	Renderer.DrawLine(x, y+h, x+w, y+h)

}

func main() {
	minit()
	rand.Seed(time.Now().UTC().UnixNano())
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	ch := make(chan int, 3)

	createWindow()
	loadFont()
	loadGopher()
	initField()
	init_generic()

	video := time.NewTicker(time.Second / 100)
	go manager(&myfield, ch)
	for processEvents(ch) {

		select {
		case <-video.C:
			redraw()

		}
		time.Sleep(1 * time.Millisecond)
	}
	sdl.Quit()
}

// processEvents from SDL
///
func processEvents(ch chan int) bool {
	rez := true
	for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
		switch ev := e.(type) {
		case *sdl.QuitEvent:
			return false
		case *sdl.DropEvent:

		case *sdl.KeyDownEvent:
			if _, ok := KeyMap[ev.Keysym.Scancode]; ok {
				//			if key, ok := KeyMap[ev.Keysym.Scancode]; ok {

			} else {
				switch ev.Keysym.Scancode {
				case sdl.SCANCODE_ESCAPE:
					rez = false

				case sdl.SCANCODE_UP:
					ch <- 1
				case sdl.SCANCODE_DOWN:
					ch <- 2
				case sdl.SCANCODE_SPACE:
					ch <- 3
				case sdl.SCANCODE_LEFT:
					ch <- 4
				case sdl.SCANCODE_RIGHT:
					ch <- 5
				}
			}
		case *sdl.KeyUpEvent:

		}
	}

	return rez
}
