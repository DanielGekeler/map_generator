package main

import "image/color"

type mappixel struct {
	c    color.Color
	x, z int
}

func render_chunk(chn chan mappixel, chunk chunk_meta, region []byte) {
	c := load_chunk(chunk, region)
	vis := visible_blocks(c)

	for x, a := range vis {
		for z, b := range a {
			color := color_id[b]
			rgb := rgb_map[color]

			chn <- mappixel{rgb, x, z}
		}
	}
}
