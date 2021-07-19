package main

import (
	"github.com/shimehituzi/pccbi/internal/codec"
	"github.com/shimehituzi/pccbi/internal/processing"
)

func main() {
	relativePath := "../../DATABASE/orig/soldier/soldier/Ply/soldier_vox10_0537.ply"
	// 1062090 点のデータ
	srcPath := relativePath[4:]
	bc, err := processing.LoadPly(srcPath, processing.YZX.Order())
	if err != nil {
		panic(err)
	}
	lpc, lbms := processing.NewLabeledPointCloud(bc)
	codec.Encode(*lpc.MakeStreamStruct().GetStream())
	fyneLoop([]processing.FyneBitMap{lpc, lbms, bc})
}
