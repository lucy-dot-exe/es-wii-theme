package main

import (
	"image/color"
	"math"
)

// The pictures on the channels: each system's controller or console, drawn in white with its details in dark
// gray and the system's color, in a 100×100 box. Drawn from shapes, never the makers' logos.

// pen draws in a 100×100 box placed at (x, y) on the canvas, s canvas units per box unit.
type pen struct {
	c       *canvas
	x, y, s float64
	accent  color.RGBA // the system's color, darkened, for details
}

var (
	body  = rgb(0xFFFFFF)
	dark  = rgb(0x3C3F4A)
	gray  = rgb(0x9A9DA6)
	light = rgb(0xD9DCE2)
	red   = rgb(0xE0453A)
	blueB = rgb(0x3A7BD5)
	green = rgb(0x3FAE5A)
	gold  = rgb(0xF2B01E)
	pink  = rgb(0xE86FA8)
)

func (p pen) X(u float64) float64 { return p.x + u*p.s }
func (p pen) Y(u float64) float64 { return p.y + u*p.s }

func (p pen) rr(x, y, w, h, r float64, col color.RGBA) {
	p.c.roundRect(p.X(x), p.Y(y), w*p.s, h*p.s, r*p.s, col)
}

func (p pen) circle(cx, cy, r float64, col color.RGBA) {
	p.rr(cx-r, cy-r, 2*r, 2*r, r, col)
}

func (p pen) ring(cx, cy, r, t float64, col color.RGBA) {
	p.c.ring(p.X(cx-r), p.Y(cy-r), 2*r*p.s, 2*r*p.s, r*p.s, t*p.s, col)
}

func (p pen) ellipse(cx, cy, rx, ry float64, col color.RGBA) {
	p.shape(cx-rx, cy-ry, 2*rx, 2*ry, func(x, y float64) bool {
		dx, dy := (x-cx)/rx, (y-cy)/ry
		return dx*dx+dy*dy <= 1
	}, col)
}

// shape fills the points of the box (in box units) that inside says are in it.
func (p pen) shape(x, y, w, h float64, inside func(x, y float64) bool, col color.RGBA) {
	p.c.fill(p.X(x), p.Y(y), w*p.s, h*p.s, func(cx, cy float64) bool {
		return inside((cx-p.x)/p.s, (cy-p.y)/p.s)
	}, solid(col))
}

// poly fills a polygon given as x, y pairs.
func (p pen) poly(col color.RGBA, pts ...float64) {
	minX, minY, maxX, maxY := 100.0, 100.0, 0.0, 0.0
	for i := 0; i < len(pts); i += 2 {
		minX, maxX = math.Min(minX, pts[i]), math.Max(maxX, pts[i])
		minY, maxY = math.Min(minY, pts[i+1]), math.Max(maxY, pts[i+1])
	}
	p.shape(minX, minY, maxX-minX, maxY-minY, func(x, y float64) bool {
		in := false
		for i, j := 0, len(pts)-2; i < len(pts); j, i = i, i+2 {
			xi, yi, xj, yj := pts[i], pts[i+1], pts[j], pts[j+1]
			if (yi > y) != (yj > y) && x < (xj-xi)*(y-yi)/(yj-yi)+xi {
				in = !in
			}
		}
		return in
	}, col)
}

// dpad is a cross, size across.
func (p pen) dpad(cx, cy, size float64, col color.RGBA) {
	t := size * 0.34
	p.rr(cx-size/2, cy-t/2, size, t, t*0.2, col)
	p.rr(cx-t/2, cy-size/2, t, size, t*0.2, col)
}

// keys draws rows of keyboard keys in the box.
func (p pen) keys(x, y, w, h float64, cols, rows int, col color.RGBA) {
	kw, kh := w/float64(cols), h/float64(rows)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			p.rr(x+float64(c)*kw+kw*0.12, y+float64(r)*kh+kh*0.14, kw*0.76, kh*0.72, kw*0.15, col)
		}
	}
	p.rr(x+w*0.25, y+h-kh*0.86, w*0.5, kh*0.72, kw*0.15, col) // the space bar, over the last row
}

// ---- controllers ----

func nesPad(p pen) {
	p.rr(4, 30, 92, 40, 4, body)
	p.rr(9, 36, 82, 28, 3, dark)
	p.dpad(23, 50, 20, body)
	p.rr(40, 51, 8, 4, 2, gray)
	p.rr(51, 51, 8, 4, 2, gray)
	p.circle(69, 54, 6, red)
	p.circle(83, 54, 6, red)
}

