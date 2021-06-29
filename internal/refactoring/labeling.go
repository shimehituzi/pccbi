package refactoring

import "sync"

type labeledPointCloud struct {
	frame  []frame
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
	numOfFrame := bc.Length[0]

	lbms := make([]labeledBitMap, numOfFrame)
	outerMatrix := make([][]chainCode, numOfFrame)

	rectMatrix := make([][]rect, numOfFrame)
	segmentMatrix := make([][]labeledBitMap, numOfFrame)
	numOfInnerMatrix := make([][]int, numOfFrame)

	inner3dMarix := make([][][]chainCode, numOfFrame)

	wg := &sync.WaitGroup{}
	for i := range bc.Data {
		wg.Add(1)
		go func(i int) {
			lbms[i], outerMatrix[i] = newLabeledBitMap(bc.Data[i])
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
	lpc.length = bc.Length
	lpc.frame = make([]frame, lpc.length[0])
	for f, frame := range lpc.frame {
		frame.img = make([][]byte, lpc.length[1])
		for i := range frame.img {
			frame.img[i] = make([]byte, lpc.length[2])
			copy(frame.img[i], bc.Data[f][i])
		}
		frame.contours = make([]contour, len(outerMatrix[f]))
		for label, contour := range frame.contours {
			contour.outer = outerMatrix[f][label]
			contour.inner = inner3dMarix[f][label]
			contour.label = label + 1
		}
	}

	return lpc, lbms
}
