package decoder

type bitmap [][]byte

type ply [][3]int
type order [3]int
type orderString int

type Voxel []bitmap
type VoxelHeader struct {
	axis, length, bias [3]int
}

type label [][]int
type labeledVoxel []label

type segment bitmap
type frame []segment
type frames []frame

type point struct {
	x, y int
}

type chainCode struct {
	start  point
	code   []byte
	points []point
}

type direction struct {
	d    point
	code byte
}

type contourBuffer [][]contour
type contour []chainCode

type StreamHeader []int
type Stream struct {
	StartPoints []int
	Flags       []byte
	Codes       []byte
}