func famicomPad(p pen) {
	p.rr(4, 32, 92, 36, 6, body)
	p.rr(9, 37, 82, 26, 4, p.accent)
	p.dpad(23, 50, 18, body)
	p.rr(40, 52, 8, 4, 2, light)
	p.rr(51, 52, 8, 4, 2, light)
	p.circle(69, 52, 6, gold)
	p.circle(83, 52, 6, gold)
}

// snesPad is the dog-bone pad. colors are its four face buttons.
func snesPad(p pen, colors [4]color.RGBA) {
	p.circle(27, 50, 21, body)
	p.circle(73, 50, 21, body)
	p.rr(27, 31, 46, 38, 4, body)
	p.dpad(27, 50, 20, dark)
	p.rr(42, 52, 7, 3.5, 1.75, gray)
	p.rr(51, 52, 7, 3.5, 1.75, gray)
	for i, d := range [][2]float64{{0, -9}, {9, 0}, {0, 9}, {-9, 0}} {
		p.circle(73+d[0], 50+d[1], 5, colors[i])
	}
}

func n64Pad(p pen) {
	p.rr(6, 26, 88, 28, 14, body)
	p.rr(10, 40, 20, 42, 10, body)
	p.rr(40, 44, 20, 44, 10, body)
	p.rr(70, 40, 20, 42, 10, body)
	p.dpad(20, 39, 15, dark)
	p.circle(50, 52, 6, gray)
	p.circle(50, 52, 3.5, dark)
	p.circle(71, 42, 4.5, blueB)
	p.circle(66, 34, 4.5, green)
	for _, d := range [][2]float64{{82, 31}, {88, 37}, {82, 43}, {76, 37}} {
		p.circle(d[0], d[1], 2.6, gold)
	}
}

func genesisPad(p pen) {
	p.shape(4, 24, 92, 56, func(x, y float64) bool {
		dx, dy := (x-50)/46, (y-50)/27
		cx, cy := (x-50)/28, (y-90)/20
		return dx*dx+dy*dy <= 1 && cx*cx+cy*cy > 1
	}, body)
	p.dpad(24, 47, 20, dark)
	p.rr(44, 36, 12, 4, 2, gray)
	for i := 0.0; i < 3; i++ {
		p.circle(60+i*12, 54-i*4, 5.5, dark)
	}
}

func saturnPad(p pen) {
	p.circle(27, 50, 21, body)
	p.circle(73, 50, 21, body)
	p.rr(27, 31, 46, 38, 4, body)
	p.dpad(27, 50, 20, dark)
	p.rr(46, 44, 8, 3.5, 1.75, gray)
	for r := 0.0; r < 2; r++ {
		for c := 0.0; c < 3; c++ {
			p.circle(62+c*10, 44+r*11-c*2, 4, []color.RGBA{dark, p.accent}[int(r)])
		}
	}
}

func dreamcastPad(p pen) {
	p.rr(10, 22, 80, 46, 18, body)
	p.rr(12, 44, 24, 42, 12, body)
	p.rr(64, 44, 24, 42, 12, body)
	p.rr(38, 26, 24, 24, 4, dark)
	p.rr(41, 29, 18, 10, 2, rgb(0x8FB7A4))
	p.circle(24, 38, 7, light)
	p.circle(24, 38, 4, gray)
	p.dpad(26, 58, 14, dark)
	for i, d := range [][2]float64{{76, 33}, {84, 41}, {76, 49}, {68, 41}} {
		p.circle(d[0], d[1], 4, []color.RGBA{gold, red, blueB, green}[i])
	}
}

func psPad(p pen) {
	p.rr(14, 28, 72, 30, 14, body)
	p.rr(10, 36, 26, 50, 13, body)
	p.rr(64, 36, 26, 50, 13, body)
	p.dpad(25, 45, 17, dark)
	p.rr(43, 44, 6, 3, 1.5, gray)
	p.rr(51, 44, 6, 3, 1.5, gray)
	p.ring(75, 37, 3.6, 1.4, green)
	p.ring(83, 45, 3.6, 1.4, red)
	p.ring(75, 53, 3.6, 1.4, blueB)
	p.ring(67, 45, 3.6, 1.4, pink)
}

func pcePad(p pen, turbo bool) {
	p.rr(5, 32, 90, 36, 6, body)
	p.dpad(22, 50, 20, dark)
	if turbo {
		p.rr(41, 37, 5, 8, 1.5, gray)
		p.rr(50, 37, 5, 8, 1.5, gray)
	}
	p.rr(40, 54, 7, 3.5, 1.75, gray)
	p.rr(50, 54, 7, 3.5, 1.75, gray)
	p.circle(68, 52, 6.5, p.accent)
	p.circle(84, 52, 6.5, p.accent)
}

