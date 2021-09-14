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
	for i := range encStream.StartPoints {
		if encStream.StartPoints[i] != decStream.StartPoints[i] {
			fmt.Println("error startPoints", i)
		}
	}
	for i := range encStream.Flags {
		if encStream.Flags[i] != decStream.Flags[i] {
			fmt.Println("error flags", i)
		}
	}
	for i := range encStream.Codes {
		if encStream.Codes[i] != decStream.Codes[i] {
			fmt.Println("error codes", i)
		}
	}

	end := time.Now()
	fmt.Println(end.Sub(start).Seconds())
}
