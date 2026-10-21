# es-theme-cartridge

**Cartridge** is an EmulationStation theme for ArkOS on 640×480 handhelds like the **R36S**, made to be
added to [Jetup13's EmulationStation OGA Theme Gallery](https://github.com/Jetup13/Emulationstation-OGA-Theme-Gallery)
like the themes already there.

Each system is a cartridge in its own color, with its name on the label. Each maker's systems share a family of
colors. A game list has a header in the system's color, the games on the left and the box art and description on
the right.

| Systems | Game list |
| --- | --- |
| ![](previews/system.png) | ![](previews/gamelist.png) |

The previews are mock-ups drawn by the generator with theme.xml's layout, not screenshots from the console.

## Installing

Run `./build.sh`. It makes `dist/es-theme-cartridge.zip`. Then either:

- add the zip with [portable-helper](../portable-helper/) (**Themes → Add theme**), or
- unzip it into the `themes` folder on the card's `EASYROMS` partition.

On the console, pick it in **Start → UI Settings → Theme Set**.

## How it's made

| Folder | What's in it |
| --- | --- |
| [theme/](theme/) | The theme, as it goes on the card: `theme.xml`, the fonts and the images |
| [generator/](generator/) | A Go program that draws every image in `theme/art/` and `theme/assets/`, and the previews |
| [previews/](previews/) | The mock-ups above |

EmulationStation's SVG renderer can't draw text, so the system cards are PNGs with the text drawn on them. They
are drawn from the table in [generator/systems.go](generator/systems.go): one row per system ArkOS lists, with
its short name, maker, full name and color. The same systems as portable-helper's `systems.go`.

To change the art, change the generator and run it (`cd generator && go run .`). The images are committed, so
the theme works straight from the repository. The colors are in both `generator/main.go` and `theme/theme.xml`.

- `art/<system>.png`: the system's card in the carousel (480×300). `<system>` is the name in `<theme>` in
  ArkOS's `es_systems.cfg`. A system without a card shows its name instead.
- `art/<system>-header.png`: the strip across the top of the system's game list (640×56).
- `assets/background.png`, `assets/noart.png` (for games without art), `assets/selector.png` (the highlight
  behind the selected game, tinted by theme.xml).

## Credits

The font is [Lato](https://www.latofonts.com/) by Łukasz Dziedzic, under the SIL Open Font License
([theme/fonts/OFL.txt](theme/fonts/OFL.txt)). Everything else is drawn by the generator. The cards use system
names, never the makers' logos.