func modernPad(p pen) {
	p.rr(14, 28, 72, 32, 15, body)
	p.rr(8, 38, 28, 46, 14, body)
	p.rr(64, 38, 28, 46, 14, body)
	p.circle(26, 44, 6.5, dark)
	p.dpad(38, 60, 12, dark)
	p.circle(60, 60, 5.5, dark)
	for i, d := range [][2]float64{{75, 36}, {83, 44}, {75, 52}, {67, 44}} {
		p.circle(d[0], d[1], 3.6, []color.RGBA{gold, red, green, blueB}[i])
	}
}

// keypadPad is a controller with a number pad (ColecoVision, Intellivision, Jaguar).
func keypadPad(p pen, wide bool) {
	if wide {
		p.rr(6, 22, 88, 40, 18, body)
		p.dpad(22, 40, 18, dark)
		for i := 0.0; i < 3; i++ {
			p.circle(64+i*10, 46-i*5, 4.5, p.accent)
		}
		p.rr(32, 54, 36, 36, 6, body)
		p.keys(36, 58, 28, 28, 3, 4, gray)
		return
	}
	p.rr(28, 6, 44, 88, 8, body)
	p.circle(50, 22, 10, dark)
	p.circle(50, 22, 4, gray)
	p.keys(33, 40, 34, 46, 3, 4, light)
	p.rr(28, 40, 4, 20, 2, p.accent)
	p.rr(68, 40, 4, 20, 2, p.accent)
}

// joystick is the one-button stick of the Atari 2600 and its kin.
func joystick(p pen) {
	p.rr(18, 50, 64, 40, 8, body)
	p.rr(24, 56, 52, 28, 5, dark)
	p.rr(45, 14, 10, 50, 5, dark)
	p.circle(50, 16, 7, p.accent)
	p.circle(30, 58, 7, red)
}

func arcadeStick(p pen) {
	p.rr(4, 46, 92, 32, 8, body)
	p.rr(22, 26, 5, 26, 2.5, gray)
	p.circle(24.5, 26, 9, red)
	for i, c := range []color.RGBA{red, gold, green, blueB} {
		p.circle(48+float64(i)*13, 60-float64(i%2)*4, 5.5, c)
	}
}

func cabinet(p pen) {
	p.poly(body, 30, 6, 70, 6, 70, 40, 78, 54, 78, 94, 22, 94, 22, 54, 30, 40)
	p.rr(32, 10, 36, 10, 2, p.accent)
	p.rr(33, 24, 34, 24, 2, dark)
	p.rr(37, 28, 12, 3, 1, gray)
	p.poly(light, 26, 54, 74, 54, 76, 62, 24, 62)
	p.circle(38, 57, 2.6, red)
	p.circle(54, 58, 2, gold)
	p.circle(62, 58, 2, blueB)
	p.rr(40, 72, 20, 14, 2, light)
}

// ---- handhelds ----

// tallHandheld is a Game Boy shaped handheld. screen is the color of its picture.
func tallHandheld(p pen, screen color.RGBA, buttons color.RGBA) {
	p.rr(24, 4, 52, 92, 6, body)
	p.rr(29, 10, 42, 36, 3, dark)
	p.rr(35, 15, 30, 26, 1, screen)
	p.dpad(37, 62, 15, dark)
	p.circle(64, 60, 5, buttons)
	p.circle(55, 66, 5, buttons)
	p.rr(41, 76, 7, 3, 1.5, gray)
	p.rr(51, 76, 7, 3, 1.5, gray)
	for i := 0.0; i < 4; i++ {
		p.rr(58+i*3.5, 82, 1.6, 9, 0.8, light)
	}
}

// wideHandheld is a handheld held sideways (Game Boy Advance, Game Gear, Lynx, ...). r is how round its ends are.
func wideHandheld(p pen, r float64, screen color.RGBA, buttons int) {
	p.rr(3, 28, 94, 44, r, body)
	p.rr(29, 31, 42, 38, 4, dark)
	p.rr(33, 35, 34, 30, 1, screen)
	p.dpad(15, 50, 15, dark)
	for i := 0; i < buttons; i++ {
		p.circle(85-float64(i)*8, 46+float64(i)*7, 4.5, p.accent)
	}
}

