// Command generator draws the theme's images: a channel tile and a header strip for every system, the
// backgrounds (the Wii-style menu with its empty slots and the dock, and the game list's panels), the selection
// frame, the list's selector and the "no art" placeholder. Run it from the generator folder:
//
//	go run .
//
// EmulationStation's SVG renderer can't draw text, so everything with words on it is drawn here to PNG.
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const (
	themeDir = "../theme"
	ss       = 3 // images are drawn this many times bigger and scaled down, for smooth edges
)

// The theme's colors. theme.xml uses the same ones.
var (
	paper     = rgb(0xF4F4F4) // the background, between its pinstripes
	pinstripe = rgb(0xE6E6E6)
	slotEdge  = rgb(0xD2D2D2)
	frameEdge = rgb(0xBDBDBD)
	white     = rgb(0xFFFFFF)
	ink       = rgb(0x505050) // text
	inkDim    = rgb(0x8C8C8C)
	blue      = rgb(0x34BEED) // the Wii's highlight blue
	blueDeep  = rgb(0x0D96D6)
)

// The layout at 640×480, in pixels. theme.xml has the same numbers as fractions of the screen.
const (
	screenW, screenH = 640, 480

	// The system view: three rows of channel slots, three to a row. The carousel is the middle row.
	tileW, tileH = 192, 98 // <carousel> logoSize 0.30 0.2042
	rowMid       = 175     // the carousel's center: pos 0 0.2625, size 1 0.2042
	rowStep      = tileH + 12
	colStep      = 640 / 3.0 // ES spaces 3 logos across the carousel's width

	// The dock along the bottom: its top edge is high at the ends, where the round buttons are, and dips in the
	// middle, under the clock.
	dockEnd, dockMid = 353, 386

	// The game list: a header, the games on the left, the art and description on the right.
	headerH          = 60 // the image; the band itself is 52
	listX, listW     = 19, 336
	panelY, panelH   = 72, 258
	detailX, detailW = 374, 247
)

var fonts = map[string]*opentype.Font{}

func main() {
	for _, name := range []string{"Regular", "Medium", "Bold", "ExtraBold"} {
		data, err := os.ReadFile(filepath.Join(themeDir, "fonts", "MPLUSRounded1c-"+name+".ttf"))
		check(err)
		f, err := opentype.Parse(data)
		check(err)
		fonts[name] = f
	}
	check(os.MkdirAll(filepath.Join(themeDir, "art"), 0o755))
	check(os.MkdirAll(filepath.Join(themeDir, "assets"), 0o755))

	for _, s := range systems {
		save(filepath.Join(themeDir, "art", s.Theme+".png"), channel(s))
		save(filepath.Join(themeDir, "art", s.Theme+"-header.png"), header(s))
	}
	save(filepath.Join(themeDir, "assets", "menu.png"), menuBackground())
	save(filepath.Join(themeDir, "assets", "gamelist.png"), gamelistBackground())
	save(filepath.Join(themeDir, "assets", "grid.png"), gridBackground())
	save(filepath.Join(themeDir, "assets", "frame.png"), selectionFrame())
	save(filepath.Join(themeDir, "assets", "selector.png"), selector(634, 58))
	save(filepath.Join(themeDir, "assets", "noart.png"), noArt(300, 300))
	fmt.Printf("drew %d systems\n", len(systems))
}

// channel is a system's tile in the carousel, drawn at twice the size it has on a 640×480 screen: a white bezel
// around a face in the system's color, with a gloss across the top and the system's name on it.
func channel(s System) image.Image {
	const k = 2
	W, H := float64(tileW*k), float64(tileH*k)
	c := newCanvas(int(W), int(H))
	col := rgb(s.Color)

	c.roundRect(0, 0, W, H, 24, frameEdge)
	c.roundRectFn(2, 2, W-4, H-4, 22, vertical(2, H-4, white, rgb(0xECECEC)))
	fx, fy, fw, fh, fr := 9.0, 9.0, W-18, H-18, 16.0
	c.roundRectFn(fx, fy, fw, fh, fr, vertical(fy, fh, shade(col, 0.32), shade(col, -0.08)))
	face := inRoundRect(fx, fy, fw, fh, fr)

	// Two big faint rings on the right, like the patterns on the Wii's channels.
	for _, r := range [][3]float64{{W * 0.86, H * 0.18, H * 0.95}, {W * 0.98, H * 0.95, H * 0.7}} {
		cx, cy, rad := r[0], r[1], r[2]
		c.fill(fx, fy, fw, fh, func(x, y float64) bool {
			d := math.Hypot(x-cx, y-cy)
			return face(x, y) && d < rad && d > rad-14
		}, solid(withAlpha(white, 34)))
	}
	// The gloss: the top half of the face lit, fading down.
	gy, gh := fy, fh*0.48
	c.fill(fx, gy, fw, gh, face, vertical(gy, gh, withAlpha(white, 120), withAlpha(white, 26)))

	c.text(strings.ToUpper(s.Maker), "Bold", 17, fx+16, fy+28, 0.08, withAlpha(white, 215), alignLeft)
	short := 64.0
	if w := textWidth(s.Short, "ExtraBold", short, 0); w > fw*0.86 {
		short *= fw * 0.86 / w
	}
	c.text(s.Short, "ExtraBold", short, W/2+2, fy+fh*0.60+3, 0, withAlpha(shade(col, -0.45), 110), alignCenter)
	c.text(s.Short, "ExtraBold", short, W/2, fy+fh*0.60, 0, white, alignCenter)
	c.fitText(s.Name, "Bold", 21, fw*0.88, W/2, fy+fh*0.88, withAlpha(white, 240), alignCenter)
	return c.done()
}

