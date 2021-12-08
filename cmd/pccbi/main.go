package main

import (
	"fmt"
	"time"

	"github.com/shimehituzi/pccbi/internal/codec"
	"github.com/shimehituzi/pccbi/internal/encoder"
)

func main() {
	start := time.Now()
	srcPath := "../../DATABASE/orig/soldier/soldier/Ply/soldier_vox10_0537.ply"[4:] // 1062090 点のデータ
	distPath := "compressed"

	encVoxel, encVoxelHeader := encoder.LoadPly(srcPath, encoder.YZX.Order())
	encContourBuffer := encoder.NewContourBuffer(encVoxel, encVoxelHeader)
	encStream, encStreamHeader := encoder.NewStream(encContourBuffer, encVoxelHeader)
	codec.Encode(encStream, encStreamHeader, distPath)

	decStream, decStreamHeader := codec.Decode(distPath)

	for i := range encStreamHeader {
		if encStreamHeader[i] != decStreamHeader[i] {
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