func clamshell(p pen) {
	p.rr(18, 4, 64, 42, 6, body)
	p.rr(26, 9, 48, 32, 2, dark)
	p.rr(29, 12, 42, 26, 1, rgb(0x9CC9E8))
	p.rr(22, 46, 56, 4, 2, light)
	p.rr(18, 50, 64, 44, 6, body)
	p.rr(33, 55, 34, 28, 2, dark)
	p.rr(35, 57, 30, 24, 1, rgb(0xC9E4F2))
	p.dpad(25, 64, 10, dark)
	for i, d := range [][2]float64{{75, 60}, {79, 64}, {75, 68}, {71, 64}} {
		p.circle(d[0], d[1], 2, []color.RGBA{dark, dark, dark, dark}[i])
	}
}

func pspShape(p pen) {
	p.rr(2, 32, 96, 36, 18, body)
	p.rr(25, 34, 50, 32, 2, dark)
	p.rr(27, 36, 46, 28, 1, rgb(0x6C8EBF))
	p.dpad(13, 46, 12, dark)
	p.circle(13, 60, 3.5, gray)
	p.ring(87, 42, 2.6, 1, green)
	p.ring(92, 48, 2.6, 1, red)
	p.ring(87, 54, 2.6, 1, blueB)
	p.ring(82, 48, 2.6, 1, pink)
}

func ngpShape(p pen) {
	p.rr(6, 26, 88, 48, 12, body)
	p.circle(20, 52, 9, dark)
	p.circle(20, 52, 5, gray)
	p.rr(34, 30, 40, 38, 4, dark)
	p.rr(38, 34, 32, 30, 1, rgb(0xA9D8C8))
	p.circle(84, 46, 4.5, p.accent)
	p.circle(84, 58, 4.5, p.accent)
}

func wonderswan(p pen) {
	p.rr(4, 30, 92, 40, 8, body)
	p.rr(22, 33, 46, 34, 3, dark)
	p.rr(25, 36, 40, 28, 1, rgb(0xB8C7B0))
	for i := 0.0; i < 4; i++ {
		p.circle(12, 38+i*8, 2.6, gray)
	}
	p.circle(80, 44, 5, p.accent)
	p.circle(88, 54, 5, p.accent)
}

func virtualBoy(p pen) {
	p.rr(8, 18, 84, 34, 14, body)
	p.rr(14, 24, 72, 18, 9, rgb(0x2A2A30))
	p.circle(36, 33, 6, red)
	p.circle(64, 33, 6, red)
	p.poly(body, 46, 52, 54, 52, 56, 70, 44, 70)
	p.poly(body, 44, 68, 56, 68, 74, 92, 26, 92)
	p.rr(30, 88, 40, 6, 3, light)
}

func vmuShape(p pen) {
	p.rr(28, 10, 44, 80, 12, body)
	p.rr(33, 18, 34, 28, 3, dark)
	p.rr(36, 21, 28, 22, 1, rgb(0x8FB7A4))
	p.dpad(40, 62, 14, dark)
	p.circle(61, 58, 4.5, p.accent)
	p.circle(61, 70, 4.5, p.accent)
}

func gameAndWatch(p pen) {
	p.rr(4, 28, 92, 44, 5, rgb(0xE8D9B0))
	p.rr(28, 32, 44, 36, 2, dark)
	p.rr(31, 35, 38, 30, 1, rgb(0xC9CDB8))
	p.dpad(15, 50, 15, dark)
	p.circle(85, 54, 6, red)
	p.rr(78, 34, 14, 5, 2.5, gold)
}

func smallHandheld(p pen, screen color.RGBA) {
	p.rr(26, 18, 48, 64, 14, body)
	p.rr(32, 24, 36, 28, 3, dark)
	p.rr(35, 27, 30, 22, 1, screen)
	p.circle(42, 66, 5, p.accent)
	p.circle(58, 66, 5, p.accent)
	p.circle(50, 58, 3, gray)
}

func arduboy(p pen) {
	p.rr(24, 6, 52, 88, 5, body)
	p.rr(30, 12, 40, 30, 2, dark)
	p.rr(33, 15, 34, 24, 1, rgb(0x1E2126))
	p.rr(37, 22, 10, 3, 1, body)
	p.rr(51, 30, 12, 3, 1, body)
	p.dpad(38, 66, 16, dark)
	p.circle(64, 62, 5, p.accent)
	p.circle(56, 70, 5, p.accent)
}

// ---- consoles, computers and media ----