// header is the band across the top of a system's game list, in the system's color, with a soft shadow under
// it. theme.xml writes the system's full name on the left; its short name is faint on the right.
func header(s System) image.Image {
	c := newCanvas(screenW, headerH)
	col := rgb(s.Color)
	W, band := float64(screenW), 52.0
	c.shadow(-20, -20, W+40, band+20, 0, 7, withAlpha(rgb(0x000000), 70))
	c.rectFn(0, 0, W, band, vertical(0, band, shade(col, 0.28), shade(col, -0.06)))
	c.rectFn(0, 0, W, band*0.5, vertical(0, band*0.5, withAlpha(white, 90), withAlpha(white, 18)))
	c.rect(0, band-2, W, 2, shade(col, -0.3))
	c.text(s.Short, "ExtraBold", band*0.72, W-18, band*0.76, 0, withAlpha(white, 70), alignRight)
	return c.done()
}

// menuBackground is the system view behind the carousel: the pinstriped background, the empty channel slots
// above and below the carousel's row, and the dock with the clock's place and two round buttons.
func menuBackground() image.Image {
	c := newCanvas(screenW, screenH)
	pinstripes(c)
	for _, row := range []float64{rowMid - rowStep, rowMid + rowStep} {
		for i := 0; i < 3; i++ {
			cx := colStep*float64(i) + colStep/2
			emptySlot(c, cx-tileW/2, row-tileH/2)
		}
	}
	dock(c, true)
	return c.done()
}

// gamelistBackground is behind every list but the grid: the white panel for the games, the one for the art and
// description, and the dock.
func gamelistBackground() image.Image {
	c := newCanvas(screenW, screenH)
	pinstripes(c)
	panel(c, listX, panelY, listW, panelH)
	panel(c, detailX, panelY, detailW, panelH)
	c.rect(detailX+14, 236, detailW-28, 1.5, rgb(0xE2E2E2)) // between the art and the description
	dock(c, false)
	return c.done()
}

func gridBackground() image.Image {
	c := newCanvas(screenW, screenH)
	pinstripes(c)
	dock(c, false)
	return c.done()
}

// selectionFrame is the blue outline around the selected channel. It stays in the middle of the carousel's
// row while the channels slide under it. Drawn at twice the size, with room for the glow around the tile.
func selectionFrame() image.Image {
	const k, m = 2, 8 * 2 // m: the margin around the tile, for the glow
	W, H := float64((tileW+16)*k), float64((tileH+16)*k)
	c := newCanvas(int(W), int(H))
	x, y, w, h := float64(m), float64(m), W-2*m, H-2*m
	c.shadow(x-2, y-2, w+4, h+4, 26, 12, withAlpha(blue, 150))
	c.fill(0, 0, W, H, inRoundRect(x+3, y+3, w-6, h-6, 21), solid(color.RGBA{})) // the glow only shows outside the tile
	c.ring(x-4, y-4, w+8, h+8, 28, 7, blue)
	c.ring(x-4, y-4, w+8, h+8, 28, 2, withAlpha(white, 160))
	return c.done()
}

// selector is drawn behind the selected game in the list: a light blue pill with a blue edge, like the Wii's
// highlighted buttons. theme.xml doesn't tint it.
func selector(w, h int) image.Image {
	c := newCanvas(w, h)
	W, H := float64(w), float64(h)
	c.roundRect(0, 0, W, H, H*0.42, blue)
	c.roundRectFn(3, 3, W-6, H-6, H*0.42-3, vertical(3, H-6, rgb(0xF0FBFF), rgb(0xCDEFFC)))
	return c.done()
}

