package decoder

import "bytes"

func NewContourBuffer(stream *Stream, header *Header) contourBuffer {
	codes := bytes.Split(uint2byte(stream.Codes), []byte{8})
	cb := make(contourBuffer, header.Length[0])

	i := 0
	for _, numCodes := range stream.NumCodesArray {
		contour := make(contour, numCodes)
		f := int(stream.StartPoints[i][0])
		for j := 0; j < int(numCodes); j++ {
			startY := int(stream.StartPoints[i][1])
			startX := int(stream.StartPoints[i][2])
			cc := chainCode{
				start:  point{startX, startY},
				code:   codes[i],
				points: []point{},
			}
			contour[j] = cc
			i++
		}
		cb[f] = append(cb[f], contour)
	}

	return cb
}

func uint2byte(uintSlice []uint) (byteSlice []byte) {
	byteSlice = make([]byte, len(uintSlice))
	for i, v := range uintSlice {
		byteSlice[i] = byte(v % 256)
	}
	return
}
