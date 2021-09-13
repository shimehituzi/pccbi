package codec

import (
	"bufio"
	"os"

	"github.com/shimehituzi/pccbi/internal/processing"
)

func Decode() *processing.Stream {
	fp, err := os.Open("compressed")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	r := bufio.NewReader(fp)

	bitbuf := NewBitbuf(false)

	bitSize := 16
	bigBitSize := 32

	headerLength := 13
	header := make([]int, headerLength)
	for i := range header {
		if i > 10 {
			header[i] = int(bitbuf.getbits(r, bigBitSize))
		} else {
			header[i] = int(bitbuf.getbits(r, bitSize))
		}
	}
	numOuter := header[9]
	numInner := header[10]
	numOuterCodes := header[11]
	numInnerCodes := header[11]

	freqLength := 9
	cumfreqLength := 10

	outerFreq := make([]uint64, freqLength)
	for i := range outerFreq {
		outerFreq[i] = uint64(bitbuf.getbits(r, bigBitSize))
	}
	outerCumfreq := make([]uint64, cumfreqLength)
	for i := range outerCumfreq {
		outerCumfreq[i] = uint64(bitbuf.getbits(r, bigBitSize))
	}
	outerTotfreq := uint64(bitbuf.getbits(r, bigBitSize))
	outerOffset := uint64(bitbuf.getbits(r, bitSize))

	outerPmodel := &Pmodel{
		freq:    outerFreq,
		cumfreq: outerCumfreq,
		totfreq: outerTotfreq,
		offset:  outerOffset,
	}

	innerFreq := make([]uint64, freqLength)
	for i := range innerFreq {
		innerFreq[i] = uint64(bitbuf.getbits(r, bigBitSize))
	}
	innerCumfreq := make([]uint64, cumfreqLength)
	for i := range innerCumfreq {
		innerCumfreq[i] = uint64(bitbuf.getbits(r, bigBitSize))
	}
	innerTotfreq := uint64(bitbuf.getbits(r, bigBitSize))
	innerOffset := uint64(bitbuf.getbits(r, bitSize))

	innerPmodel := &Pmodel{
		freq:    innerFreq,
		cumfreq: innerCumfreq,
		totfreq: innerTotfreq,
		offset:  innerOffset,
	}

	outerStartPoints := make([]int, numOuter)
	for i := range outerStartPoints {
		outerStartPoints[i] = int(bitbuf.getbits(r, bitSize))
	}
	innerStartPoints := make([]int, numInner)
	for i := range innerStartPoints {
		innerStartPoints[i] = int(bitbuf.getbits(r, bitSize))
	}

	rc := newRangeCoder()
	rc.startdec(r)

	outerCodes := make([]byte, numOuterCodes)
	for i := 0; i < numOuterCodes; i++ {
		outerCodes[i] = byte(rc.decode(r, outerPmodel))
	}
	innerCodes := make([]byte, numInnerCodes)
	for i := 0; i < numInnerCodes; i++ {
		innerCodes[i] = byte(rc.decode(r, innerPmodel))
	}

	stream := &processing.Stream{
		Header:           header,
		OuterStartPoints: outerStartPoints,
		InnerStartPoints: innerStartPoints,
		OuterCodes:       outerCodes,
		InnerCodes:       innerCodes,
	}

	return stream
}
