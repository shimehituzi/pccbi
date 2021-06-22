package labeling

type Contour struct {
	ChainCode ChainCode
	Label     int
}

// type Segment struct {
// 	Countour []Contour
// 	Label    int
// }

type LabeledBitMap struct {
	Image    [][]byte
	Countour []Contour
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
		lbm.Countour = append(
			lbm.Countour,
			Contour{ChainCode: *cc, Label: i},
		)
		if isExistPoint(tmp) {
			break
		}
	}

	return lbm
}

func isExistPoint(bm [][]byte) bool {
	for y := range bm {
		for x := range bm[y] {
			if bm[y][x] == 1 {
				return true
			}
		}
	}
	return false
}
