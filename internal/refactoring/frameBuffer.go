package refactoring

type frameBuffer [][]contour
type contour struct {
	outer chainCode
	inner []chainCode
}
type chainCode struct {
	start  point
	code   []byte
	points []point
}
type point struct {
	x, y int
}

func NewFrameBuffer() {
	// voxel から frames を生成

	// segment の外輪郭と内輪郭を入手
}

type segment bitmap
type frame []segment
type frames []frame

func newFrames() {
	// Voxel をフレーム毎にラベリングして各セグメント毎に img を作成
}
