package main

import "github.com/shimehituzi/pccbi/internal/plyio"

func LordPly(srcfile string) *plyio.BitMaps {
	ply := plyio.NewPly()
	if err := ply.ReadPlyFile(srcfile); err != nil {
		panic(err)
	}

	points, err := plyio.NewPoints([3]int{1, 2, 0})
	if err != nil {
		panic(err)
	}
	if err := points.ReadPly(ply); err != nil {
		panic(err)
	}

	bms := plyio.NewBitMaps()
	bms.ReadPoints(points)

	return bms
}
