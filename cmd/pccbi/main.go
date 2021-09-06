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
	lv, _ := refactoring.NewLabels(voxel)
	frames := refactoring.NewFrames(voxel)
	fyneLoop([]refactoring.FyneBitMap{lv, frames})
}
