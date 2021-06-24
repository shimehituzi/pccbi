package main

import (
	"fmt"

	"github.com/shimehituzi/pccbi/internal/labeling"
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
	lbms := labeling.NewLabeledBitMaps(bms)
	contours := (*lbms)[0].Contour
	segments := labeling.NewSegments(contours)
	fmt.Println(segments[1].Label)
	FyneLoop(lbms)
}
