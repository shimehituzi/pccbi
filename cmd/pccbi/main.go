package main

import (
	"fmt"
	"time"

	"github.com/shimehituzi/pccbi/internal/codec"
	"github.com/shimehituzi/pccbi/internal/decoder"
	"github.com/shimehituzi/pccbi/internal/encoder"
)

func main() {
	start := time.Now()

	srcPath := "../../DATABASE/orig/soldier/soldier/Ply/soldier_vox10_0537.ply"[4:] // 1062090 点のデータ
	axis := encoder.YZX
	distPath := "compressed"

	encPly := encoder.NewPly(srcPath)
	encHeader := encoder.NewHeader(encPly, axis)
	encVoxel := encoder.NewVoxel(encPly, encHeader)
	encContourBuffer := encoder.NewContourBuffer(encVoxel, encHeader)
	encStream := encoder.NewStream(encContourBuffer)

	codec.Encode(encStream, encHeader, distPath)
	decStream, decHeader := codec.Decode(distPath)

	decoder.NewContourBuffer(decStream, decHeader)

	end := time.Now()
	fmt.Println(end.Sub(start).Seconds())
}