func consoleBox(p pen) {
	p.rr(6, 40, 88, 38, 7, body)
	p.rr(22, 44, 56, 7, 3, dark)
	p.rr(28, 24, 44, 22, 3, p.accent) // a game in the slot
	p.rr(34, 28, 32, 10, 2, light)
	p.circle(16, 66, 4, red)
	p.rr(62, 63, 22, 6, 3, light)
}

func vectrex(p pen) {
	p.rr(24, 4, 52, 74, 6, body)
	p.rr(30, 10, 40, 60, 3, dark)
	p.poly(rgb(0x7FD4F0), 40, 50, 50, 22, 60, 50, 57, 50, 50, 31, 43, 50)
	p.rr(18, 80, 64, 14, 5, body)
	p.circle(30, 87, 3.5, dark)
	for i := 0.0; i < 4; i++ {
		p.circle(52+i*7, 87, 2.4, red)
	}
}

func computer(p pen, screen color.RGBA) {
	p.rr(18, 4, 64, 48, 5, body)
	p.rr(24, 10, 52, 36, 2, dark)
	p.rr(27, 13, 46, 30, 1, screen)
	p.rr(31, 18, 20, 3, 1, body)
	p.rr(31, 25, 30, 3, 1, light)
	p.rr(42, 52, 16, 6, 1, light)
	p.rr(6, 62, 88, 28, 5, body)
	p.keys(10, 66, 80, 20, 12, 3, light)
}

// keyboardComputer is a computer in its keyboard (C64, Amiga 500, ZX Spectrum, MSX). keyColor is its keys'.
func keyboardComputer(p pen, keyColor color.RGBA) {
	p.poly(body, 10, 34, 90, 34, 97, 74, 3, 74)
	p.rr(12, 38, 76, 6, 2, p.accent)
	p.keys(10, 47, 80, 24, 13, 3, keyColor)
}

func disc(p pen) {
	p.circle(50, 50, 40, body)
	p.ring(50, 50, 30, 1.2, light)
	p.ring(50, 50, 22, 1.2, light)
	p.circle(50, 50, 13, light)
	p.circle(50, 50, 6, p.accent)
	p.shape(10, 10, 80, 80, func(x, y float64) bool {
		a := math.Atan2(y-50, x-50)
		d := math.Hypot(x-50, y-50)
		return d < 39 && d > 14 && a > -2.6 && a < -1.9
	}, rgb(0xC7E3F5))
}

func floppy(p pen) {
	p.poly(body, 14, 10, 78, 10, 90, 22, 90, 90, 14, 90)
	p.rr(30, 10, 38, 26, 2, light)
	p.rr(54, 14, 9, 18, 1, dark)
	p.rr(24, 50, 56, 34, 3, p.accent)
	p.rr(28, 56, 48, 4, 2, light)
	p.rr(28, 64, 36, 4, 2, light)
}

func cartridge(p pen) {
	p.poly(body, 16, 8, 84, 8, 84, 92, 70, 92, 70, 86, 30, 86, 30, 92, 16, 92)
	for i := 0.0; i < 5; i++ {
		p.rr(30+i*9, 12, 4, 8, 2, light)
	}
	p.rr(24, 26, 52, 48, 4, p.accent)
	p.rr(30, 34, 40, 6, 3, body)
	p.rr(30, 46, 28, 4, 2, light)
}

func dish(p pen) {
	p.shape(10, 14, 70, 70, func(x, y float64) bool {
		dx, dy := (x-44)/32, (y-44)/26
		return dx*dx+dy*dy <= 1 && x+y < 104
	}, body)
	p.rr(56, 22, 5, 30, 2.5, gray)
	p.circle(58, 20, 5, p.accent)
	p.poly(body, 40, 66, 52, 66, 62, 92, 30, 92)
}

// pixelScreen is a small monitor with a pixel sprite on it, for the fantasy consoles.
func pixelScreen(p pen) {
	p.rr(14, 8, 72, 62, 7, body)
	p.rr(20, 14, 60, 50, 3, dark)
	sprite := []string{
		"..X...X..",
		"...X.X...",
		"..XXXXX..",
		".XX.X.XX.",
		"XXXXXXXXX",
		"X.XXXXX.X",
		"X.X...X.X",
		"...X.X...",
	}
	const px = 4.4
	for r, row := range sprite {
		for c, ch := range row {
			if ch == 'X' {
				p.rr(30.2+float64(c)*px, 21+float64(r)*px, px, px, 0, p.accent)
			}
		}
	}
	p.rr(40, 70, 20, 8, 1, light)
	p.rr(26, 78, 48, 12, 5, body)
}

