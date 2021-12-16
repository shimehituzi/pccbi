package codec

import "sync"

// =================
//      encode
// =================

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

// =================
//      decode
// =================

func DecContour(contour Contour, header *Header) Voxel {
	voxel := make(Voxel, header.Length[0])
	for f := range contour {
		voxel[f] = getFrame(contour[f], header.Length[1], header.Length[2])
	}

	return voxel
}

func getFrame(contourFrame [][]chaincode, height, width int) bitmap {
	frame := make(bitmap, height)
	for y := range frame {
		frame[y] = make([]byte, width)
	}
	segments := make([]bitmap, len(contourFrame))
	for l := range segments {
		segments[l] = make(bitmap, height)
		for y := range segments[l] {
			segments[l][y] = make([]byte, width)
		}
	}

	for l := range segments {
		contourSegment := contourFrame[l]
		if len(contourSegment) > 0 {
			outer := contourSegment[0]
			inners := contourSegment[1:]
			fillSegment(segments[l], outer, inners)
		} else {
			panic("outer contour is not exist")
		}
	}

	for l := range segments {
		for y := range segments[l] {
			for x, v := range segments[l][y] {
				if v == 1 {
					frame[y][x] = 1
				}
			}
		}
	}

	return frame
}

func fillSegment(img bitmap, outer chaincode, inners []chaincode) {
	outerPoints := outer.getPoints()
	for _, p := range outerPoints {
		img[p.y][p.x] = 1
	}
	for _, p := range outerPoints {
		nearest4points := [4]point{
			{p.x, p.y - 1},
			{p.x - 1, p.y},
			{p.x + 1, p.y},
			{p.x, p.y + 1},
		}
		for _, p := range nearest4points {
			if p.checkValue(img, 0) && p.isInside(outerPoints) {
				fillArea(img, p, 0, 1)
			}
		}
	}
	innerPointsArray := make([][]point, len(inners))
	for i, inner := range inners {
		innerPointsArray[i] = inner.getPoints()
	}
	for _, innerPoints := range innerPointsArray {
		for _, p := range innerPoints {
			img[p.y][p.x] = 2
		}
	}
	for _, innerPoints := range innerPointsArray {
		for _, p := range innerPoints {
			nearest4points := [4]point{
				{p.x, p.y - 1},
				{p.x - 1, p.y},
				{p.x + 1, p.y},
				{p.x, p.y + 1},
			}
			for _, p := range nearest4points {
				if p.checkValue(img, 1) && p.isInside(innerPoints) {
					fillArea(img, p, 1, 2)
				}
			}
		}
	}
}
