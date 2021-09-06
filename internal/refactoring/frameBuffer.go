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

func NewFrameBuffer(voxel *voxel) frames {
	return newFrames(voxel)
}

type segment bitmap
type frame []segment
type frames []frame

func newFrames(voxel *voxel) frames {
	lv, numLabels := NewLabels(voxel)

	frames := make(frames, voxel.header.length[0])
	for i := range frames {
		frames[i] = make(frame, numLabels[i])
		for j := range frames[i] {
			frames[i][j] = make(segment, voxel.header.length[1])
			for k := range frames[i][j] {
				frames[i][j][k] = make([]byte, voxel.header.length[2])
			}
		}
	}

	for f := range lv {
		for y := range lv[f] {
			for x, label := range lv[f][y] {
				if label != 0 {
					l := label - 1
					frames[f][l][y][x] = 1
				}
			}
		}
	}

	return frames
}
