package grid

type scale int

const (
	Chunk  scale = 16  // Chunks are 16 blocks wide
	Region scale = 512 // regions are 512 blocks wide
)

// Calculate the position of a block in a "higher grid"
// this function is 1d!
// Example: in which chunk is X:1337 ?
// InGrid(1337, grid.Chunk) will return 83
func AbsGrid(p int, s scale) int {
	a := p / int(s)
	if p < 0 { // fix of by one error for negative numbers
		a--
	}
	return a
}

// TODO: Add function RelGrid to calculate relative position inside the grid
// (Distance to the start of the grid unit a point resides in)
