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
	recPath := "rec.ply"

	// Encode
	encVoxel, encHeader := codec.ReadPly(srcPath, axis)
	encContour := codec.EncContour(encVoxel, encHeader)
	encStream := codec.EncStream(encContour)
	bitstream.Encode(encStream, encHeader, distPath)

	// Decode
	decStream, decHeader := bitstream.Decode(distPath)
	decContour := codec.DecStream(decStream, decHeader)
	decVoxel := codec.DecContour(decContour, decHeader)
	codec.WritePly(recPath, decVoxel, decHeader)

	Test(encContour, decContour)

	end := time.Now()
	fmt.Println(end.Sub(start).Seconds())
}

func Test(encContour, decContour codec.Contour) {
	if len(encContour) != len(decContour) {
		panic("The contour length is different")
	}
	for f := range encContour {
		if len(encContour[f]) != len(decContour[f]) {
			panic("The contour[f] length is different")
		}
		for l := range encContour[f] {
			if len(encContour[f][l]) != len(decContour[f][l]) {
				panic("The contour[f][l] length is different")
			}
			for i := range encContour[f][l] {
				encChaincode := encContour[f][l][i]
				decChainCode := decContour[f][l][i]
				if !codec.ComparePoint(encChaincode.Start, decChainCode.Start) {
					panic("The start point is different")
				}
				if len(encChaincode.Code) != len(decChainCode.Code) {
					panic("The chaincode.code length is different")
				}
				if len(encChaincode.Points) != len(decChainCode.Points) {
					panic("The chaincode.points length is different")
				}
				for j := range encChaincode.Code {
					if encChaincode.Code[j] != decChainCode.Code[j] {
						panic("The Chaincode is different")
					}
				}
				for j := range encChaincode.Points {
					if !codec.ComparePoint(encChaincode.Points[j], decChainCode.Points[j]) {
						panic("The Points is different")
					}
				}
			}
		}
	}
}
