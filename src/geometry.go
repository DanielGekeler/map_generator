package main

import (
	"github.com/Tnze/go-mc/save"
)

// Get the top most blocks (visible from the top)
// returns a slice of the namespaced block IDs
func visible_blocks(c save.Column) (vis [16][16]string) {
	sections := sort_subchunks(c.Level.Sections)
	index := top_subchunk(sections)
	top := sections[index]

	blocks := blocks_in_section(top)
	vis = y_hunter(blocks)

	vis = add_missing(vis, sections, index)

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

// iterate over a grid of blocks and return a list of missing positions
func find_missing(grid [16][16]string) (pos [][2]int) {
	for xi, x := range grid {
		for zi, z := range x {
			if z == "" {
				pos = append(pos, [2]int{xi, zi})
			}
		}
	}
	return
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

// Recursive function that searches missing blocks in a [][]string of namespaced block IDs
// each iteration of add_missing searches a lower subchunk then the one before it
// until it is complete or the bottom of the world is reached
func add_missing(blocks [16][16]string, sections []save.Chunk, index int) [16][16]string {
	missing := find_missing(blocks)
	if len(missing) == 0 || index == 0 { // return if complete
		return blocks
	}
	index -= 1
	new_section := sections[index]
	new_blocks := blocks_in_section(new_section)

	// iterate over the list of missing blocks
	for _, v := range missing {
		// x and z coordinates (realtive to chunk border)
		// of the missing block
		x := v[0]
		z := v[1]

		// iterate over the Y axis to find the highest non air block
		for y := 15; y >= 0; y-- {
			i := xyz_to_index(x, y, z)
			if new_blocks[i] != "minecraft:air" {
				// store the block that got found
				blocks[x][z] = new_blocks[i]
				// exit the Y loop (search for the next missing block)
				break
			}
		}
	}
	return add_missing(blocks, sections, index)
}
