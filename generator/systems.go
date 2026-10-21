package main

// System is one of the systems EmulationStation can show. Theme is the name ArkOS's es_systems.cfg gives it in
// <theme>, which is what the theme's ${system.theme} expands to, so its channel is art/<Theme>.png.
type System struct {
	Theme string
	Short string // the big text on the channel
	Maker string
	Name  string
	Color uint32 // the channel's color, 0xRRGGBB
}

// The systems ArkOS lists, the same ones portable-helper's systems.go knows, and ES's collections. Each maker's
// systems share a family of colors, so the carousel groups by color as well as by name.
var systems = []System{
	// Nintendo
	{"nes", "NES", "Nintendo", "Entertainment System", 0xC8102E},
	{"famicom", "FC", "Nintendo", "Famicom", 0xB0243A},
	{"fds", "FDS", "Nintendo", "Famicom Disk System", 0xD9A21B},
	{"snes", "SNES", "Nintendo", "Super Nintendo", 0x5B4C9E},
	{"sfc", "SFC", "Nintendo", "Super Famicom", 0x7A5BB5},
	{"snes-hacks", "SFC+", "Nintendo", "Super Famicom Hacks", 0x6A4FA8},
	{"snesmsu1", "MSU-1", "Nintendo", "Super NES CD", 0x4E5BA8},
	{"sufami", "ST", "Nintendo", "SuFami Turbo", 0xD4613A},
	{"satellaview", "BS-X", "Nintendo", "Satellaview", 0x3C8DBC},
	{"sgb", "SGB", "Nintendo", "Super Game Boy", 0x8A5BB0},
	{"n64", "N64", "Nintendo", "Nintendo 64", 0x2A9D48},
	{"n64dd", "64DD", "Nintendo", "Nintendo 64DD", 0x3B7D4F},
	{"gb", "GB", "Nintendo", "Game Boy", 0x6F8A3A},
	{"gbc", "GBC", "Nintendo", "Game Boy Color", 0x8E3BA8},
	{"gba", "GBA", "Nintendo", "Game Boy Advance", 0x3D3FA8},
	{"nds", "DS", "Nintendo", "Nintendo DS", 0x5C6670},
	{"virtualboy", "VB", "Nintendo", "Virtual Boy", 0xA81E1E},
	{"gameandwatch", "G&W", "Nintendo", "Game & Watch", 0xB8862F},
	{"pokemonmini", "MINI", "Nintendo", "Pokémon Mini", 0x2F7FC1},

	// Sega
	{"sg-1000", "SG", "Sega", "SG-1000", 0x2B5FA8},
	{"mastersystem", "SMS", "Sega", "Master System", 0xC0392B},
	{"megadrive", "MD", "Sega", "Mega Drive", 0x1F4E9C},
	{"genesis", "GEN", "Sega", "Genesis", 0x26303D},
	{"segacd", "CD", "Sega", "Sega CD", 0x2F6CB5},
	{"sega32x", "32X", "Sega", "32X", 0xB5372F},
	{"saturn", "SS", "Sega", "Saturn", 0x4A4F6B},
	{"dreamcast", "DC", "Sega", "Dreamcast", 0xE06A1B},
	{"vmu", "VMU", "Sega", "Dreamcast VMU", 0x5E8FC4},
	{"gamegear", "GG", "Sega", "Game Gear", 0x2C3E8F},
	{"naomi", "NAOMI", "Sega", "NAOMI", 0x3471B0},
	{"atomiswave", "AW", "Sammy", "Atomiswave", 0x2B8C8C},

	// Sony
	{"psx", "PS", "Sony", "PlayStation", 0x6B7380},
	{"psp", "PSP", "Sony", "PlayStation Portable", 0x23395B},
	{"pspminis", "MINIS", "Sony", "PSP Minis", 0x34507A},

	// NEC
	{"pcengine", "PCE", "NEC", "PC Engine", 0xE0592A},
	{"pcenginecd", "PCE-CD", "NEC", "PC Engine CD", 0xC94A26},
	{"turbografx", "TG16", "NEC", "TurboGrafx-16", 0x3A3A44},
	{"turbografxcd", "TG-CD", "NEC", "TurboGrafx-CD", 0x4A4A58},
	{"supergrafx", "SGX", "NEC", "SuperGrafx", 0xA8432A},
	{"pcfx", "PC-FX", "NEC", "PC-FX", 0x8A5A3C},
	{"pc98", "PC-98", "NEC", "PC-9801", 0x7A6A55},

	// SNK
	{"neogeo", "NEO", "SNK", "Neo Geo", 0xC9A227},
	{"neogeocd", "NEO-CD", "SNK", "Neo Geo CD", 0xB08A1F},
	{"ngp", "NGP", "SNK", "Neo Geo Pocket", 0x5A5F66},
	{"ngpc", "NGPC", "SNK", "Neo Geo Pocket Color", 0x1F8FA8},

	// Arcade
	{"arcade", "ARC", "Arcade", "Various Platforms", 0x7B2CBF},
	{"mame", "MAME", "Arcade", "MAME 2010", 0x9D4EDD},
	{"mame2003", "MAME", "Arcade", "MAME 2003", 0x8A3FD1},
	{"cps1", "CPS1", "Capcom", "CP System", 0x2160A8},
	{"cps2", "CPS2", "Capcom", "CP System II", 0x1E7A9C},
	{"cps3", "CPS3", "Capcom", "CP System III", 0x1B8F7F},

	// Atari
	{"atari2600", "2600", "Atari", "Atari 2600", 0x8B5A2B},
	{"atari5200", "5200", "Atari", "Atari 5200", 0x5C4033},
	{"atari7800", "7800", "Atari", "Atari 7800", 0x9E3B2E},
	{"atari800", "800", "Atari", "Atari 800", 0xA0703C},
	{"atarixegs", "XEGS", "Atari", "XE Game System", 0x6E6E6E},
	{"atarist", "ST", "Atari", "Atari ST", 0x7A7F87},
	{"atarilynx", "LYNX", "Atari", "Lynx", 0xD4782F},
	{"atarijaguar", "JAG", "Atari", "Jaguar", 0xB52D2D},

	// Bandai, Watara and the other handhelds
	{"wonderswan", "WS", "Bandai", "WonderSwan", 0x7D8A99},
	{"wonderswancolor", "WSC", "Bandai", "WonderSwan Color", 0x2E86C1},
	{"supervision", "SV", "Watara", "Supervision", 0x4F7F4F},
	{"megaduck", "DUCK", "Welback", "Mega Duck", 0xC7A23A},
	{"arduboy", "ARDU", "Arduboy", "Arduboy", 0x2B9E8F},

	// Other consoles
	{"3do", "3DO", "Panasonic", "3DO", 0xB3262E},
	{"coleco", "CV", "Coleco", "ColecoVision", 0x3A3F55},
	{"intellivision", "INTV", "Mattel", "Intellivision", 0x9C7A3C},
	{"odyssey2", "O²", "Magnavox", "Odyssey²", 0x3B6E3B},
	{"videopac", "VP", "Philips", "Videopac", 0x2C6E9C},
	{"vectrex", "VEC", "GCE", "Vectrex", 0x252830},
	{"channelf", "F", "Fairchild", "Channel F", 0xB5762F},
	{"gx4000", "GX", "Amstrad", "GX4000", 0xC03A5A},

	// Computers
	{"amiga", "AMIGA", "Commodore", "Amiga", 0xD25B2A},
	{"amigacd32", "CD32", "Commodore", "Amiga CD32", 0xB34A28},
	{"c64", "C64", "Commodore", "Commodore 64", 0x7E6E57},
	{"c128", "C128", "Commodore", "Commodore 128", 0x6E6455},
	{"c16", "C16", "Commodore", "Commodore 16", 0x55524D},
	{"vic20", "VIC", "Commodore", "VIC-20", 0x9A8A6E},
	{"amstradcpc", "CPC", "Amstrad", "CPC", 0xA83A4A},
	{"msx", "MSX", "Microsoft", "MSX", 0x3D6FB0},
	{"msx2", "MSX2", "Microsoft", "MSX2", 0x2F5C99},
	{"dos", "DOS", "Microsoft", "MS-DOS", 0x3A3A3A},
	{"zxspectrum", "ZX", "Sinclair", "ZX Spectrum", 0x2B2B2B},
	{"zx81", "ZX81", "Sinclair", "ZX81", 0x444444},
	{"x1", "X1", "Sharp", "X1", 0x8C3A3A},
	{"x68000", "X68K", "Sharp", "X68000", 0x4F5A6B},

	// Engines, fantasy consoles and ports
	{"ports", "PORTS", "ArkOS", "Ports", 0x3A7D6B},
	{"doom", "DOOM", "id Software", "Doom", 0x9C2A1A},
	{"scummvm", "SCUMM", "ScummVM", "Adventure Games", 0x4F9A3A},
	{"easyrpg", "RPG", "EasyRPG", "RPG Maker 2000/2003", 0x3A6FA8},
	{"openbor", "BOR", "OpenBOR", "Beats of Rage", 0xA83A2A},
	{"solarus", "SOL", "Solarus", "Action RPGs", 0xD4A23A},
	{"love2d", "LÖVE", "LÖVE", "Lua Games", 0xE0457B},
	{"pico-8", "P8", "Lexaloffle", "PICO-8", 0xFF004D},
	{"tic80", "TIC", "Nesbox", "TIC-80", 0x41A6F6},
	{"wasm4", "W4", "WASM-4", "Fantasy Console", 0x6A5ACD},
	{"lowresnx", "NX", "LowRes", "LowRes NX", 0x3AA87A},
	{"uzebox", "UZE", "Uzebox", "Uzebox", 0x5A3AA8},
	{"retropie", "OPT", "ArkOS", "Options", 0x50566B},

	// EmulationStation's collections
	{"auto-allgames", "ALL", "Collection", "All Games", 0x34BEED},
	{"auto-favorites", "FAV", "Collection", "Favorites", 0xF2B01E},
	{"auto-lastplayed", "RECENT", "Collection", "Last Played", 0x5AB55E},
	{"auto-neverplayed", "NEW", "Collection", "Never Played", 0x9B6FD6},
	{"auto-at2players", "2P", "Collection", "2 Players", 0xE8743B},
	{"auto-at4players", "4P", "Collection", "4 Players", 0xD94F6A},
	{"custom-collections", "MY", "Collection", "Collections", 0x3E8FD8},
}
