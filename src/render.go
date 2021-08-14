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

func render_chunk(chunk chunk_meta, region []byte, chn chan mappixel, min, max pos2d) {
	c := load_chunk(chunk, region)
	vis := visible_blocks(c)

	for xi, x := range vis {
		stopx := xi + (chunk.x * 16)
		if stopx < min.X || stopx > max.X {
			continue
		}
		for zi, z := range x {
			stopz := zi + (chunk.z * 16)
			if stopz < min.Z || stopz > max.Z {
				continue
			}

			color := color_id[z]

			a := xi + (16 * chunk.x) - min.X
			b := zi + (16 * chunk.z) - min.Z
			chn <- mappixel{a, b, color}
		}
	}
}
