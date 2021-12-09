package main

import (
	"fmt"
	"time"

	"github.com/shimehituzi/pccbi/internal/bitstream"
	"github.com/shimehituzi/pccbi/internal/codec"
)

func main() {
	start := time.Now()

	axis := codec.YZX
	srcPath := "../../DATABASE/orig/soldier/soldier/Ply/soldier_vox10_0537.ply"[4:] // 1062090 点のデータ
	distPath := "compressed"

	encVoxel, encHeader := codec.ReadPly(srcPath, axis)
	encContour := codec.EncContour(encVoxel, encHeader)
	encStream := codec.EncStream(encContour)

	bitstream.Encode(encStream, encHeader, distPath)
	decStream, decHeader := bitstream.Decode(distPath)

	codec.DecStream(decStream, decHeader)

	end := time.Now()
	fmt.Println(end.Sub(start).Seconds())
}
