package processing

import (
	"bufio"
	"fmt"
	"os"
)

type streamStruct struct {
	header           header
	outerStartPoints []point
	innerStartPoints []point
	outerCodes       [][]byte
	innerCodes       [][]byte
}

type Stream struct {
	Header           []int
	OuterStartPoints []int
	InnerStartPoints []int
	OuterCodes       []byte
	InnerCodes       []byte
	Codes            []byte
}

func (lpc *labeledPointCloud) MakeStreamStruct() *streamStruct {

	header := lpc.header

	outerStartPoints := make([]point, lpc.header.numOuterContours)
	outerCodes := make([][]byte, lpc.header.numOuterContours)
	innerStartPoints := make([]point, lpc.header.numInnerContours)
	innerCodes := make([][]byte, lpc.header.numInnerContours)

	innerIndex := 0
	outerIndex := 0
	for _, frame := range lpc.frames {
		for _, contour := range frame.contours {
			for _, inner := range contour.inners {
				innerStartPoints[innerIndex] = inner.start
				innerCodes[innerIndex] = inner.code
				innerIndex++
			}
			outerStartPoints[outerIndex] = contour.outer.start
			outerCodes[outerIndex] = contour.outer.code
			outerIndex++
		}
	}

	return &streamStruct{
		header,
		outerStartPoints,
		innerStartPoints,
		outerCodes,
		innerCodes,
	}
}

func (ss *streamStruct) GetStream() *Stream {
	header := []int{
		ss.header.axis[0],   //axisX
		ss.header.axis[1],   //axisY
		ss.header.axis[2],   //axisZ
		ss.header.length[0], //frames
		ss.header.length[1], //height
		ss.header.length[2], //width
		ss.header.bias[0],   //biasFrames
		ss.header.bias[1],   //biasHeight
		ss.header.bias[2],   //biasWidth
		ss.header.numOuterContours,
		ss.header.numInnerContours,
	}

	outerStartPoints := make([]int, ss.header.numOuterContours*2)
	for i, point := range ss.outerStartPoints {
		outerStartPoints[i*2] = point.x
		outerStartPoints[i*2+1] = point.y
	}
	innerStartPoints := make([]int, ss.header.numInnerContours*2)
	for i, point := range ss.innerStartPoints {
		innerStartPoints[i*2] = point.x
		innerStartPoints[i*2+1] = point.y
	}

	outerCodes := []byte{}
	for _, code := range ss.outerCodes {
		outerCodes = append(outerCodes, code...)
		outerCodes = append(outerCodes, 8)
	}

	innerCodes := []byte{}
	for _, code := range ss.innerCodes {
		innerCodes = append(innerCodes, code...)
		innerCodes = append(innerCodes, 8)
	}

	codes := []byte{}
	// dummyO := make([]byte, len(outerCodes))
	// for i := range dummyO {
	// 	elem := byte(0)
	// 	if i%10 == 1 {
	// 		elem = byte(i%8) + byte(1)
	// 	}
	// 	dummyO[i] = elem
	// }
	// codes = append(codes, dummyO...)
	// dummyI := make([]byte, len(innerCodes))
	// for i := range dummyI {
	// 	elem := byte(0)
	// 	if i%10 == 1 {
	// 		elem = byte(i%8) + byte(1)
	// 	}
	// 	dummyI[i] = elem
	// }
	// codes = append(codes, dummyI...)
	codes = append(codes, outerCodes...)
	codes = append(codes, innerCodes...)

	return &Stream{
		header,
		outerStartPoints,
		innerStartPoints,
		outerCodes,
		innerCodes,
		codes,
	}
}

func (s *Stream) Write() {
	fp, err := os.Create("compressed")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	w := bufio.NewWriter(fp)

	w.WriteString("header\n")
	for _, header := range s.Header {
		w.WriteString(fmt.Sprint(header))
		w.WriteString(" ")
	}

	w.WriteString("\nOuterStartPoints\n")
	for _, point := range s.OuterStartPoints {
		w.WriteString(fmt.Sprint(point))
		w.WriteString(" ")
	}
	w.WriteString("\nInnerStartPoints\n")

	for _, point := range s.InnerStartPoints {
		w.WriteString(fmt.Sprint(point))
		w.WriteString(" ")
	}

	w.WriteString("\nOuterCodes\n")
	for _, code := range s.OuterCodes {
		w.WriteString(fmt.Sprint(int(code)))
	}

	w.WriteString("\nInnerCodes\n")
	for _, code := range s.InnerCodes {
		w.WriteString(fmt.Sprint(int(code)))
	}

	w.Flush()
}
