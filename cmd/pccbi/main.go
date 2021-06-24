package main

import (
	"github.com/shimehituzi/pccbi/internal/plyio"
)

func main() {
	relativePath := "../../DATABASE/orig/soldier/soldier/Ply/soldier_vox10_0537.ply"
	srcfile := relativePath[4:]
	sortOrders := plyio.SortOrders{1, 2, 0}
	bms, err := plyio.LordPly(srcfile, sortOrders)
	if err != nil {
		panic(err)
	}
	// lbms := labeling.NewLabeledBitMaps(bms)
	FyneLoop(bms)
}
