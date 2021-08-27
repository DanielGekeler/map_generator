package main

import (
	"image"
	"image/png"
	"math"
	"os"
)

// chn: channel for the pixel data
// filename: path to export the png image
// x,z: size of the image
func draw_map(chn chan mappixel, filename string, x, z, pixels int) {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{x, z}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	//pixels := x * z

	for v := range chn {
		if v == nilpixel {
			pixels -= 256
			continue
		}
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
	c, err := load_chunk(chunk, region)
	if err != nil || len(c.Level.Sections) == 0 {
		chn <- nilpixel
		return
	}
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

func calc_pixels(a, b pos2d) int {
	a0 := block_pos_to_chunk(a)
	b0 := block_pos_to_chunk(b)
	return (b0.X - a0.X + 1) * 16 * (b0.Z - a0.Z + 1) * 16
}

// create a image.Rectangle from 2 specified corners
// given corners can be on any side as long as they opose each other
func image_area(a, b pos2d) (area image.Rectangle) {
	min_x := math.Min(float64(a.X), float64(b.X))
	min_z := math.Min(float64(a.Z), float64(b.Z))
	area.Min = image.Point{int(min_x), int(min_z)}

	max_x := math.Max(float64(a.X), float64(b.X))
	max_z := math.Max(float64(a.Z), float64(b.Z))
	area.Max = image.Point{int(max_x), int(max_z)}
	return
}