// noArt stands in for the box art of a game that hasn't been scraped: an empty disc slot.
func noArt(w, h int) image.Image {
	c := newCanvas(w, h)
	W, H := float64(w), float64(h)
	c.roundRect(0, 0, W, H, 22, slotEdge)
	c.roundRectFn(3, 3, W-6, H-6, 19, vertical(3, H-6, rgb(0xFAFAFA), rgb(0xE9E9E9)))
	cx, cy := W/2, H*0.44
	c.fill(0, 0, W, H, func(x, y float64) bool {
		d := math.Hypot(x-cx, y-cy)
		return d < H*0.25 && d > H*0.045
	}, solid(rgb(0xDADADA)))
	c.ring(cx-H*0.25, cy-H*0.25, H*0.5, H*0.5, H*0.25, 3, rgb(0xC8C8C8))
	c.text("No Image", "Bold", H*0.085, W/2, H*0.86, 0, rgb(0xA8A8A8), alignCenter)
	return c.done()
}

// ---- the pieces of the backgrounds ----

// pinstripes fills the canvas with the Wii menu's background: light gray, with a darker line every third pixel.
func pinstripes(c *canvas) {
	for y := 0; y < c.h*ss; y++ {
		col := paper
		if (y/ss)%3 == 0 {
			col = pinstripe
		}
		col = shade(col, 0.4*(1-float64(y)/float64(c.h*ss))) // lighter toward the top
		for x := 0; x < c.w*ss; x++ {
			c.img.SetRGBA(x, y, col)
		}
	}
}

// emptySlot is a channel slot with nothing in it.
func emptySlot(c *canvas, x, y float64) {
	c.roundRect(x, y, tileW, tileH, 12, slotEdge)
	c.roundRectFn(x+1.5, y+1.5, tileW-3, tileH-3, 10.5, vertical(y, tileH, rgb(0xFBFBFB), rgb(0xE8E8E8)))
}

// panel is a white rounded box with a gray edge and a soft shadow.
func panel(c *canvas, x, y, w, h float64) {
	c.shadow(x, y+2, w, h, 14, 6, withAlpha(rgb(0x000000), 40))
	c.roundRect(x, y, w, h, 14, rgb(0xCFCFCF))
	c.roundRect(x+1.5, y+1.5, w-3, h-3, 12.5, white)
}

// dockEdge is the y of the dock's top edge at x.
func dockEdge(x float64) float64 {
	W := float64(screenW)
	d := math.Min(x, W-x) / W // 0 at the screen's sides, 0.5 in the middle
	t := smoothstep(0.17, 0.31, d)
	return dockEnd + (dockMid-dockEnd)*t
}

// dock draws the bar along the bottom. withButtons adds the two round buttons at its ends (the system view's;
// the game lists leave the ends empty).
func dock(c *canvas, withButtons bool) {
	W, H := float64(screenW), float64(screenH)
	below := func(off float64) func(x, y float64) bool {
		return func(x, y float64) bool { return y >= dockEdge(x)+off }
	}
	// A soft shadow above the edge.
	for i := 0.0; i < 6; i++ {
		c.fill(0, dockEnd-8, W, H, below(-6+i), solid(withAlpha(rgb(0x000000), 4)))
	}
	c.fill(0, dockEnd-2, W, H, below(0), solid(rgb(0x9FDDF3)))
	c.fill(0, dockEnd, W, H, below(2), vertical(dockEnd, H-dockEnd, rgb(0xFDFDFD), rgb(0xD5D6D8)))
	c.fill(0, dockEnd, W, H, func(x, y float64) bool { return below(2)(x, y) && !below(5)(x, y) }, solid(withAlpha(white, 220)))
	if !withButtons {
		return
	}
	for _, cx := range []float64{44, W - 44} {
		cy, r := 394.0, 29.0
		c.shadow(cx-r, cy-r+2, 2*r, 2*r, r, 5, withAlpha(rgb(0x000000), 50))
		c.roundRect(cx-r, cy-r, 2*r, 2*r, r, rgb(0xB8B8B8))
		c.roundRectFn(cx-r+1.5, cy-r+1.5, 2*r-3, 2*r-3, r-1.5, vertical(cy-r, 2*r, white, rgb(0xE2E2E2)))
		c.ring(cx-r+4, cy-r+4, 2*r-8, 2*r-8, r-4, 2.5, blue)
	}
	// The left button's icon: three lines, a menu. The right one's: a memory card.
	for i := -1.0; i <= 1; i++ {
		c.roundRect(44-10, 394+i*6.5-1.5, 20, 3.5, 1.75, blueDeep)
	}
	cx, cy := W-44, 394.0
	c.fill(cx-9, cy-11, 18, 22, func(x, y float64) bool {
		return inRoundRect(cx-9, cy-11, 18, 22, 2.5)(x, y) && !(x-(cx+9-6)+(cy-11)-y > 0) // a corner clipped off
	}, solid(blueDeep))
	for i := 0.0; i < 3; i++ {
		c.rect(cx-5+i*3.6, cy-8, 2, 5, white)
	}
}

