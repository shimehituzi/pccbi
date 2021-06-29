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

func NewLabeledPointCloud(bc *bitCube) (*labeledPointCloud, labeledBitMaps, labeledBitMaps) {
	lpc := new(labeledPointCloud)
	lpc.length = bc.Length
	lbms := make([]labeledBitMap, lpc.length[0])
	tmp := make([]labeledBitMap, lpc.length[0])
	outers := make([][]chainCode, lpc.length[0])

	wg := &sync.WaitGroup{}
	for i := range bc.Data {
		wg.Add(1)
		go func(i int) {
			lbms[i], outers[i] = newLabeledBitMap(bc.Data[i])
			wg.Done()
		}(i)
	}
	wg.Wait()

	for i := range lbms {
		wg.Add(1)
		go func(i int) {
			tmp[i] = fillLabeledBitMap(lbms[i], outers[i])
			wg.Done()
		}(i)
	}
	wg.Wait()

	return lpc, tmp, lbms
}
