package main

import (
	"github.com/shimehituzi/pccbi/internal/labeling"
	"github.com/shimehituzi/pccbi/internal/plyio"
)

func main() {
	var bm *plyio.BitMap
	bm = &plyio.BitMap{
		[]byte{0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0},
		[]byte{0, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0},
		[]byte{0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0},
		[]byte{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		[]byte{1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1},
		[]byte{1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0},
		[]byte{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0},
		[]byte{0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0},
		[]byte{0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 0},
		[]byte{0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0},
		[]byte{0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0},
	}
	labeling.NewLabeledBitMap(bm)
}