func pointer(p pen) {
	p.rr(40, 10, 54, 34, 10, body)
	p.poly(body, 52, 40, 64, 40, 50, 54)
	for i := 0.0; i < 3; i++ {
		p.circle(56+i*12, 27, 4, p.accent)
	}
	p.poly(body, 10, 40, 10, 92, 22, 80, 31, 98, 39, 94, 30, 76, 46, 76)
	p.poly(dark, 15, 52, 15, 80, 22, 73, 29, 88, 32, 86, 25, 71, 34, 71)
}

func chest(p pen) {
	p.rr(10, 22, 80, 30, 14, body)
	p.rr(10, 44, 80, 44, 4, body)
	p.rr(10, 44, 80, 6, 0, p.accent)
	p.rr(43, 40, 14, 18, 3, gold)
	p.circle(50, 49, 2.6, dark)
	p.rr(22, 22, 6, 66, 0, light)
	p.rr(72, 22, 6, 66, 0, light)
}

func sword(p pen) {
	p.poly(body, 46, 6, 54, 6, 56, 62, 50, 70, 44, 62)
	p.rr(48.5, 12, 3, 46, 1.5, light)
	p.rr(28, 62, 44, 8, 4, gold)
	p.rr(46, 70, 8, 18, 3, p.accent)
	p.circle(50, 91, 5, gold)
}

func heart(p pen) {
	p.circle(33, 38, 20, body)
	p.circle(67, 38, 20, body)
	p.poly(body, 14, 44, 86, 44, 50, 88)
	p.circle(30, 34, 6, withAlpha(p.accent, 90))
}

func gear(p pen) {
	for i := 0; i < 8; i++ {
		a := float64(i) * math.Pi / 4
		pts := []float64{}
		for _, c := range [][2]float64{{-7, 26}, {7, 26}, {6, 42}, {-6, 42}} {
			x, y := c[0], c[1]
			pts = append(pts, 50+x*math.Cos(a)-y*math.Sin(a), 50+x*math.Sin(a)+y*math.Cos(a))
		}
		p.poly(body, pts...)
	}
	p.circle(50, 50, 30, body)
	p.circle(50, 50, 12, p.accent)
}

// ---- the collections ----

func star(p pen, cx, cy, r float64, col color.RGBA) {
	pts := []float64{}
	for i := 0; i < 10; i++ {
		rad := r
		if i%2 == 1 {
			rad = r * 0.45
		}
		a := -math.Pi/2 + float64(i)*math.Pi/5
		pts = append(pts, cx+rad*math.Cos(a), cy+rad*math.Sin(a))
	}
	p.poly(col, pts...)
}

func allGames(p pen) {
	for i := 0; i < 4; i++ {
		x, y := 12+float64(i%2)*40, 12+float64(i/2)*40
		p.rr(x, y, 36, 36, 8, body)
		p.rr(x+6, y+6, 24, 24, 4, []color.RGBA{p.accent, light, light, p.accent}[i])
	}
}

func clock(p pen) {
	p.circle(50, 50, 40, body)
	p.ring(50, 50, 40, 4, light)
	for i := 0; i < 12; i++ {
		a := float64(i) * math.Pi / 6
		p.circle(50+31*math.Cos(a), 50+31*math.Sin(a), 2, gray)
	}
	p.rr(47, 24, 6, 30, 3, dark)
	p.poly(dark, 49, 47, 53, 51, 72, 66, 68, 70)
	p.circle(50, 50, 5, p.accent)
}

func sparkle(p pen) {
	four := func(cx, cy, r float64, col color.RGBA) {
		w := r * 0.28
		p.poly(col, cx, cy-r, cx+w, cy-w, cx+r, cy, cx+w, cy+w, cx, cy+r, cx-w, cy+w, cx-r, cy, cx-w, cy-w)
	}
	four(42, 54, 36, body)
	four(78, 22, 16, body)
	four(80, 80, 10, body)
}

func people(p pen, n int) {
	place := [][3]float64{{32, 30, 1}, {68, 30, 1}}
	if n == 4 {
		place = [][3]float64{{14, 28, 0.6}, {38, 24, 0.6}, {62, 24, 0.6}, {86, 28, 0.6}}
	}
	for _, q := range place {
		cx, top, k := q[0], q[1], q[2]
		p.circle(cx, top+12*k, 11*k, body)
		p.rr(cx-17*k, top+26*k, 34*k, 36*k, 14*k, body)
	}
}

