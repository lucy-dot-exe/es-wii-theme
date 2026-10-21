# es-theme-cartridge

**Cartridge** is an EmulationStation theme for ArkOS on 640×480 handhelds like the **R36S**, made to be
added to [Jetup13's EmulationStation OGA Theme Gallery](https://github.com/Jetup13/Emulationstation-OGA-Theme-Gallery)
like the themes already there.

It looks like the Wii's menu: a pinstriped white background with rows of channel slots, each system a glossy
channel in its own color in the middle row, a blue frame around the selected one, and the curved dock along the
bottom with the clock. Each maker's systems share a family of colors. A game list has a band in the system's color
across the top, the games in a white panel on the left with a light blue selector, and the box art and description
in a panel on the right.

| Systems | Game list |
| --- | --- |
| ![](previews/system.png) | ![](previews/gamelist.png) |

The previews are screenshots of ArkOS's EmulationStation built for macOS (see [es-preview](../es-preview/)), not
from the console.

## Installing

Run `./build.sh`. It makes `dist/es-theme-cartridge.zip`. Then either:

- add the zip with [portable-helper](../portable-helper/) (**Themes → Add theme**), or
- unzip it into the `themes` folder on the card's `EASYROMS` partition.

On the console, pick it in **Start → UI Settings → Theme Set**.

## How it's made

| Folder | What's in it |
| --- | --- |
| [theme/](theme/) | The theme, as it goes on the card: `theme.xml`, the fonts and the images |
| [generator/](generator/) | A Go program that draws every image in `theme/art/` and `theme/assets/` |
| [previews/](previews/) | The screenshots above |

EmulationStation's SVG renderer can't draw text, so the channels are PNGs with the text drawn on them. They are
drawn from the table in [generator/systems.go](generator/systems.go): one row per system ArkOS lists (the same
systems as portable-helper's `systems.go`) and per ES collection (all games, favorites, last played, ...), with its
short name, maker, full name and color.

To change the art, change the generator and run it (`cd generator && go run .`). The images are committed, so
the theme works straight from the repository. The colors are in both `generator/main.go` and `theme/theme.xml`.

- `art/<system>.png`: the system's channel in the carousel (384×196, twice its size at 640×480). `<system>` is
  the name in `<theme>` in ArkOS's `es_systems.cfg`. A system without a channel shows its name instead.
- `art/<system>-header.png`: the band across the top of the system's game list (640×60, with its shadow).
- `assets/menu.png`: the system view's background, with the empty slots and the dock. It and `assets/frame.png`
  (the blue frame) are `extra="static"`, so they stay still while the channels slide.
- `assets/gamelist.png` (with the two panels) and `assets/grid.png`: the game lists' backgrounds.
- `assets/selector.png` (the pill behind the selected game) and `assets/noart.png` (for games without art).

The carousel shows three channels a row and the empty slots are drawn in the background at the same places, so
theme.xml's carousel and frame and the generator's layout constants have to change together.

## Credits

The font is [M PLUS Rounded 1c](https://fonts.google.com/specimen/M+PLUS+Rounded+1c) by the Rounded M+ Project,
under the SIL Open Font License ([theme/fonts/OFL.txt](theme/fonts/OFL.txt)), cut down to its Latin letters and
symbols. Everything else is drawn by the generator. The channels use system names, never the makers' logos, and
nothing is taken from the Wii itself.
