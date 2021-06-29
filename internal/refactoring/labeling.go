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
	outer  chainCode
	innter []chainCode
	label  int
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
	lpc := new(labeledPointCloud)
	lpc.length = bc.Length
	lbms := make([]labeledBitMap, lpc.length[0])
	outerMatrix := make([][]chainCode, lpc.length[0])
	segmentMatrix := make([][]labeledBitMap, lpc.length[0])

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
			segmentMatrix[i] = getFilledSegments(lbms[i], outerMatrix[i])
			wg.Done()
		}(i)
	}
	wg.Wait()

	return lpc, lbms
}
