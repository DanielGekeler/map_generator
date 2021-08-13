package main

import (
	"fmt"
	"os"
)

const filepath = "region/r.0.0.mca"

func main() {
	fmt.Println("Starting")

	raw_region, _ := os.ReadFile(filepath) // fully read a region file => []byte
	chunks := parse_chunks_from_region(raw_region)

	chunk := chunks[34]
	c := load_chunk(chunk, raw_region)
	vis := visible_blocks(c)
	fmt.Println(chunk.x, chunk.z)
	for _, v := range vis[15] {
		fmt.Println(v)
	}
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
