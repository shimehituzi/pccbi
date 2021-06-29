package main

import (
	"github.com/shimehituzi/pccbi/internal/refactoring"
)

func main() {
	relativePath := "../../DATABASE/orig/soldier/soldier/Ply/soldier_vox10_0537.ply"
	srcPath := relativePath[4:]
	bc, err := refactoring.LoadPly(srcPath, refactoring.YZX.Order())
	if err != nil {
		panic(err)
	}
	_, lbms := refactoring.NewLabeledPointCloud(bc)
	fyneLoop(lbms)
}
