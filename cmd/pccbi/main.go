package main

import (
	"github.com/shimehituzi/pccbi/internal/labeling"
	"github.com/shimehituzi/pccbi/internal/plyio"
)

func main() {
	relativePath := "../../DATABASE/orig/soldier/soldier/Ply/soldier_vox10_0537.ply"
	srcfile := relativePath[4:]
	bms, err := plyio.LordPly(srcfile)
	if err != nil {
		panic(err)
	}
	lbms := labeling.NewLabeledBitMaps(bms)
	FyneLoop(lbms)
}
