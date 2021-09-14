package codec

import (
	"bufio"
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
	flagPmodel := newEncPmodel(stream.Flags, 0, 1)
	codePmodel := newEncPmodel(stream.Codes, 0, 8)

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

	// チェーンコードのスタートポイントの符号化
	for _, point := range stream.StartPoints {
		bb.putbits(w, bitSize, uint(point))
	}

	// ビットフラグの確率モデル
	for _, freq := range flagPmodel.freq {
		bb.putbits(w, bigBitSize, uint(freq))
	}
	// チェーンコードの確率モデル
	for _, freq := range codePmodel.freq {
		bb.putbits(w, bigBitSize, uint(freq))
	}

	rc := newRangeCoder()
	// ビットフラグの算術符号化
	for _, flag := range stream.Flags {
		rc.encode(w, flagPmodel, uint64(flag))
	}
	// チェーンコードの算術符号化
	for _, code := range stream.Codes {
		rc.encode(w, codePmodel, uint64(code))
	}
	rc.finishenc(w)

	w.Flush()
}

func newEncPmodel(val []byte, min, max uint) *Pmodel {
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
