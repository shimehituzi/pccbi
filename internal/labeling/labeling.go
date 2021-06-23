package labeling

import (
	"sync"

	"github.com/shimehituzi/pccbi/internal/bitmap"
	"github.com/shimehituzi/pccbi/internal/plyio"
)

type Contour struct {
	ChainCode ChainCode
	Label     int
}

// type Segment struct {
// 	Countour []Contour
// 	Label    int
// }

type LabeledBitMap struct {
	Image   [][]byte
	Contour []Contour
}

type LabeledBitMaps []LabeledBitMap

func NewLabeledBitMaps(bms *plyio.BitMaps) *LabeledBitMaps {
	lbms := new(LabeledBitMaps)
	*lbms = make([]LabeledBitMap, len(bms.Data))

	wg := &sync.WaitGroup{}
	for i, bm := range bms.Data {
		wg.Add(1)
		go func(i int, bm bitmap.BitMap) {
			(*lbms)[i] = *NewLabeledBitMap(bm)
			wg.Done()
		}(i, bm)
	}
	wg.Wait()

	return lbms
}

func NewLabeledBitMap(bm [][]byte) *LabeledBitMap {
	lbm := new(LabeledBitMap)
	lbm.Image = make([][]byte, len(bm))
	tmp := make([][]byte, len(bm))
	for i := range lbm.Image {
		lbm.Image[i] = make([]byte, len(bm[i]))
		tmp[i] = make([]byte, len(bm[i]))
		copy(lbm.Image[i], bm[i])
		copy(tmp[i], bm[i])
	}

	for i := 0; ; i++ {
		cc := CountourTracking(tmp)
		for _, point := range cc.Points {
			tmp[point.Y][point.X] = 0
		}
		lbm.Contour = append(
			lbm.Contour,
			Contour{ChainCode: *cc, Label: i},
		)
		if isNotExistPoint(tmp) {
			break
		}
	}

	return lbm
}

func isNotExistPoint(bm [][]byte) bool {
	for y := range bm {
		for x := range bm[y] {
			if bm[y][x] == 1 {
				return false
			}
		}
	}
	return true
}
