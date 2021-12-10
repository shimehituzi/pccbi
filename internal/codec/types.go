package codec

type Ply [][3]int

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

type chaincode struct {
	Start  point
	Code   []byte
	Points []point
}

type direction struct {
	d    point
	code byte
}

type Contour [][][]chaincode

type Stream struct {
	StartPoints   [][3]uint
	Codes         []uint
	NumCodesArray []uint
}
