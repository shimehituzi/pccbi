package decoder

type ply [][3]int

type Axis int

type Header struct {
	Axis         Axis
	Length, Bias [3]int
}

type bitmap [][]byte
type Voxel []bitmap

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

type Stream struct {
	StartPoints   [][3]uint
	Codes         []uint
	NumCodesArray []uint
}
