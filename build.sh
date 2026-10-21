#!/bin/sh
# Draws the theme's images and zips the theme into dist/es-theme-cartridge.zip, ready to install with
# portable-helper (Themes → Add theme) or to unzip into the card's themes folder.
set -e
cd "$(dirname "$0")"
(cd generator && go run .)
rm -rf dist
mkdir -p dist/es-theme-cartridge
cp -R theme/ dist/es-theme-cartridge/
(cd dist && zip -qr -X es-theme-cartridge.zip es-theme-cartridge && rm -rf es-theme-cartridge)
echo dist/es-theme-cartridge.zip
