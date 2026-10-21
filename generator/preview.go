package main

import (
	"image"

	xdraw "golang.org/x/image/draw"
)

// The previews are mock-ups of the theme at the R36S's 640×480, laid out with the same numbers as theme.xml, so
// a change can be checked without a console. They're drawn by hand, not by EmulationStation: screenshots from the
// console are still what the gallery wants.

const screenW, screenH = 640, 480

// image draws img into the box (x, y, w, h), in canvas units.
func (c *canvas) image(img image.Image, x, y, w, h float64) {
	r := image.Rect(int(x*ss), int(y*ss), int((x+w)*ss), int((y+h)*ss))
	xdraw.CatmullRom.Scale(c.img, r, img, img.Bounds(), xdraw.Over, nil)
}

// previewSystem is the system view, with cur selected between prev and next.
func previewSystem(prev, cur, next System) image.Image {
	c := newCanvas(screenW, screenH)
	c.image(background(screenW, screenH), 0, 0, screenW, screenH)
	c.text("12:00", "Lato-Bold", screenH*0.035, screenW*0.97, screenH*0.06, 0, withAlpha(inkDim, 255), alignRight)

	// <carousel>: pos 0 0.20, size 1 0.50, logoSize 0.28 0.175, logoScale 1.3, maxLogoCount 3.
	lw, lh := screenW*0.28, screenH*0.175
	cy := screenH * (0.20 + 0.50/2)
	step := screenW / 3.0
	for i, s := range []System{prev, next} {
		cx := screenW/2 + step*float64(i*2-1)
		c.image(card(s, 480, 300), cx-lw/2, cy-lh/2, lw, lh)
	}
	c.image(card(cur, 480, 300), screenW/2-lw*1.3/2, cy-lh*1.3/2, lw*1.3, lh*1.3)

	// <text name="systemInfo">: pos 0 0.70, size 1 0.06.
	c.text("42 GAMES AVAILABLE", "Lato-Bold", screenH*0.04, screenW/2, screenH*0.745, 0.08, inkDim, alignCenter)
	helpBar(c)
	return c.done()
}

// previewGamelist is a system's detailed game list, with no art scraped for the selected game.
func previewGamelist(s System) image.Image {
	c := newCanvas(screenW, screenH)
	c.image(background(screenW, screenH), 0, 0, screenW, screenH)
	c.image(header(s, screenW, 56), 0, 0, screenW, 56)
	c.text(s.Maker+" "+s.Name, "Lato-Black", screenH*0.055, screenW*0.04, 56*0.5+screenH*0.055*0.36, 0, ink, alignLeft)

	// <textlist name="gamelist">: pos 0.03 0.15, size 0.53 0.75, fontSize 0.042, lineSpacing 1.6.
	games := []string{"Astro Courier", "Bramble Knights", "Cobalt Circuit", "Deep Mine Dash", "Ember Valley",
		"Frost Relay", "Glimmer Tactics", "Harbor Lights", "Iron Orchard", "Jade Runner", "Kite Tournament"}
	size := screenH * 0.042
	row := size * 1.6
	x, y, w := screenW*0.03, screenH*0.15, screenW*0.53
	for i, g := range games {
		top := y + float64(i)*row
		col := ink
		if i == 2 {
			c.roundRect(x, top, w, row, row*0.3, accent)
			col = bgBottom
		}
		c.text(g, "Lato-Bold", size, x+screenW*0.015, top+row*0.5+size*0.36, 0, col, alignLeft)
	}

	// <image name="md_image">: origin 0.5 0.5, pos 0.785 0.40, maxSize 0.38 0.44 (the placeholder is square).
	side := screenH * 0.44
	c.image(noArt(300, 300), screenW*0.785-side/2, screenH*0.40-side/2, side, side)

	// <text name="md_description">: pos 0.60 0.66, size 0.37 0.25, fontSize 0.032.
	desc := []string{"Deliver parcels across a", "tiny solar system before", "the sun goes down."}
	for i, l := range desc {
		c.text(l, "Lato-Regular", screenH*0.032, screenW*0.60, screenH*0.69+float64(i)*screenH*0.045, 0, inkDim, alignLeft)
	}
	helpBar(c)
	return c.done()
}

// helpBar stands in for EmulationStation's help prompts: <helpsystem> pos 0.03 0.935.
func helpBar(c *canvas) {
	x := screenW * 0.03
	for _, p := range [][2]string{{"A", "SELECT"}, {"B", "BACK"}, {"START", "MENU"}} {
		bw := textWidth(p[0], "Lato-Black", screenH*0.028, 0) + 10
		c.roundRect(x, screenH*0.935, bw, screenH*0.045, 6, withAlpha(inkDim, 255))
		c.text(p[0], "Lato-Black", screenH*0.028, x+bw/2, screenH*0.968, 0, bgBottom, alignCenter)
		x += bw + 6
		c.text(p[1], "Lato-Bold", screenH*0.03, x, screenH*0.97, 0.05, withAlpha(inkDim, 255), alignLeft)
		x += textWidth(p[1], "Lato-Bold", screenH*0.03, 0.05) + 16
	}
}
