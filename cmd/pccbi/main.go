package main

import (
	"fmt"
	"time"

	"github.com/shimehituzi/pccbi/internal/bitstream"
	"github.com/shimehituzi/pccbi/internal/codec"
)

func main() {
	start := time.Now()

	// argument
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

	// Test
	TestHeader(encHeader, decHeader)
	TestVoxel(encVoxel, decVoxel)
	TestContour(encContour, decContour)
	TestStream(encStream, decStream)

	end := time.Now()
	fmt.Println(end.Sub(start).Seconds())
}

func TestHeader(encHeader, decHeader *codec.Header) {
	if encHeader.Axis != decHeader.Axis {
		panic("The Header.Axis is different")
	}
	for i := range encHeader.Length {
		if encHeader.Length[i] != decHeader.Length[i] {
			panic("The Header.Length is different")
		}
	}
	for i := range encHeader.Bias {
		if encHeader.Bias[i] != decHeader.Bias[i] {
			panic("The Header.Bias is different")
		}
	}
}

func TestVoxel(encVoxel, decVoxel codec.Voxel) {
	if len(encVoxel) != len(decVoxel) {
		panic("The Voxel length is different")
	}
	if len(encVoxel[0]) != len(decVoxel[0]) {
		panic("The Voxel[0] length is different")
	}
	if len(encVoxel[0][0]) != len(decVoxel[0][0]) {
		panic("The Voxel[0][0] length is different")
	}
	for f := range encVoxel {
		for y := range encVoxel[f] {
			for x := range encVoxel[f][y] {
				if encVoxel[f][y][x] != decVoxel[f][y][x] {
					panic("The Voxel is different")
				}
			}
		}
	}
}

func TestContour(encContour, decContour codec.Contour) {
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

func TestStream(encStream, decStream *codec.Stream) {
	if len(encStream.StartPoints) != len(decStream.StartPoints) {
		panic("The Stream.StartPoints Length is different")
	}
	if len(encStream.Codes) != len(decStream.Codes) {
		panic("The Stream.Codes Length is different")
	}
	if len(encStream.NumCodesArray) != len(decStream.NumCodesArray) {
		panic("The Stream.NumCodesArray Length is different")
	}
	for i := range encStream.StartPoints {
		for j := range encStream.StartPoints[i] {
			if encStream.StartPoints[i][j] != decStream.StartPoints[i][j] {
				panic("The Stream.StartPoints is different")
			}
		}
	}
	for i := range encStream.Codes {
		if encStream.Codes[i] != decStream.Codes[i] {
			panic("The Stream.Codes is different")
		}
	}
	for i := range encStream.NumCodesArray {
		if encStream.NumCodesArray[i] != decStream.NumCodesArray[i] {
			panic("The Stream.NumCodesArray is different")
		}
	}
}
