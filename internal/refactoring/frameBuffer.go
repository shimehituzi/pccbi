package refactoring

import "sync"

type frameBuffer [][]contour

type contour []chainCode

type chainCode struct {
	start  point
	code   []byte
	points []point
}
type point struct {
	x, y int
}

func NewFrameBuffer(voxel *voxel) {
	frames := NewFrames(voxel)

	fb := make(frameBuffer, len(frames))
	for i := range fb {
		fb[i] = make([]contour, len(frames[i]))
	}

	for f, frame := range frames {
		wg := &sync.WaitGroup{}
		for l, seg := range frame {
			wg.Add(1)
			go func(f, l int, seg segment) {
				fb[f][l] = newContours(seg)
				wg.Done()
			}(f, l, seg)
		}
		wg.Wait()
	}
}

func newContours(s segment) []chainCode {
	return []chainCode{}
}

type segment bitmap
type frame []segment
type frames []frame

func NewFrames(voxel *voxel) frames {
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
