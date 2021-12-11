package bitstream

import (
	"bufio"
	"os"

	"github.com/shimehituzi/pccbi/internal/codec"
)

func Encode(pccPath string, stream *codec.Stream, header *codec.Header) (int, int) {
	fp, err := os.Create(pccPath)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	w := bufio.NewWriter(fp)

	bb := NewBitbuf(true)
	numCodesArrayFreqmax := uint(0)
	for _, v := range stream.NumCodesArray {
		if v > numCodesArrayFreqmax {
			numCodesArrayFreqmax = v
		}
	}
	numCodesArrayPmodel := newEncPmodel(stream.NumCodesArray, 0, numCodesArrayFreqmax)
	codePmodel := newEncPmodel(stream.Codes, 0, 8)

	bitSize := 16
	bigBitSize := 32

	headerBits := 0

	// ヘッダー情報
	headerBits += bb.putbits(w, bitSize, uint(header.Axis))
	for _, v := range header.Length {
		headerBits += bb.putbits(w, bitSize, uint(v))
	}
	for _, v := range header.Bias {
		headerBits += bb.putbits(w, bitSize, uint(v))
	}
	headerBits += bb.putbits(w, bigBitSize, uint(len(stream.Codes)))
	headerBits += bb.putbits(w, bitSize, uint(len(stream.StartPoints)))
	headerBits += bb.putbits(w, bitSize, uint(len(stream.NumCodesArray)))
	headerBits += bb.putbits(w, bitSize, numCodesArrayFreqmax)

	// チェーンコードのスタートポイントの符号化
	for _, point := range stream.StartPoints {
		for i := 0; i < 3; i++ {
			headerBits += bb.putbits(w, bitSize, uint(point[i]))
		}
	}

	// セグメント内の輪郭数の確率モデル
	for _, freq := range numCodesArrayPmodel.freq {
		headerBits += bb.putbits(w, bigBitSize, uint(freq))
	}
	// チェーンコードの確率モデル
	for _, freq := range codePmodel.freq {
		headerBits += bb.putbits(w, bigBitSize, uint(freq))
	}

	rc := newRangeCoder()
	// セグメント内の輪郭数の算術符号化
	for _, numCodesArray := range stream.NumCodesArray {
		rc.encode(w, numCodesArrayPmodel, uint64(numCodesArray))
	}
	// チェーンコードの算術符号化
	for _, code := range stream.Codes {
		rc.encode(w, codePmodel, uint64(code))
	}
	dataBits := int(rc.finishenc(w))

	w.Flush()

	return dataBits, headerBits
}

func newEncPmodel(val []uint, min, max uint) *Pmodel {
	freq := make([]uint64, max+1)
	for _, v := range val {
		freq[v]++
	}
	cumfreq := make([]uint64, max+2)
	cumfreq[0] = 0
	for i := range freq {
		cumfreq[i+1] = cumfreq[i] + freq[i]
	}
	offset := cumfreq[min]
	totfreq := cumfreq[max+1] - offset

	return &Pmodel{freq, cumfreq, totfreq, offset}
}
