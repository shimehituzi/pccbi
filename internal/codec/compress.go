package codec

import (
	"bufio"
	"os"

	"github.com/shimehituzi/pccbi/internal/bitstream"
)

const bitSize = 16
const bigBitSize = 32
const codeFreqMax = 8

// =================
//      encode
// =================

func Encode(dstPath string, stream *Stream, header *Header) (int, int) {
	fp, err := os.Create(dstPath)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	w := bufio.NewWriter(fp)

	bb := bitstream.NewBitbuf(true)
	numCodesArrayFreqmax := uint(0)
	for _, v := range stream.NumCodesArray {
		if v > numCodesArrayFreqmax {
			numCodesArrayFreqmax = v
		}
	}
	numCodesArrayPmodel := bitstream.NewEncPmodel(stream.NumCodesArray, 0, numCodesArrayFreqmax)
	codePmodel := bitstream.NewEncPmodel(stream.Codes, 0, codeFreqMax)

	// ヘッダー情報
	headerBits := 0
	headerBits += bb.Putbits(w, bitSize, uint(header.Axis))
	for _, v := range header.Length {
		headerBits += bb.Putbits(w, bitSize, uint(v))
	}
	for _, v := range header.Bias {
		headerBits += bb.Putbits(w, bitSize, uint(v))
	}
	headerBits += bb.Putbits(w, bigBitSize, uint(len(stream.Codes)))
	headerBits += bb.Putbits(w, bitSize, uint(len(stream.StartPoints)))
	headerBits += bb.Putbits(w, bitSize, uint(len(stream.NumCodesArray)))
	headerBits += bb.Putbits(w, bitSize, numCodesArrayFreqmax)

	// チェーンコードのスタートポイントの符号化
	for _, point := range stream.StartPoints {
		for i := 0; i < 3; i++ {
			headerBits += bb.Putbits(w, bitSize, uint(point[i]))
		}
	}

	// セグメント内の輪郭数の確率モデル
	for _, freq := range numCodesArrayPmodel.Freq() {
		headerBits += bb.Putbits(w, bigBitSize, uint(freq))
	}
	// チェーンコードの確率モデル
	for _, freq := range codePmodel.Freq() {
		headerBits += bb.Putbits(w, bigBitSize, uint(freq))
	}

	rc := bitstream.NewRangeCoder()
	// セグメント内の輪郭数の算術符号化
	for _, numCodesArray := range stream.NumCodesArray {
		rc.Encode(w, numCodesArrayPmodel, uint64(numCodesArray))
	}
	// チェーンコードの算術符号化
	for _, code := range stream.Codes {
		rc.Encode(w, codePmodel, uint64(code))
	}
	dataBits := int(rc.Finishenc(w))

	w.Flush()

	return dataBits, headerBits
}

// =================
//      decode
// =================

func Decode(dstPath string) (*Stream, *Header) {
	fp, err := os.Open(dstPath)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	r := bufio.NewReader(fp)

	bitbuf := bitstream.NewBitbuf(false)

	header := new(Header)
	header.Axis = Axis(bitbuf.Getbits(r, bitSize))
	for i := range header.Length {
		header.Length[i] = int(bitbuf.Getbits(r, bitSize))
	}
	for i := range header.Bias {
		header.Bias[i] = int(bitbuf.Getbits(r, bitSize))
	}
	codesLength := bitbuf.Getbits(r, bigBitSize)
	startPointsLength := bitbuf.Getbits(r, bitSize)
	numCodesArrayLength := bitbuf.Getbits(r, bitSize)
	numCodesArrayFreqMax := bitbuf.Getbits(r, bitSize)

	startPoints := make([][3]uint, startPointsLength)
	for i := range startPoints {
		for j := 0; j < 3; j++ {
			startPoints[i][j] = bitbuf.Getbits(r, bitSize)
		}
	}

	numCodesArrayFreq := make([]uint64, numCodesArrayFreqMax+1)
	for i := range numCodesArrayFreq {
		numCodesArrayFreq[i] = uint64(bitbuf.Getbits(r, bigBitSize))
	}
	numCodesArrayPmodel := bitstream.NewDecPmodel(numCodesArrayFreq, 0, uint(numCodesArrayFreqMax))

	codesFreq := make([]uint64, codeFreqMax+1)
	for i := range codesFreq {
		codesFreq[i] = uint64(bitbuf.Getbits(r, bigBitSize))
	}
	codesPmodel := bitstream.NewDecPmodel(codesFreq, 0, uint(codeFreqMax))

	rc := bitstream.NewRangeCoder()
	rc.Startdec(r)

	numCodesArray := make([]uint, numCodesArrayLength)
	for i := range numCodesArray {
		numCodesArray[i] = uint(rc.Decode(r, numCodesArrayPmodel))
	}
	codes := make([]uint, codesLength)
	for i := range codes {
		codes[i] = uint(rc.Decode(r, codesPmodel))
	}

	stream := &Stream{
		StartPoints:   startPoints,
		NumCodesArray: numCodesArray,
		Codes:         codes,
	}

	return stream, header
}