func folder(p pen) {
	p.rr(8, 18, 36, 16, 5, body)
	p.rr(8, 26, 84, 60, 7, body)
	p.rr(8, 36, 84, 50, 7, light)
	star(p, 50, 61, 15, p.accent)
}

// icons says how to draw each system. Systems not in it get a cartridge.
var icons = map[string]func(p pen){
	"nes":          nesPad,
	"famicom":      famicomPad,
	"fds":          floppy,
	"snes":         func(p pen) { snesPad(p, [4]color.RGBA{rgb(0x8F84C4), rgb(0x8F84C4), rgb(0x5B4C9E), rgb(0x5B4C9E)}) },
	"sfc":          func(p pen) { snesPad(p, [4]color.RGBA{green, red, gold, blueB}) },
	"snes-hacks":   func(p pen) { snesPad(p, [4]color.RGBA{green, red, gold, blueB}) },
	"snesmsu1":     disc,
	"sufami":       cartridge,
	"satellaview":  dish,
	"sgb":          func(p pen) { tallHandheld(p, rgb(0xB9A6E0), rgb(0x8A5BB0)) },
	"n64":          n64Pad,
	"n64dd":        floppy,
	"gb":           func(p pen) { tallHandheld(p, rgb(0x9BBC0F), rgb(0xA8285E)) },
	"gbc":          func(p pen) { tallHandheld(p, rgb(0x8FD0E8), p.accent) },
	"gba":          func(p pen) { wideHandheld(p, 22, rgb(0x9FB8E8), 2) },
	"nds":          clamshell,
	"virtualboy":   virtualBoy,
	"gameandwatch": gameAndWatch,
	"pokemonmini":  func(p pen) { smallHandheld(p, rgb(0xB8C7B0)) },

	"sg-1000":      joystick,
	"mastersystem": func(p pen) { pcePad(p, false) },
	"megadrive":    genesisPad,
	"genesis":      genesisPad,
	"segacd":       disc,
	"sega32x":      consoleBox,
	"saturn":       saturnPad,
	"dreamcast":    dreamcastPad,
	"vmu":          vmuShape,
	"gamegear":     func(p pen) { wideHandheld(p, 8, rgb(0x9FB8E8), 2) },
	"naomi":        cabinet,
	"atomiswave":   cabinet,

	"psx":      psPad,
	"psp":      pspShape,
	"pspminis": pspShape,

	"pcengine":     func(p pen) { pcePad(p, true) },
	"pcenginecd":   disc,
	"turbografx":   func(p pen) { pcePad(p, true) },
	"turbografxcd": disc,
	"supergrafx":   func(p pen) { pcePad(p, true) },
	"pcfx":         psPad,
	"pc98":         func(p pen) { computer(p, rgb(0x3A5A9C)) },

	"neogeo":   arcadeStick,
	"neogeocd": disc,
	"ngp":      ngpShape,
	"ngpc":     ngpShape,

	"arcade":   cabinet,
	"mame":     cabinet,
	"mame2003": cabinet,
	"cps1":     arcadeStick,
	"cps2":     arcadeStick,
	"cps3":     arcadeStick,

	"atari2600":   joystick,
	"atari5200":   func(p pen) { keypadPad(p, false) },
	"atari7800":   joystick,
	"atari800":    func(p pen) { keyboardComputer(p, light) },
	"atarixegs":   func(p pen) { keyboardComputer(p, light) },
	"atarist":     func(p pen) { keyboardComputer(p, light) },
	"atarilynx":   func(p pen) { wideHandheld(p, 10, rgb(0x9FB8E8), 2) },
	"atarijaguar": func(p pen) { keypadPad(p, true) },

	"wonderswan":      wonderswan,
	"wonderswancolor": wonderswan,
	"supervision":     func(p pen) { wideHandheld(p, 14, rgb(0xB8C7B0), 2) },
	"megaduck":        func(p pen) { tallHandheld(p, rgb(0xB8C7B0), p.accent) },
	"arduboy":         arduboy,

	"3do":           consoleBox,
	"coleco":        func(p pen) { keypadPad(p, false) },
	"intellivision": func(p pen) { keypadPad(p, false) },
	"odyssey2":      func(p pen) { keyboardComputer(p, light) },
	"videopac":      func(p pen) { keyboardComputer(p, light) },
	"vectrex":       vectrex,
	"channelf":      joystick,
	"gx4000":        consoleBox,

	"amiga":      func(p pen) { keyboardComputer(p, light) },
	"amigacd32":  disc,
	"c64":        func(p pen) { keyboardComputer(p, rgb(0x6E6458)) },
	"c128":       func(p pen) { keyboardComputer(p, rgb(0x6E6458)) },
	"c16":        func(p pen) { keyboardComputer(p, dark) },
	"vic20":      func(p pen) { keyboardComputer(p, rgb(0x6E6458)) },
	"amstradcpc": func(p pen) { computer(p, rgb(0x2B3C8F)) },
	"msx":        func(p pen) { keyboardComputer(p, gray) },
	"msx2":       func(p pen) { keyboardComputer(p, gray) },
	"dos":        func(p pen) { computer(p, rgb(0x1E2A5A)) },
	"zxspectrum": func(p pen) { keyboardComputer(p, dark) },
	"zx81":       func(p pen) { keyboardComputer(p, dark) },
	"x1":         func(p pen) { computer(p, rgb(0x2B2B33)) },
	"x68000":     func(p pen) { computer(p, rgb(0x2B2B33)) },

	"ports":    modernPad,
	"doom":     modernPad,
	"scummvm":  pointer,
	"easyrpg":  chest,
	"openbor":  modernPad,
	"solarus":  sword,
	"love2d":   heart,
	"pico-8":   pixelScreen,
	"tic80":    pixelScreen,
	"wasm4":    pixelScreen,
	"lowresnx": pixelScreen,
	"uzebox":   pixelScreen,
	"retropie": gear,

	"auto-allgames":      allGames,
	"auto-favorites":     func(p pen) { star(p, 50, 52, 44, body) },
	"auto-lastplayed":    clock,
	"auto-neverplayed":   sparkle,
	"auto-at2players":    func(p pen) { people(p, 2) },
	"auto-at4players":    func(p pen) { people(p, 4) },
	"custom-collections": folder,
}