// ---- drawing ----

type canvas struct {
	img  *image.RGBA
	w, h int
}

func newCanvas(w, h int) *canvas {
	return &canvas{img: image.NewRGBA(image.Rect(0, 0, w*ss, h*ss)), w: w, h: h}
}

func (c *canvas) done() image.Image {
	out := image.NewRGBA(image.Rect(0, 0, c.w, c.h))
	xdraw.CatmullRom.Scale(out, out.Bounds(), c.img, c.img.Bounds(), xdraw.Src, nil)
	return out
}

// fill paints the pixels of the box (x, y, w, h) that inside says are in the shape, in the color paint gives
// each one. Coordinates are in canvas units. A fully transparent color cuts a hole instead.
func (c *canvas) fill(x, y, w, h float64, inside func(x, y float64) bool, paint func(x, y float64) color.RGBA) {
	b := c.img.Bounds()
	for py := max(int(y*ss), b.Min.Y); py < min(int(math.Ceil((y+h)*ss)), b.Max.Y); py++ {
		for px := max(int(x*ss), b.Min.X); px < min(int(math.Ceil((x+w)*ss)), b.Max.X); px++ {
			fx, fy := (float64(px)+0.5)/ss, (float64(py)+0.5)/ss
			if !inside(fx, fy) {
				continue
			}
			col := paint(fx, fy)
			if col.A == 0 {
				c.img.SetRGBA(px, py, color.RGBA{})
			} else {
				c.img.SetRGBA(px, py, over(c.img.RGBAAt(px, py), col))
			}
		}
	}
}

func (c *canvas) rect(x, y, w, h float64, col color.RGBA) {
	c.roundRect(x, y, w, h, 0, col)
}

func (c *canvas) rectFn(x, y, w, h float64, paint func(x, y float64) color.RGBA) {
	c.roundRectFn(x, y, w, h, 0, paint)
}

func (c *canvas) roundRect(x, y, w, h, r float64, col color.RGBA) {
	c.roundRectFn(x, y, w, h, r, solid(col))
}

func (c *canvas) roundRectFn(x, y, w, h, r float64, paint func(x, y float64) color.RGBA) {
	c.fill(x, y, w, h, inRoundRect(x, y, w, h, r), paint)
}

// ring draws the edge of a rounded rectangle, t thick, inside its bounds.
func (c *canvas) ring(x, y, w, h, r, t float64, col color.RGBA) {
	outer, inner := inRoundRect(x, y, w, h, r), inRoundRect(x+t, y+t, w-2*t, h-2*t, math.Max(r-t, 0))
	c.fill(x, y, w, h, func(px, py float64) bool { return outer(px, py) && !inner(px, py) }, solid(col))
}

// shadow paints a soft shadow of a rounded rectangle, fading out over blur around its edge.
func (c *canvas) shadow(x, y, w, h, r, blur float64, col color.RGBA) {
	sdf := roundRectDistance(x, y, w, h, r)
	c.fill(x-blur, y-blur, w+2*blur, h+2*blur, func(px, py float64) bool { return sdf(px, py) < blur }, func(px, py float64) color.RGBA {
		a := 1 - smoothstep(-blur, blur, sdf(px, py))
		return withAlpha(col, uint8(math.Max(float64(col.A)*a, 1)))
	})
}

func inRoundRect(x, y, w, h, r float64) func(px, py float64) bool {
	sdf := roundRectDistance(x, y, w, h, r)
	return func(px, py float64) bool { return sdf(px, py) <= 0 }
}

