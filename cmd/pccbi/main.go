package main

import (
	"fmt"
	"time"

	"github.com/shimehituzi/pccbi/internal/codec"
	"github.com/shimehituzi/pccbi/internal/encoder"
)

func main() {
	start := time.Now()
	relativePath := "../../DATABASE/orig/soldier/soldier/Ply/soldier_vox10_0537.ply"
	// 1062090 点のデータ
	srcPath := relativePath[4:]
	voxel, err := encoder.LoadPly(srcPath, encoder.YZX.Order())
	if err != nil {
		panic(err)
	}
	cb := encoder.NewContourBuffer(voxel)

	encStream := encoder.NewStream(voxel, cb)
	codec.Encode(encStream)

	decStream := codec.Decode()

	for i := range encStream.Header {
		if encStream.Header[i] != decStream.Header[i] {
			fmt.Println("error header", i)
		}
	}
	for i := range encStream.OuterStartPoints {
		if encStream.OuterStartPoints[i] != decStream.OuterStartPoints[i] {
			fmt.Println("error outerStartPoints", i)
		}
	}
	for i := range encStream.InnerStartPoints {
		if encStream.InnerStartPoints[i] != decStream.InnerStartPoints[i] {
			fmt.Println("error innerStartPoints", i)
		}
	}
	for i := range encStream.OuterCodes {
		if encStream.OuterCodes[i] != decStream.OuterCodes[i] {
			fmt.Println("error outerCodes", i)
		}
	}
	for i := range encStream.InnerCodes {
		if encStream.InnerCodes[i] != decStream.InnerCodes[i] {
			fmt.Println("error innerCodes", i)
		}
	}

	end := time.Now()
	fmt.Println(end.Sub(start).Seconds())
}
