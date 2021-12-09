package main

import (
	"fmt"
	"time"

	"github.com/shimehituzi/pccbi/internal/bitstream"
	"github.com/shimehituzi/pccbi/internal/codec"
)

func main() {
	start := time.Now()

	srcPath := "../../DATABASE/orig/soldier/soldier/Ply/soldier_vox10_0537.ply"[4:] // 1062090 点のデータ
	axis := codec.YZX
	distPath := "compressed"

	encPly := codec.NewPly(srcPath)
	encHeader := codec.NewHeader(encPly, axis)
	encVoxel := codec.NewVoxel(encPly, encHeader)
	encContour := codec.NewContour(encVoxel, encHeader)
	encStream := codec.NewStream(encContour)

	bitstream.Encode(encStream, encHeader, distPath)
	decStream, decHeader := bitstream.Decode(distPath)

	codec.NewDecContour(decStream, decHeader)

	end := time.Now()
	fmt.Println(end.Sub(start).Seconds())
}
