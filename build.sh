#!/bin/sh
# Draws the theme's images and animations (it needs ffmpeg) and zips the theme into dist/es-wii-theme.zip, ready to
# unzip into the card's themes folder.
set -e
cd "$(dirname "$0")"
(cd generator && go run .)
rm -rf dist
mkdir -p dist/es-wii-theme
cp -R theme/ dist/es-wii-theme/
(cd dist && zip -qr -X es-wii-theme.zip es-wii-theme && rm -rf es-wii-theme)
echo dist/es-wii-theme.zip
