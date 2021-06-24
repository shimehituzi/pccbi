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

type Segment struct {
	Contours []Contour
	Label    int
}

type LabeledBitMap struct {
	Image   [][]byte
	Segment []Segment
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

	contours := []Contour{}
	for i := 0; ; i++ {
		cc := CountourTracking(tmp)
		for _, point := range cc.Points {
			tmp[point.Y][point.X] = 0
		}
		contours = append(
			contours,
			Contour{ChainCode: *cc, Label: i},
		)
		if isNotExistPoint(tmp) {
			break
		}
	}

	lbm.Segment = NewSegments(contours)

	return lbm
}

func NewSegments(contours []Contour) []Segment {
	segments := []Segment{}

	count := 0
	for i := range contours {
		if label := isAdjacentSegment(contours[i], segments); label == -1 {
			count++
			segments = append(segments, Segment{
				Contours: []Contour{contours[i]},
				Label:    count,
			})
		} else {
			segments[label-1].Contours = append(segments[label-1].Contours, contours[i])
		}
	}

	return segments
}

func isAdjacentSegment(contour Contour, segments []Segment) int {
	if len(segments) == 0 {
		return -1
	}
	for _, point := range contour.ChainCode.Points {
		for _, d := range getDirection() {
			adjacnet := Point{
				point.X + d.Dx,
				point.Y + d.Dy,
			}
			for _, segment := range segments {
				for _, arroundContour := range segment.Contours {
					for _, arroundPoint := range arroundContour.ChainCode.Points {
						if adjacnet == arroundPoint {
							return segment.Label
						}
					}
				}
			}
		}
	}
	return -1
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
