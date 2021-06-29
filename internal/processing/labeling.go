package processing

import (
	"sync"
)

type labeledPointCloud struct {
	frames []frame
	length [3]int
}

type frame struct {
	img      [][]byte
	contours []contour
}

type contour struct {
	outer chainCode
	inner []chainCode
	label int
}

type chainCode struct {
	start  point
	code   []byte
	points []point
}

type point struct {
	x, y int
}

type rect struct {
	max point
	min point
}

func NewLabeledPointCloud(bc *bitCube) (*labeledPointCloud, labeledBitMaps) {
	numOfFrame := bc.length[0]

	lbms := make([]labeledBitMap, numOfFrame)
	outerMatrix := make([][]chainCode, numOfFrame)

	rectMatrix := make([][]rect, numOfFrame)
	segmentMatrix := make([][]labeledBitMap, numOfFrame)
	numOfInnerMatrix := make([][]int, numOfFrame)

	inner3dMarix := make([][][]chainCode, numOfFrame)

	wg := &sync.WaitGroup{}
	for i := range bc.data {
		wg.Add(1)
		go func(i int) {
			lbms[i], outerMatrix[i] = newLabeledBitMap(bc.data[i])
			wg.Done()
		}(i)
	}
	wg.Wait()

	for i := range lbms {
		wg.Add(1)
		go func(i int) {
			segmentMatrix[i], rectMatrix[i], numOfInnerMatrix[i] = getFilledSegments(lbms[i], outerMatrix[i])
			wg.Done()
		}(i)
	}
	wg.Wait()

	for i := range segmentMatrix {
		wg.Add(1)
		go func(i int) {
			inner3dMarix[i] = getInnerContour(segmentMatrix[i], rectMatrix[i], numOfInnerMatrix[i])
			wg.Done()
		}(i)
	}
	wg.Wait()

	lpc := new(labeledPointCloud)
	lpc.length = bc.length
	lpc.frames = make([]frame, lpc.length[0])
	for f := range lpc.frames {
		lpc.frames[f].img = make([][]byte, lpc.length[1])
		for i := range lpc.frames[f].img {
			lpc.frames[f].img[i] = make([]byte, lpc.length[2])
		}
		lpc.frames[f].contours = make([]contour, len(outerMatrix[f]))
		for label := range lpc.frames[f].contours {
			lpc.frames[f].contours[label].outer = outerMatrix[f][label]
			lpc.frames[f].contours[label].inner = inner3dMarix[f][label]
			lpc.frames[f].contours[label].label = label + 1

			for _, point := range lpc.frames[f].contours[label].outer.points {
				lpc.frames[f].img[point.y][point.x] = 1
			}
			for _, inner := range lpc.frames[f].contours[label].inner {
				for _, point := range inner.points {
					lpc.frames[f].img[point.y][point.x] = 2
				}
			}
		}
	}

	return lpc, lbms
}
