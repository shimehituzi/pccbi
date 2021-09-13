package main

import (
	"github.com/shimehituzi/pccbi/internal/processing"
)

func main() {
	relativePath := "../../DATABASE/orig/soldier/soldier/Ply/soldier_vox10_0537.ply"
	// 1062090 点のデータ
	srcPath := relativePath[4:]
	voxel, err := processing.LoadPly(srcPath, processing.YZX.Order())
	if err != nil {
		panic(err)
	}
	frames := processing.NewFrames(voxel)
	cb := processing.NewContourBuffer(voxel)
	fc := processing.NewFyneContour(cb, voxel)
	fyneLoop([]processing.FyneBitMap{fc, frames})
}
