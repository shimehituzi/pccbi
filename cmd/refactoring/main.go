package main

import (
	"fmt"

	"github.com/shimehituzi/pccbi/internal/refactoring"
)

func main() {
	relativePath := "../../DATABASE/orig/soldier/soldier/Ply/soldier_vox10_0537.ply"
	srcPath := relativePath[4:]
	bc, err := refactoring.LoadPly(srcPath, refactoring.XYZ.Order())
	if err != nil {
		panic(err)
	}
	fmt.Println(bc.Bias)
	fmt.Println(bc.Length)
}
