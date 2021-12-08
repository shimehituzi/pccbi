package codec

import (
	"bufio"
	"os"

	"github.com/shimehituzi/pccbi/internal/decoder"
)

func Decode(distPath string) (*decoder.Stream, decoder.StreamHeader) {
	fp, err := os.Open(distPath)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	r := bufio.NewReader(fp)

	bitbuf := NewBitbuf(false)

	bitSize := 16
	bigBitSize := 32

	headerLength := 12
	header := make([]int, headerLength)
	for i := range header {
		if i > 10 {
			header[i] = int(bitbuf.getbits(r, bigBitSize))
		} else {
			header[i] = int(bitbuf.getbits(r, bitSize))
		}
	}
	numPoints := header[9]
	numFlags := header[10]
	numCodes := header[11]

	startPoints := make([]int, numPoints)
	for i := range startPoints {
		startPoints[i] = int(bitbuf.getbits(r, bitSize))
	}

	flagFreqMax := 1
	codeFreqMax := 8

	flagFreq := make([]uint64, flagFreqMax+1)
	for i := range flagFreq {
		flagFreq[i] = uint64(bitbuf.getbits(r, bigBitSize))
	}
	flagPmodel := newDecPmodel(flagFreq, 0, uint(flagFreqMax))

	codeFreq := make([]uint64, codeFreqMax+1)
	for i := range codeFreq {
		codeFreq[i] = uint64(bitbuf.getbits(r, bigBitSize))
	}
	codePmodel := newDecPmodel(codeFreq, 0, uint(codeFreqMax))

	rc := newRangeCoder()
	rc.startdec(r)

	flags := make([]byte, numFlags)
	for i := range flags {
		flags[i] = byte(rc.decode(r, flagPmodel))
	}
	codes := make([]byte, numCodes)
	for i := range codes {
		codes[i] = byte(rc.decode(r, codePmodel))
	}

	stream := &decoder.Stream{
		StartPoints: startPoints,
		Flags:       flags,
		Codes:       codes,
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
