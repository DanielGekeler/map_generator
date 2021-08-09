package main

import (
	"math"

	"github.com/Tnze/go-mc/save"
)

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

func load_chunk(chunk chunk_meta, region []byte) save.Column {
	// calculate offsets
	a := (chunk.offset * 4096) + 4
	b := chunk.length

	data := region[a : a+b] // the raw bytes of the chunk data

	var c save.Column // Column means the whole chunk (0-255)...
	if err := c.Load(data); err != nil {
		panic(err)
	}
	return c
}