// drawFaint draws s's picture as a pale shape in tints of its color, at alpha, for the channel's background.
func drawFaint(c *canvas, s System, x, y, size, alpha float64) {
	col := rgb(s.Color)
	ob, ol := body, light
	body, light = shade(col, 0.55), shade(col, 0.75)
	defer func() { body, light = ob, ol }()
	layer := newCanvas(c.w, c.h)
	draw := icons[s.Theme]
	if draw == nil {
		draw = cartridge
	}
	draw(pen{c: layer, x: x, y: y, s: size / 100, accent: shade(col, -0.1)})
	b := c.img.Bounds()
	for py := b.Min.Y; py < b.Max.Y; py++ {
		for px := b.Min.X; px < b.Max.X; px++ {
			l := layer.img.RGBAAt(px, py)
			if l.A == 0 {
				continue
			}
			a := float64(l.A) / 255
			st := color.RGBA{uint8(float64(l.R) / a), uint8(float64(l.G) / a), uint8(float64(l.B) / a), uint8(float64(l.A) * alpha)}
			c.img.SetRGBA(px, py, over(c.img.RGBAAt(px, py), st))
		}
	}
}

// drawIcon draws s's picture in the box (x, y, size×size) on c, with a soft shadow under it.
func drawIcon(c *canvas, s System, x, y, size float64) {
	layer := newCanvas(c.w, c.h)
	draw := icons[s.Theme]
	if draw == nil {
		draw = cartridge
	}
	draw(pen{c: layer, x: x, y: y, s: size / 100, accent: shade(rgb(s.Color), -0.15)})
	c.composite(layer, size*0.025, size*0.035, withAlpha(shade(rgb(s.Color), -0.6), 90))
}

// composite paints layer over c, first its shadow: its shape in col, offset by (dx, dy).
func (c *canvas) composite(layer *canvas, dx, dy float64, col color.RGBA) {
	ox, oy := int(dx*ss), int(dy*ss)
	b := c.img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if sx, sy := x-ox, y-oy; sx >= 0 && sy >= 0 {
				if a := layer.img.RGBAAt(sx, sy).A; a > 0 {
					c.img.SetRGBA(x, y, over(c.img.RGBAAt(x, y), withAlpha(col, uint8(int(col.A)*int(a)/255))))
				}
			}
		}
	}
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			s := layer.img.RGBAAt(x, y)
			if s.A == 0 {
				continue
			}
			// layer is premultiplied; over() takes a straight color.
			a := float64(s.A) / 255
			st := color.RGBA{uint8(float64(s.R) / a), uint8(float64(s.G) / a), uint8(float64(s.B) / a), s.A}
			c.img.SetRGBA(x, y, over(c.img.RGBAAt(x, y), st))
		}
	}
}
