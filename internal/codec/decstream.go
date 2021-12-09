package codec

import "bytes"

func DecStream(stream *Stream, header *Header) contour {
	codes := bytes.Split(uint2byte(stream.Codes), []byte{8})
	contour := make(contour, header.Length[0])

	i := 0
	for _, numCodes := range stream.NumCodesArray {
		cs := make([]chaincode, numCodes)
		f := int(stream.StartPoints[i][0])
		for j := 0; j < int(numCodes); j++ {
			startY := int(stream.StartPoints[i][1])
			startX := int(stream.StartPoints[i][2])
			cc := chaincode{
				start:  point{startX, startY},
				code:   codes[i],
				points: []point{},
			}
			cs[j] = cc
			i++
		}
		contour[f] = append(contour[f], cs)
	}

	return contour
}

func uint2byte(uintSlice []uint) (byteSlice []byte) {
	byteSlice = make([]byte, len(uintSlice))
	for i, v := range uintSlice {
		byteSlice[i] = byte(v % 256)
	}
	return
}
