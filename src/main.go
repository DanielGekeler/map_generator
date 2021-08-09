package main

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/Tnze/go-mc/save"
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

func split_bytes(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:])
	}
	return chunks
}

func bytes_to_int(input []byte) int {
	return int(binary.BigEndian.Uint32(input))
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

// decompress_chunk takes the specific part of the region and returns a []byte of the raw NBT data
/*func decompress_chunk(chunk chunk_meta, region []byte) []byte {
	a := (chunk.offset * 4096) + 5
	b := chunk.length

	compressed := bytes.NewBuffer(region[a : a+b])
	r, _ := zlib.NewReader(compressed)
	raw_nbt, _ := ioutil.ReadAll(r)
	r.Close()
	return raw_nbt
}*/

// Get the top most blocks (visible from the top)
// returns a slice of the namespaced block IDs
func visible_blocks(c save.Column) []string {
	sections := sort_subchunks(c.Level.Sections)
	top_index := top_subchunk(sections)
	top := sections[top_index]
	bit_length := index_bit_length(top.Palette)

	var blocks []string
	for _, v := range top.BlockStates {
		x := nbt_to_block(v, top.Palette, bit_length)
		blocks = append(blocks, x...)
	}

	fmt.Println(blocks)

	return nil
}

// Sort subchunks (16x16x16) by Y index
func sort_subchunks(sections []save.Chunk) []save.Chunk {
	ret := make([]save.Chunk, len(sections))
	for _, v := range sections {
		if v.Palette != nil {
			ret[v.Y] = v
		}
	}
	return ret[:]
}

// Find the highest subchunk with blocks other than air
func top_subchunk(sections []save.Chunk) int {
	for i := len(sections) - 1; i >= 0; i-- {
		if len(sections[i].Palette) > 1 {
			return i
		}
	}
	return 0
}
