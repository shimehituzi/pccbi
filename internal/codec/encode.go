package codec

import (
	"bufio"
	"math"
	"os"

	"github.com/shimehituzi/pccbi/internal/encoder"
)

func Encode(stream *encoder.Stream) {
	fp, err := os.Create("compressed")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	w := bufio.NewWriter(fp)

	bb := NewBitbuf(true)
	outerPmodel := newEncPmodel(stream.OuterCodes)
	innerPmodel := newEncPmodel(stream.InnerCodes)

	bitSize := 16
	bigBitSize := 32

	// ヘッダー情報
	for i, header := range stream.Header {
		if i > 10 {
			bb.putbits(w, bigBitSize, uint(header))
		} else {
			bb.putbits(w, bitSize, uint(header))
		}
	}

	// 外輪郭チェーンコードの確率モデル
	for _, freq := range outerPmodel.freq {
		bb.putbits(w, bigBitSize, uint(freq))
	}
	for _, cumfreq := range outerPmodel.cumfreq {
		bb.putbits(w, bigBitSize, uint(cumfreq))
	}
	bb.putbits(w, bigBitSize, uint(outerPmodel.totfreq))
	bb.putbits(w, bitSize, uint(outerPmodel.offset))

	// 内輪郭チェーンコードの確率モデル
	for _, freq := range innerPmodel.freq {
		bb.putbits(w, bigBitSize, uint(freq))
	}
	for _, cumfreq := range innerPmodel.cumfreq {
		bb.putbits(w, bigBitSize, uint(cumfreq))
	}
	bb.putbits(w, bigBitSize, uint(innerPmodel.totfreq))
	bb.putbits(w, bitSize, uint(innerPmodel.offset))

	// チェーンコードのスタートポイントの符号化
	for _, outerPoint := range stream.OuterStartPoints {
		bb.putbits(w, bitSize, uint(outerPoint))
	}
	for _, innerPoint := range stream.InnerStartPoints {
		bb.putbits(w, bitSize, uint(innerPoint))
	}

	// チェーンコードの算術符号化
	rc := newRangeCoder()
	for _, outerCode := range stream.OuterCodes {
		rc.encode(w, outerPmodel, uint64(outerCode))
	}
	for _, innerCode := range stream.InnerCodes {
		rc.encode(w, innerPmodel, uint64(innerCode))
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
