package codec

import (
	"bytes"
)

func DecStream(stream *Stream, header *Header) Contour {
	// stream.Codes を 8 で split する
	codes := bytes.Split(uint2byte(stream.Codes), []byte{8})

	// contour[frame][segment][外輪郭or内輪郭] の形に整形
	contour := make(Contour, header.Length[0])
	i := 0
	for _, numCodes := range stream.NumCodesArray {
		cs := make([]chaincode, numCodes)
		f := int(stream.StartPoints[i][0])
		for j := 0; j < int(numCodes); j++ {
			startY := int(stream.StartPoints[i][1])
			startX := int(stream.StartPoints[i][2])
			cc := chaincode{
				Start: point{startX, startY},
				Code:  decDecorrelation(codes[i]),
			}
			cs[j] = cc
			i++
		}
		contour[f] = append(contour[f], cs)
	}

	return contour
}
