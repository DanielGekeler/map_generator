package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"

	"github.com/Tnze/go-mc/save"
)

const filepath = "region/r.0.0.mca"

func main() {
	fmt.Println("Starting")

	raw_region, _ := os.ReadFile(filepath) // fully read a region file => []byte

	chunks := parse_chunks_from_region(raw_region)

	chunk := chunks[0]

	visible_blocks(chunk, raw_region)
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

func parse_chunks_from_region(region []byte) []chunk_meta {
	locations := split_bytes(region[:4096], 4)
	time_data := split_bytes(region[4096:8192], 4)

	ret := make([]chunk_meta, 1024)

	for i := 0; i < len(ret); i++ {
		pos := append([]byte{0}, locations[i][:3]...)
		sector_length := locations[i][3]

		offset := bytes_to_int(pos)
		chunk_time := bytes_to_int(time_data[i])

		x, z := calculate_chunk_pos(i)

		length := bytes_to_int(region[offset*4096 : offset*4096+4])
		compression := int(region[offset*4096+4])

		ret[i] = chunk_meta{offset: offset,
			sectors: int(sector_length),
			time:    chunk_time,
			x:       x, z: z,
			length:      length,
			compression: compression,
		}
	}

	return ret
}

func calculate_chunk_pos(index int) (int, int) {
	return index % 32, index / 32
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

// Calculate how many bits are needed to index the elements in the pallete of a chunk section
func index_bit_length(palette []save.Block) int {
	bits := 4
	for math.Pow(2.0, float64(bits)) < float64(len(palette)) {
		bits++
	}
	return bits
}

// Parse namespaced block IDs from nbt data
// Block IDs in each section are asigned an index by the pallete
// multiple indexes are stored in one int64
// bit_length is the number of bits needed for an index number in the pallete
// this function returns a []string of the namespaced block IDs
func nbt_to_block(long int64, pallete []save.Block, bit_length int) (block_id []string) {
	// mask works like a subnetmask to get the last (bit_length) bits
	mask := byte(math.Pow(2.0, float64(bit_length)) - 1)

	// shift the input number by a diffrent amount
	// with each iteration to get all indexes
	for i := 0; i < 64/bit_length; i++ {
		shifted := long >> (i * bit_length) // shifting
		block := shifted & int64(mask)      // get only the last (bit_length) bits
		block_id = append(block_id, pallete[block].Name)
	}
	return
}

// Get the top most blocks (visible from the top)
// returns a slice of the namespaced block IDs
func visible_blocks(chunk chunk_meta, region []byte) []string {
	// calculate offsets
	a := (chunk.offset * 4096) + 4
	b := chunk.length

	data := region[a : a+b] // the raw bytes of the chunk data

	var c save.Column // Column means the whole chunk (0-255)...
	if err := c.Load(data); err != nil {
		panic(err)
	}

	sections := sort_subchunks(c.Level.Sections)
	top_index := top_subchunk(sections)
	top := sections[top_index]

	bit_length := index_bit_length(top.Palette)
	for _, v := range top.BlockStates {
		x := nbt_to_block(v, top.Palette, bit_length)
		fmt.Println(x)
	}

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
