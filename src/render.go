package main

import (
	"image"
	"image/png"
	"os"
)

// chn: channel for the pixel data
// filename: path to export the png image
// x,z: size of the image
func draw_map(chn chan mappixel, filename string, x, z int) {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{x, z}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	pixels := x * z

	for v := range chn {
		img.Set(v.x, v.z, rgb_map[v.color])
		pixels--
		if pixels == 0 {
			break
		}
	}

	f, _ := os.Create(filename)
	png.Encode(f, img)
	f.Close()
}
