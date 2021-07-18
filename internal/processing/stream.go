package processing

import (
	"bufio"
	"fmt"
	"os"
)

type streamStruct struct {
	Header           header
	OuterStartPoints []point
	InnerStartPoints []point
	OuterCodes       [][]byte
	InnerCodes       [][]byte
}

type stream struct {
	Header           []int
	OuterStartPoints []int
	InnerStartPoints []int
	OuterCodes       []byte
	InnerCodes       []byte
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

func (ss *streamStruct) GetStream() *stream {
	header := []int{
		ss.Header.axis[0],   //axisX
		ss.Header.axis[1],   //axisY
		ss.Header.axis[2],   //axisZ
		ss.Header.length[0], //frames
		ss.Header.length[1], //height
		ss.Header.length[2], //width
		ss.Header.bias[0],   //biasFrames
		ss.Header.bias[1],   //biasHeight
		ss.Header.bias[2],   //biasWidth
		ss.Header.numOuterContours,
		ss.Header.numInnerContours,
	}

	outerStartPoints := make([]int, ss.Header.numOuterContours*2)
	for i, point := range ss.OuterStartPoints {
		outerStartPoints[i*2] = point.x
		outerStartPoints[i*2+1] = point.y
	}
	innerStartPoints := make([]int, ss.Header.numInnerContours*2)
	for i, point := range ss.InnerStartPoints {
		innerStartPoints[i*2] = point.x
		innerStartPoints[i*2+1] = point.y
	}

	outerCodes := []byte{}
	for _, code := range ss.OuterCodes {
		outerCodes = append(outerCodes, code...)
		outerCodes = append(outerCodes, 8)
	}

	innerCodes := []byte{}
	for _, code := range ss.InnerCodes {
		innerCodes = append(innerCodes, code...)
		innerCodes = append(innerCodes, 4)
	}

	return &stream{
		header,
		outerStartPoints,
		innerStartPoints,
		outerCodes,
		innerCodes,
	}
}

func (s *stream) Write() {
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
}
