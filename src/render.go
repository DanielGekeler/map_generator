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

func render_chunk(chunk chunk_meta, region []byte, chn chan mappixel, begin pos2d) {
	c := load_chunk(chunk, region)
	vis := visible_blocks(c)

	x_off := (16 * chunk.x) - begin.X
	z_off := (16 * chunk.z) - begin.Z

	for xi, x := range vis {
		for zi, z := range x {
			color := color_id[z]

			a := xi + x_off
			b := zi + z_off
			chn <- mappixel{a, b, color}
		}
	}
}
