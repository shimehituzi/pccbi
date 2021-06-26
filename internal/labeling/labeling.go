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

	lbm.FillInnerArea()

	return lbm
}

type rectMinMax struct {
	minX, maxX, minY, maxY int
}

func (lbm *LabeledBitMap) FillInnerArea() {

	imgs := make([][][]byte, len(lbm.Segment))
	wg := &sync.WaitGroup{}
	for s, segment := range lbm.Segment {
		wg.Add(1)
		go func(s int, segment Segment) {
			imgs[s] = make([][]byte, len(lbm.Image))
			for i := range imgs[s] {
				imgs[s][i] = make([]byte, len(lbm.Image[i]))
			}
			rect := rectMinMax{
				minX: len(imgs[s][0]),
				maxX: 0,
				minY: len(imgs[s]),
				maxY: 0,
			}
			for _, contour := range segment.Contours {
				for _, point := range contour.ChainCode.Points {
					imgs[s][point.Y][point.X] = 1

					if rect.minX > point.X {
						rect.minX = point.X
					}
					if rect.maxX < point.X {
						rect.maxX = point.X
					}
					if rect.minY > point.Y {
						rect.minY = point.Y
					}
					if rect.maxY < point.Y {
						rect.maxY = point.Y
					}
				}
			}
			for y := rect.minY; y <= rect.maxY; y++ {
				for x := rect.minX; x <= rect.maxX; x++ {
					FillArea(imgs[s], Point{x, y}, segment.Contours[0], rect)
				}
			}
			wg.Done()
		}(s, segment)
	}
	wg.Wait()

	for _, img := range imgs {
		for y := range img {
			for x := range img[y] {
				if lbm.Image[y][x] == 0 && img[y][x] == 2 {
					lbm.Image[y][x] = 2
				}
			}
		}
	}
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
