package codec

import (
	"bufio"
	"math"
	"os"

	"github.com/shimehituzi/pccbi/internal/processing"
)

func Encode(stream processing.Stream, filename string) {
	fp, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	w := bufio.NewWriter(fp)

	bb := NewBitbuf(true)
	codesPmodel := newEncPmodel(stream.Codes)

	bitSize := 16
	bigBitSize := 32

	// ヘッダー情報
	for _, header := range stream.Header {
		bb.putbits(w, bitSize, uint(header))
	}
	bb.putbits(w, bigBitSize, uint(len(stream.OuterCodes)))
	bb.putbits(w, bigBitSize, uint(len(stream.InnerCodes)))

	// チェーンコードの確率モデル
	for _, freq := range codesPmodel.freq {
		bb.putbits(w, bitSize, uint(freq))
	}
	for _, cumfreq := range codesPmodel.cumfreq {
		bb.putbits(w, bitSize, uint(cumfreq))
	}
	bb.putbits(w, bigBitSize, uint(codesPmodel.totfreq))
	bb.putbits(w, bitSize, uint(codesPmodel.offset))

	// チェーンコードのスタートポイントの符号化
	for _, outerPoint := range stream.OuterStartPoints {
		bb.putbits(w, bitSize, uint(outerPoint))
	}
	for _, innerPoint := range stream.InnerStartPoints {
		bb.putbits(w, bitSize, uint(innerPoint))
	}

	// チェーンコードの算術符号化
	rc := newRangeCoder()
	for _, code := range stream.Codes {
		rc.encode(w, codesPmodel, uint64(code))
	}
	rc.finishenc(w)

	w.Flush()
}

func newEncPmodel(val []byte) *Pmodel {
	var (
		min byte = math.MaxUint8
		max byte = 0
	)
	for _, v := range val {
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}

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
