// Command generator draws the theme's images: a cartridge card and a header strip for every system, the
// background, the "no art" placeholder, and the previews in previews/. Run it from the generator folder:
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
	themeDir   = "../theme"
	previewDir = "../previews"
	ss         = 3 // images are drawn this many times bigger and scaled down, for smooth edges
)

// The theme's colors. theme.xml uses the same ones.
var (
	bgTop    = rgb(0x1E2030)
	bgBottom = rgb(0x12131C)
	ink      = rgb(0xF4F1EA)
	inkDim   = rgb(0x9A9DB0)
	accent   = rgb(0xF5B841)
)

var fonts = map[string]*opentype.Font{}

func main() {
	for _, name := range []string{"Lato-Regular", "Lato-Bold", "Lato-Black"} {
		data, err := os.ReadFile(filepath.Join(themeDir, "fonts", name+".ttf"))
		check(err)
		f, err := opentype.Parse(data)
		check(err)
		fonts[name] = f
	}
	check(os.MkdirAll(filepath.Join(themeDir, "art"), 0o755))
	check(os.MkdirAll(filepath.Join(themeDir, "assets"), 0o755))
	check(os.MkdirAll(previewDir, 0o755))

	for _, s := range systems {
		save(filepath.Join(themeDir, "art", s.Theme+".png"), card(s, 480, 300))
		save(filepath.Join(themeDir, "art", s.Theme+"-header.png"), header(s, 640, 56))
	}
	save(filepath.Join(themeDir, "assets", "background.png"), background(640, 480))
	save(filepath.Join(themeDir, "assets", "noart.png"), noArt(300, 300))
	save(filepath.Join(themeDir, "assets", "selector.png"), selector(600, 60))

	byTheme := map[string]System{}
	for _, s := range systems {
		byTheme[s.Theme] = s
	}
	save(filepath.Join(previewDir, "system.png"), previewSystem(byTheme["snes"], byTheme["gba"], byTheme["megadrive"]))
	save(filepath.Join(previewDir, "gamelist.png"), previewGamelist(byTheme["gba"]))
	fmt.Printf("drew %d systems\n", len(systems))
}

// card is a system's carousel image: a cartridge in the system's color with a label on it.
func card(s System, w, h int) image.Image {
	c := newCanvas(w, h)
	body := rgb(s.Color)
	W, H := float64(w), float64(h)

	// The cartridge's shell: a rounded body with a notch cut out of the bottom, where the pins would be.
	c.roundRect(0, 0, W, H, 26, body)
	c.roundRect(W*0.30, H*0.93, W*0.40, H*0.2, 10, color.RGBA{})
	// Grip ridges on the top edge.
	for i := 0; i < 6; i++ {
		x := W*0.5 - 6*13/2.0 + float64(i)*13
		c.roundRect(x, H*0.035, 6, H*0.07, 3, shade(body, -0.22))
	}
	// The label, a lighter inset with the system on it.
	lx, ly, lw, lh := W*0.07, H*0.16, W*0.86, H*0.68
	c.roundRect(lx, ly, lw, lh, 16, shade(body, 0.12))
	c.roundRect(lx, ly+lh*0.80, lw, lh*0.20, 16, shade(body, -0.18))
	c.rect(lx, ly+lh*0.80, lw, lh*0.10, shade(body, -0.18))

	c.text(strings.ToUpper(s.Maker), "Lato-Bold", H*0.075, lx+W*0.04, ly+lh*0.17, 0.12, withAlpha(ink, 200), alignLeft)
	c.fitText(s.Short, "Lato-Black", H*0.30, lw*0.88, W*0.5, ly+lh*0.62, ink, alignCenter)
	c.fitText(s.Name, "Lato-Bold", H*0.07, lw*0.9, W*0.5, ly+lh*0.935, ink, alignCenter)
	return c.done()
}

// header is the strip across the top of a system's game list, in the system's color.
func header(s System, w, h int) image.Image {
	c := newCanvas(w, h)
	body := rgb(s.Color)
	c.rect(0, 0, float64(w), float64(h), body)
	c.rect(0, float64(h)-4, float64(w), 4, shade(body, -0.3))
	// The system's short name, faint, on the right; theme.xml writes its full name on the left.
	c.text(s.Short, "Lato-Black", float64(h)*0.8, float64(w)-18, float64(h)*0.78, 0, withAlpha(ink, 60), alignRight)
	return c.done()
}

func background(w, h int) image.Image {
	c := newCanvas(w, h)
	for y := 0; y < h*ss; y++ {
		t := float64(y) / float64(h*ss-1)
		col := mix(bgTop, bgBottom, t)
		for x := 0; x < w*ss; x++ {
			// Faint diagonal stripes.
			if (x+y)/(18*ss)%2 == 0 {
				c.img.Set(x, y, shade(col, 0.025))
			} else {
				c.img.Set(x, y, col)
			}
		}
	}
	return c.done()
}

// noArt stands in for the box art of a game that hasn't been scraped.
func noArt(w, h int) image.Image {
	c := newCanvas(w, h)
	W, H := float64(w), float64(h)
	c.roundRect(0, 0, W, H, 18, withAlpha(ink, 22))
	c.roundRect(6, 6, W-12, H-12, 14, withAlpha(bgBottom, 200))
	c.text("NO ART", "Lato-Black", H*0.11, W/2, H*0.54, 0.15, withAlpha(ink, 90), alignCenter)
	return c.done()
}

// selector is drawn behind the selected game in the list. theme.xml tints it.
func selector(w, h int) image.Image {
	c := newCanvas(w, h)
	c.roundRect(0, 0, float64(w), float64(h), float64(h)*0.3, color.RGBA{255, 255, 255, 255})
	return c.done()
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

func (c *canvas) rect(x, y, w, h float64, col color.RGBA) {
	c.roundRect(x, y, w, h, 0, col)
}

// roundRect fills a rounded rectangle. A fully transparent color cuts a hole instead.
func (c *canvas) roundRect(x, y, w, h, r float64, col color.RGBA) {
	x0, y0 := int(x*ss), int(y*ss)
	x1, y1 := int(math.Ceil((x+w)*ss)), int(math.Ceil((y+h)*ss))
	r *= ss
	cut := col.A == 0
	b := c.img.Bounds()
	for py := max(y0, b.Min.Y); py < min(y1, b.Max.Y); py++ {
		for px := max(x0, b.Min.X); px < min(x1, b.Max.X); px++ {
			fx, fy := float64(px)+0.5, float64(py)+0.5
			dx := math.Max(math.Max(float64(x0)+r-fx, fx-(float64(x1)-r)), 0)
			dy := math.Max(math.Max(float64(y0)+r-fy, fy-(float64(y1)-r)), 0)
			if dx*dx+dy*dy > r*r {
				continue
			}
			if cut {
				c.img.SetRGBA(px, py, color.RGBA{})
			} else {
				c.img.SetRGBA(px, py, over(c.img.RGBAAt(px, py), col))
			}
		}
	}
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