// roundRectDistance gives how far a point is outside a rounded rectangle (negative inside).
func roundRectDistance(x, y, w, h, r float64) func(px, py float64) float64 {
	r = math.Min(r, math.Min(w, h)/2)
	cx, cy, hx, hy := x+w/2, y+h/2, w/2-r, h/2-r
	return func(px, py float64) float64 {
		qx, qy := math.Abs(px-cx)-hx, math.Abs(py-cy)-hy
		return math.Hypot(math.Max(qx, 0), math.Max(qy, 0)) + math.Min(math.Max(qx, qy), 0) - r
	}
}

func solid(col color.RGBA) func(x, y float64) color.RGBA {
	return func(float64, float64) color.RGBA { return col }
}

// vertical is a gradient from top (at y) to bottom (at y+h).
func vertical(y, h float64, top, bottom color.RGBA) func(x, y float64) color.RGBA {
	return func(_, py float64) color.RGBA { return mix(top, bottom, math.Min(math.Max((py-y)/h, 0), 1)) }
}

func smoothstep(a, b, x float64) float64 {
	t := math.Min(math.Max((x-a)/(b-a), 0), 1)
	return t * t * (3 - 2*t)
}

type align int

const (
	alignLeft align = iota
	alignCenter
	alignRight
)

func face(name string, size float64) font.Face {
	f, err := opentype.NewFace(fonts[name], &opentype.FaceOptions{Size: size * ss, DPI: 72, Hinting: font.HintingNone})
	check(err)
	return f
}

// textWidth measures s in canvas units, with tracking (extra space between letters, as a fraction of the size).
func textWidth(s, fontName string, size, tracking float64) float64 {
	f := face(fontName, size)
	defer f.Close()
	w := font.MeasureString(f, s)
	return float64(w)/64/ss + tracking*size*float64(len([]rune(s))-1)
}

// text draws s with its baseline at y.
func (c *canvas) text(s, fontName string, size, x, y, tracking float64, col color.RGBA, a align) {
	w := textWidth(s, fontName, size, tracking)
	switch a {
	case alignCenter:
		x -= w / 2
	case alignRight:
		x -= w
	}
	f := face(fontName, size)
	defer f.Close()
	d := &font.Drawer{Dst: c.img, Src: image.NewUniform(color.NRGBA(col)), Face: f} // col isn't premultiplied
	px := x * ss
	for _, r := range s {
		d.Dot = fixed.Point26_6{X: fixed.Int26_6(px * 64), Y: fixed.Int26_6(y * ss * 64)}
		d.DrawString(string(r))
		adv, _ := f.GlyphAdvance(r)
		px += float64(adv)/64 + tracking*size*ss
	}
}

// fitText draws s at size, or smaller if that's what it takes to fit in maxWidth.
func (c *canvas) fitText(s, fontName string, size, maxWidth, x, y float64, col color.RGBA, a align) {
	if w := textWidth(s, fontName, size, 0); w > maxWidth {
		size *= maxWidth / w
	}
	c.text(s, fontName, size, x, y, 0, col, a)
}

// ---- colors ----

func rgb(v uint32) color.RGBA {
	return color.RGBA{uint8(v >> 16), uint8(v >> 8), uint8(v), 255}
}

func withAlpha(c color.RGBA, a uint8) color.RGBA {
	c.A = a
	return c
}

// shade lightens a color toward white (amount > 0) or darkens it toward black (amount < 0).
func shade(c color.RGBA, amount float64) color.RGBA {
	target := 255.0
	if amount < 0 {
		target, amount = 0, -amount
	}
	f := func(v uint8) uint8 { return uint8(float64(v) + (target-float64(v))*amount) }
	return color.RGBA{f(c.R), f(c.G), f(c.B), c.A}
}

func mix(a, b color.RGBA, t float64) color.RGBA {
	f := func(x, y uint8) uint8 { return uint8(float64(x) + (float64(y)-float64(x))*t) }
	return color.RGBA{f(a.R, b.R), f(a.G, b.G), f(a.B, b.B), f(a.A, b.A)}
}

// over paints src (not premultiplied) over dst (premultiplied, as image.RGBA stores it).
func over(dst, src color.RGBA) color.RGBA {
	a := float64(src.A) / 255
	f := func(d, s uint8) uint8 { return uint8(float64(s)*a + float64(d)*(1-a)) }
	return color.RGBA{f(dst.R, src.R), f(dst.G, src.G), f(dst.B, src.B), uint8(float64(src.A) + float64(dst.A)*(1-a))}
}

func save(path string, img image.Image) {
	f, err := os.Create(path)
	check(err)
	defer f.Close()
	enc := png.Encoder{CompressionLevel: png.BestCompression}
	check(enc.Encode(f, img))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
