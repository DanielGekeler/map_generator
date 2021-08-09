package main

import (
	"fmt"

	"github.com/Tnze/go-mc/save"
)

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
