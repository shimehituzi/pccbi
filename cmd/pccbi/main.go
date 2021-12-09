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

	decContour := codec.DecStream(decStream, decHeader)

	for f := range encContour {
		for l := range encContour[f] {
			for i := range encContour[f][l] {
				encChaincode := encContour[f][l][i]
				decChainCode := decContour[f][l][i]
				if !codec.ComparePoint(encChaincode.Start, decChainCode.Start) {
					panic("The start point is different")
				}
				for j := range encChaincode.Code {
					if encChaincode.Code[j] != decChainCode.Code[j] {
						panic("The Chaincode is different")
					}
				}
			}
		}
	}

	end := time.Now()
	fmt.Println(end.Sub(start).Seconds())
}
