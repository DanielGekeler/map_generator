package main

import (
	"os"
)

const filepath = "region/r.0.0.mca"

func main() {
	pos1 := pos2d{0, 0}
	pos2 := pos2d{511, 511}
	pixelpipe := make(chan mappixel)

	for _, region_path := range needed_regions() {

	}
	raw_region, _ := os.ReadFile(filepath) // fully read a region file => []byte
	chunks := parse_chunks_from_region(raw_region)

	for _, c := range needed_chunks(pos1, pos2) {
		i := calculate_chunk_index(c.X, c.Z)
		go render_chunk(chunks[i], raw_region, pixelpipe, pos1)
	}

	pixels := calc_pixels(pos1, pos2)
	draw_map(pixelpipe, "img/test6.png", pos2.X-pos1.X+1, pos2.Z-pos1.Z+1, pixels)
}

type chunk_meta struct {
	offset  int // chunk data offset in 4KiB sectors in region file
	sectors int // number of sectors
	time    int // last modification time of a chunk in epoch seconds
	x, z    int // x and z chunk coordinates inside the region file

	length int // length of the (compressed) data in bytes
	// 1: GZip (RFC1952) (unused in practice)
	// 2: Zlib (RFC1950) DEFAULT
	// 3: uncompressed (unused in practice)
	compression int
}

// chunk2d is a 2d slice of namespaced block IDs
// used to store a flat slice of a chunk
// or blocks visible from the top
type chunk2d [16][16]string

// pos2d is used to strore a 2 dimensional position
type pos2d struct{ X, Z int }

// describe a single pixel on a map (pos & color)
type mappixel struct{ x, z, color int }

var nilpixel = mappixel{-1, -1, -1}
