package codec

import "sync"

func EncContour(voxel Voxel, header *Header) Contour {
	frames := NewFrames(voxel, header)

	contour := make(Contour, len(frames))
	for i := range contour {
		contour[i] = make([][]chaincode, len(frames[i]))
	}

	for f, frame := range frames {
		for l, img := range frame {
			contour[f][l] = newContourSegment(bitmap(img))
		}
	}

	return contour
}

func NewFrames(voxel Voxel, header *Header) frames {
	lv, numLabels := newLabels(voxel, header)

	frames := make(frames, header.Length[0])
	for i := range frames {
		frames[i] = make(frame, numLabels[i])
		for j := range frames[i] {
			frames[i][j] = make(segment, header.Length[1])
			for k := range frames[i][j] {
				frames[i][j][k] = make([]byte, header.Length[2])
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

func newLabels(voxel Voxel, header *Header) (labeledVoxel, []int) {
	lv := make(labeledVoxel, header.Length[0])
	numLabels := make([]int, header.Length[0])
	wg := &sync.WaitGroup{}
	for i := range lv {
		wg.Add(1)
		go func(i int) {
			lv[i], numLabels[i] = newLabel(voxel[i])
			wg.Done()
		}(i)
	}
	wg.Wait()

	return lv, numLabels
}

func newContourSegment(orig bitmap) []chaincode {
	img := make(bitmap, len(orig))
	for i := range orig {
		img[i] = make([]byte, len(orig[i]))
		copy(img[i], orig[i])
	}

	cs := []chaincode{}

	// 外輪郭
	outer := contourTracking(img, 1, false)
	if outer == nil {
		panic("cannot produce chaincode")
	}
	cs = append(cs, *outer)
	outerPoints := outer.getPoints()

	// 塗り潰し
	filledOutside := false
	label := byte(2)

	for y := range img {
		for x, v := range img[y] {
			if v == 0 {
				p := point{x, y}
				if filledOutside {
					fillArea(img, p, 0, label)
					label++
				} else {
					if p.isInside(outerPoints) {
						fillArea(img, p, 0, label)
						label++
					} else {
						fillArea(img, p, 0, 1)
						filledOutside = true
					}
				}
			}
		}
	}

	// 内輪郭
	for l := byte(2); l < label; l++ {
		inner := contourTracking(img, l, true)
		if outer == nil {
			panic("cannot produce chainCode")
		}
		cs = append(cs, *inner)
	}

	return cs
}

func contourTracking(img bitmap, value byte, inner bool) *chaincode {
	for y := range img {
		for x, v := range img[y] {
			if v == value {
				start := point{x, y}
				return newChaincode(img, start, value, inner)
			}
		}
	}
	return nil
}
