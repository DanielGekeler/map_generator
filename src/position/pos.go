package position

import (
	"errors"
	"fmt"
)

type scale int

const (
	Block  scale = 1   // Blocks are 1 block wide
	Chunk  scale = 16  // Chunks are 16 blocks wide
	Region scale = 512 // Regions are 512 blocks wide
)

var (
	errConversion = errors.New("INTERNAL ERROR: cannot convert to lower grid")
)

// Pos stores a 2d position on a map / world
type Pos struct {
	X, Z int // coordinates

	/* The Grid of the Position
	this allows to store block positions and
	chunk positions in the same struct type
	and easyily convert between them.
	Grid should never change!*/
	Grid scale
}

// Get the the coordinates of a Pos{} in a higher grid
// Example: In which region is chunk -23 51? {-1 1}
func (p Pos) InGrid(g scale) (r Pos) {
	if p.Grid == g {
		return p
	}

	if g < p.Grid {
		panic(fmt.Errorf("%w (%d -> %d)", errConversion, p.Grid, g))
	}

	f := int(g / p.Grid)

	r.X = scalefloor(p.X, f)
	r.Z = scalefloor(p.Z, f)
	r.Grid = g
	return
}

// internal function for InGrid()
func scalefloor(a, b int) (c int) {
	c = a / b
	if a < 0 {
		c--
	}
	return
}

// Get the offset of p to the begin of the closest grid unit g
// This can be used to find the relative position of a chunk in a region file
func (p Pos) Offset(g scale) (x, z int) {
	f := int(g / p.Grid)

	x = p.X % f
	if p.X < 0 {
		x += f - 1
	}

	z = p.Z % f
	if p.Z < 0 {
		z += f - 1
	}
	return
}
