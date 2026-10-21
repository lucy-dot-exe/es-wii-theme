package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"os/exec"

	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/math/f64"
)

// The channels' animations, like the Wii's channels that never stand still: the picture bobs gently and, once a
// loop, a shine sweeps across the face. theme.xml plays the selected system's over its channel. The first frame is
// the still channel, so the video fades in over it without a jump.

const (
	animFPS     = 30
	animSeconds = 3
	bob         = 3.5 // how far the picture moves up and down, in pixels of the double-size channel
)

// frame draws the channel at t seconds into its loop.
func (l layers) frame(t float64) *image.RGBA {
	out := image.NewRGBA(l.base.Bounds())
	draw.Draw(out, out.Bounds(), l.base, image.Point{}, draw.Src)
	dy := -bob * math.Sin(2*math.Pi*t/animSeconds)
	xdraw.CatmullRom.Transform(out, f64.Aff3{1, 0, 0, 0, 1, dy}, l.icon, l.icon.Bounds(), xdraw.Over, nil)
	draw.Draw(out, out.Bounds(), l.top, image.Point{}, draw.Over)
	l.shine(out, t)
	return out
}

// shine is a soft diagonal band of light crossing the face between 1 and 1.8 seconds into the loop.
func (l layers) shine(img *image.RGBA, t float64) {
	const start, length, width = 1.0, 0.8, 26.0
	if t < start || t > start+length {
		return
	}
	b := img.Bounds()
	p := (t - start) / length
	center := -80 + p*(float64(b.Dx())+200) // where the band crosses y = 0, moving right
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			fx, fy := float64(x)+0.5, float64(y)+0.5
			if !l.face(fx, fy) {
				continue
			}
			d := (fx + 0.45*fy - center) / width
			a := math.Exp(-d*d) * 0.42
			if a < 0.004 {
				continue
			}
			c := img.RGBAAt(x, y)
			f := func(v uint8) uint8 { return uint8(float64(v) + (255-float64(v))*a) }
			img.SetRGBA(x, y, color.RGBA{f(c.R), f(c.G), f(c.B), c.A})
		}
	}
}

// animate encodes the channel's loop to path as an H.264 MP4, which ES plays with libVLC.
func animate(l layers, path string) {
	b := l.base.Bounds()
	cmd := exec.Command("ffmpeg", "-y", "-loglevel", "error",
		"-f", "rawvideo", "-pix_fmt", "rgba", "-s", fmt.Sprintf("%dx%d", b.Dx(), b.Dy()), "-r", fmt.Sprint(animFPS), "-i", "-",
		"-an", "-c:v", "libx264", "-preset", "slow", "-crf", "16", "-pix_fmt", "yuv420p", "-movflags", "+faststart", path)
	in, err := cmd.StdinPipe()
	check(err)
	check(cmd.Start())
	for i := 0; i < animFPS*animSeconds; i++ {
		f := l.frame(float64(i) / animFPS)
		// A video has no transparency, so the channel's rounded corners get the menu's background gray. (Rounding
		// the video in theme.xml instead would stop ES drawing it, as with the grid's tiles.)
		for p := 0; p < len(f.Pix); p += 4 {
			if a := f.Pix[p+3]; a < 255 {
				k := 255 - int(a)
				f.Pix[p], f.Pix[p+1], f.Pix[p+2], f.Pix[p+3] = uint8(int(f.Pix[p])+k*0xF4/255), uint8(int(f.Pix[p+1])+k*0xF4/255), uint8(int(f.Pix[p+2])+k*0xF4/255), 255
			}
		}
		_, err := in.Write(f.Pix)
		check(err)
	}
	check(in.Close())
	check(cmd.Wait())
}
