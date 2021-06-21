package labeling

import (
	"github.com/shimehituzi/pccbi/internal/plyio"
)

type Countour struct {
	ChainCode ChainCode
	Label     int
	Level     int
}

type Segment struct {
	Countour []Countour
	Label    int
}

type LabeledBitMap struct {
	Image    [][]byte
	Segments []Segment
}

func NewLabeledBitMap(bm *plyio.BitMap) *LabeledBitMap {
	lbm := new(LabeledBitMap)
	lbm.Image = make([][]byte, bm.Bounds().Dy())
	for y := range lbm.Image {
		lbm.Image[y] = make([]byte, bm.Bounds().Dx())
		for x := range lbm.Image[y] {
			lbm.Image[y][x] = (*bm)[y][x]
		}
	}

	CountourTracking(lbm)

	return lbm
}
