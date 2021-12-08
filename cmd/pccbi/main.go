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
	order := encoder.YZX.Order()
	distPath := "compressed"

	encVoxel, encHeader := encoder.LoadPly(srcPath, order)
	encContourBuffer := encoder.NewContourBuffer(encVoxel, encHeader)
	encStream := encoder.NewStream(encContourBuffer)
	codec.Encode(encStream, encHeader, distPath)

	decStream, decHeader := codec.Decode(distPath)

	for i := range encHeader.Axis {
		if encHeader.Axis[i] != decHeader.Axis[i] {
			fmt.Println("error header axis", i)
		}
	}
	for i := range encHeader.Length {
		if encHeader.Length[i] != decHeader.Length[i] {
			fmt.Println("error header axis", i)
		}
	}
	for i := range encHeader.Bias {
		if encHeader.Bias[i] != decHeader.Bias[i] {
			fmt.Println("error header axis", i)
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
