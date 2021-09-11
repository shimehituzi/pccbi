package main

import (
	"github.com/shimehituzi/pccbi/internal/refactoring"
)

func main() {
	relativePath := "../../DATABASE/orig/soldier/soldier/Ply/soldier_vox10_0537.ply"
	// 1062090 点のデータ
	srcPath := relativePath[4:]
	voxel, err := refactoring.LoadPly(srcPath, refactoring.YZX.Order())
	if err != nil {
		panic(err)
	}
	frames := refactoring.NewFrames(voxel)
	cb := refactoring.NewContourBuffer(voxel)
	fc := refactoring.NewFyneContour(cb, voxel)
	fyneLoop([]refactoring.FyneBitMap{fc, frames})
}
