package main

import (
	"flag"

	"github.com/shimehituzi/pccbi/internal/plyio"
)

func main() {
	var srcfile string
	flag.StringVar(&srcfile, "s", "", "入力ファイルのパス")
	flag.Parse()

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

	frames := plyio.NewFrames()
	frames.ReadPoints(points)

	bm := plyio.NewBitMaps()
	bm.ReadPoints(points)
}
