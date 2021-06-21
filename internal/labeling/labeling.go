package labeling

import (
	"fmt"

	"github.com/shimehituzi/pccbi/internal/plyio"
)

type ChainCode struct {
	Start struct{ X, Y int }
	code  []byte
}

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

	lbm.CountourTracking()

	return lbm
}

func (lbm *LabeledBitMap) CountourTracking() {
	cc := new(ChainCode)
	cc.code = []byte{0}
	direction := direction()
	for y := range lbm.Image {
		for x := range lbm.Image[y] {
			if lbm.Image[y][x] == 1 {
				cc.Start = struct{ X, Y int }{x, y}
				current := cc.Start
				for i := 0; ; i++ {
					prevD := cc.code[i]
					for j := 0; j < 8; j++ {
						currentD := (5 + int(prevD) + j) % 8
						d := direction[currentD]
						dy := current.Y + d.Dy
						dx := current.X + d.Dx
						if dy < 0 || dx < 0 || len(lbm.Image) <= dy || len(lbm.Image[0]) <= dx {
							continue
						}
						if lbm.Image[dy][dx] == 1 {
							cc.code = append(cc.code, byte(currentD))
							current = struct {
								X int
								Y int
							}{dx, dy}
							break
						}
					}
					if (current.X == cc.Start.X) && (current.Y == cc.Start.Y) {
						break
					}
				}
				break
			}
		}
		break
	}
	fmt.Println("cc.start", cc.Start)
	fmt.Println("cc.code", cc.code[1:])
}

type Direction struct {
	Dx, Dy int
}

func direction() [8]Direction {
	return [8]Direction{
		{1, 0}, {1, 1}, {0, 1}, {-1, 1},
		{-1, 0}, {-1, -1}, {0, -1}, {1, -1},
	}
}
