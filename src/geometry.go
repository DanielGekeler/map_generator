package main

import (
	"fmt"
	"image"
	"math"

	"github.com/Tnze/go-mc/save"
)

// Get the top most blocks (visible from the top)
// returns a slice of the namespaced block IDs
func visible_blocks(c save.Column) (vis chunk2d) {
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
// return a chunk2d object of the lowest block in each XZ postition
// also: great function name
func y_hunter(blocks []string) (ret chunk2d) {
	for x := 0; x < 16; x++ { // iterate over the x axis
		for z := 0; z < 16; z++ { // z axis
			for y := 15; y >= 0; y-- { // y axis from top to bottom
				i := xyz_to_index(x, y, z)
				if b := blocks[i]; !is_transparent(b) {
					ret[x][z] = b
					break
				}
			}
		}
	}
	return
}

// iterate over a chunk2d and return a list of missing positions
func find_missing(grid chunk2d) (pos [][2]int) {
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

// Recursive function that searches missing blocks in a block2d object
// each iteration of add_missing searches a lower subchunk then the one before it
// until it is complete or the bottom of the world is reached
// resulting in a chunk2d object with all blocks visible from the top
func add_missing(blocks chunk2d, sections []save.Chunk, index int) chunk2d {
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
			if !is_transparent(new_blocks[i]) {
				// store the block that got found
				blocks[x][z] = new_blocks[i]
				// exit the Y loop (search for the next missing block)
				break
			}
		}
	}
	return add_missing(blocks, sections, index)
}

// calculate the chunk coordinates of all chunks
// within a given area. pos1,pos2 are block coordinates
// pos1 needs to be north west
// pos2 needs to be south east
// returns a slice of chunk coordinates
func needed_chunks(pos1, pos2 pos2d) (ret []pos2d) {
	x1 := int(math.Floor(float64(pos1.X) / 16.0))
	z1 := int(math.Floor(float64(pos1.Z) / 16.0))

	x2 := int(math.Floor(float64(pos2.X) / 16.0))
	z2 := int(math.Floor(float64(pos2.Z) / 16.0))

	for i := x1; i <= x2; i++ {
		for u := z1; u <= z2; u++ {
			ret = append(ret, pos2d{i, u})
		}
	}
	return
}

/*func needed_regions(chunks []pos2d) (ret []pos2d) {
	a := pos2d{chunks[0].X, chunks[0].Z}
	l := len(chunks) - 1
	b := pos2d{chunks[l].X, chunks[l].Z}

	x1 := int(math.Floor(float64(a.X) / 32.0))
	z1 := int(math.Floor(float64(a.Z) / 32.0))

	x2 := int(math.Floor(float64(b.X) / 32.0))
	z2 := int(math.Floor(float64(b.Z) / 32.0))

	for i := x1; i <= x2; i++ {
		for u := z1; u <= z2; u++ {
			ret = append(ret, pos2d{i, u})
		}
	}
	return
}*/

// get the filenames of all region files inside area
func needed_regions(area image.Rectangle) (files []string) {
	a := block_pos_to_region(point_to_pos2d(area.Min))
	b := block_pos_to_region(point_to_pos2d(area.Max))

	for i := a.X; i <= b.X; i++ { // iterate over x
		for h := a.Z; h <= b.Z; h++ { // iterate over z
			files = append(files, fmt.Sprintf("r.%v.%v.mca", i, h))
		}
	}
	return
}

// calculate in which region file a given block coordinates is
func block_pos_to_region(block pos2d) pos2d {
	// x axis
	x := block.X / 512
	if block.X < 0 { // fix of by one error for negative numbers
		x -= 1
	}

	// z axis
	z := block.Z / 512
	if block.Z < 0 {
		z -= 1
	}
	return pos2d{x, z}
}

// calculate in which chunk, a given block is
// block pos => chunk pos
func block_pos_to_chunk(block pos2d) pos2d {
	x := int(math.Floor(float64(block.X) / 16.0))
	z := int(math.Floor(float64(block.Z) / 16.0))
	return pos2d{x, z}
}
