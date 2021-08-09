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

	chunk := chunks[0]
	c := load_chunk(chunk, raw_region)
	visible_blocks(c)
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
