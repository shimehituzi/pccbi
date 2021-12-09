package bitstream

import (
	"bufio"
	"os"

	"github.com/shimehituzi/pccbi/internal/codec"
)

func Decode(distPath string) (*codec.Stream, *codec.Header) {
	fp, err := os.Open(distPath)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	r := bufio.NewReader(fp)

	bitbuf := NewBitbuf(false)

	bitSize := 16
	bigBitSize := 32

	header := new(codec.Header)
	header.Axis = codec.Axis(bitbuf.getbits(r, bitSize))
	for i := range header.Length {
		header.Length[i] = int(bitbuf.getbits(r, bitSize))
	}
	for i := range header.Bias {
		header.Bias[i] = int(bitbuf.getbits(r, bitSize))
	}
	codesLength := bitbuf.getbits(r, bigBitSize)
	startPointsLength := bitbuf.getbits(r, bitSize)
	numCodesArrayLength := bitbuf.getbits(r, bitSize)
	numCodesArrayFreqMax := bitbuf.getbits(r, bitsize)

	startPoints := make([][3]uint, startPointsLength)
	for i := range startPoints {
		for j := 0; j < 3; j++ {
			startPoints[i][j] = bitbuf.getbits(r, bitSize)
		}
	}

	codeFreqMax := 8

	numCodesArrayFreq := make([]uint64, numCodesArrayFreqMax+1)
	for i := range numCodesArrayFreq {
		numCodesArrayFreq[i] = uint64(bitbuf.getbits(r, bigBitSize))
	}
	numCodesArrayPmodel := newDecPmodel(numCodesArrayFreq, 0, uint(numCodesArrayFreqMax))

	codesFreq := make([]uint64, codeFreqMax+1)
	for i := range codesFreq {
		codesFreq[i] = uint64(bitbuf.getbits(r, bigBitSize))
	}
	codesPmodel := newDecPmodel(codesFreq, 0, uint(codeFreqMax))

	rc := newRangeCoder()
	rc.startdec(r)

	numCodesArray := make([]uint, numCodesArrayLength)
	for i := range numCodesArray {
		numCodesArray[i] = uint(rc.decode(r, numCodesArrayPmodel))
	}
	codes := make([]uint, codesLength)
	for i := range codes {
		codes[i] = uint(rc.decode(r, codesPmodel))
	}

	stream := &codec.Stream{
		StartPoints:   startPoints,
		NumCodesArray: numCodesArray,
		Codes:         codes,
	}

	return stream, header
}

func newDecPmodel(freq []uint64, min, max uint) *Pmodel {
	cumfreq := make([]uint64, max+2)
	cumfreq[0] = 0
	for i := range freq {
		cumfreq[i+1] = cumfreq[i] + freq[i]
	}
	offset := cumfreq[min]
	totfreq := cumfreq[max+1] - offset

	return &Pmodel{freq, cumfreq, totfreq, offset}
}
