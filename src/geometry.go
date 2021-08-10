package main

import (
	"github.com/Tnze/go-mc/save"
)

// Get the top most blocks (visible from the top)
// returns a slice of the namespaced block IDs
func visible_blocks(c save.Column) (vis [16][16]string) {
	sections := sort_subchunks(c.Level.Sections)
	top_index := top_subchunk(sections)
	top := sections[top_index]

	blocks := blocks_in_section(top)

	vis = y_hunter(blocks)
	return
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

// calculate the index of a block inside a subchunk
// given xyz coordinates inside a subchunk
func xyz_to_index(x, y, z int) int {
	return (y * 16 * 16) + (z * 16) + x
}

// iterate over the slice of the blocks in a subchunk
// return a 2d slice of the lowest block in each XZ postition
// also: great function name
func y_hunter(blocks []string) (ret [16][16]string) {
	for x := 0; x < 16; x++ { // iterate over the x axis
		for z := 0; z < 16; z++ { // z axis
			for y := 15; y >= 0; y-- { // y axis from top to bottom
				i := xyz_to_index(x, y, z)
				if b := blocks[i]; b != "minecraft:air" {
					ret[x][z] = b
					break
				}
			}
		}
	}
	return
}

// iterate over a 2d slice of namespaced block IDs
// return false if one isn't populated
func grid_complete(grid [16][16]string) bool {
	for _, x := range grid {
		for _, z := range x {
			if z == "" {
				return false
			}
		}
	}
	return true
}

// get a slice of the blocks in a subchunk
func blocks_in_section(section save.Chunk) (blocks []string) {
	bit_length := index_bit_length(section.Palette)

	for _, v := range section.BlockStates {
		x := nbt_to_block(v, section.Palette, bit_length)
		blocks = append(blocks, x...)
	}
	return
}
