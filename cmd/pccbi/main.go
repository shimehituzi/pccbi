package main

import (
	"flag"
	"fmt"

	"github.com/shimehituzi/pccbi/internal/plyio"
)

func main() {
	var srcPath string
	flag.StringVar(&srcPath, "s", "", "入力ファイルのパス")
	flag.Parse()

	ply := plyio.NewPly()
	ply.Read(srcPath)

	points, err := plyio.NewPoints([3]int{1, 2, 0})
	if err != nil {
		panic(err)
	}
	points.Read(ply)

	for i := range ply.Data {
		if i < 30 {
			fmt.Println(points.Data[i])
		}
	}
}
