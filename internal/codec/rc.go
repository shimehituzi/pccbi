package codec

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
)

const SIZE = 64
const TOP = uint64(1) << (SIZE - 8)
const BOT = uint64(1) << (SIZE - 16)

type RangeCoder struct {
	range_, low, code uint64
}

type Pmodel struct {
	freq, cumfreq   []uint64
	totfreq, offset uint64
}

func newRangeCoder() *RangeCoder {
	return &RangeCoder{
		range_: math.MaxUint64,
		low:    0,
		code:   0,
	}
}

func (rc *RangeCoder) encode(w *bufio.Writer, pm *Pmodel, val uint64) {
	rc.range_ /= pm.totfreq
	rc.low += pm.cumfreq[val] * rc.range_
	rc.range_ *= pm.freq[val]

	for (rc.low ^ (rc.low + rc.range_)) < TOP {
		w.WriteByte(byte(rc.low >> (SIZE - 8)))
		rc.code += 8
		if rc.code > 1e8 {
			panic("L T")
		}
		rc.range_ <<= 8
		rc.low <<= 8
	}

	for rc.range_ < BOT {
		w.WriteByte(byte(rc.low >> (SIZE - 8)))
		rc.code += 8
		if rc.code > 1e8 {
			panic("L T")
		}
		rc.range_ = ((-rc.low) & (BOT - 1)) << 8
		rc.low <<= 8
	}
}

func (rc *RangeCoder) finishenc(w *bufio.Writer) (bits uint64) {
	for i := 0; i < SIZE; i += 8 {
		w.WriteByte(byte(rc.low >> (SIZE - 8)))
		rc.code += 8
		rc.low <<= 8
	}
	w.Flush()

	bits = rc.code
	return
}

func (rc *RangeCoder) startdec(r *bufio.Reader) {
	for i := 0; i < SIZE; i += 8 {
		buf, err := r.ReadByte()
		if err != nil && err != io.EOF {
			panic(err)
		}
		rc.code = (rc.code << 8) | uint64(buf)
	}
}

func (rc *RangeCoder) decode(r *bufio.Reader, pm *Pmodel) (val uint64) {
	var rfreq uint64
	rc.range_ /= pm.totfreq
	rfreq = (rc.code - rc.low) / rc.range_
	if rfreq >= pm.totfreq {
		panic("C D")
	}

	var max, middle uint64
	val = 0
	max = uint64(len(pm.freq) - 1)
	for val < max {
		middle = (val + max) / 2
		if pm.cumfreq[middle+1]-pm.offset <= rfreq {
			val = middle + 1
		} else {
			max = middle
		}
	}

	rc.low += (pm.cumfreq[val] - pm.offset) * rc.range_
	rc.range_ *= pm.freq[val]
	for (rc.low ^ (rc.low + rc.range_)) < TOP {
		buf, err := r.ReadByte()
		if err != nil && err != io.EOF {
			panic(err)
		}
		rc.code = (rc.code << 8) | uint64(buf)
		rc.range_ <<= 8
		rc.low <<= 8
	}
	for rc.range_ < BOT {
		buf, err := r.ReadByte()
		if err != nil && err != io.EOF {
			panic(err)
		}
		rc.code = (rc.code << 8) | uint64(buf)
		rc.range_ = ((-rc.low) & (BOT - 1)) << 8
		rc.low <<= 8
	}

	return val
}

// ------------Test------------

func newTestEncPmodel(val []uint64) *Pmodel {
	var (
		min uint64 = math.MaxUint64
		max uint64 = 0
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

func newTestDecPmodel(freq []uint64, min, max uint64) *Pmodel {
	cumfreq := make([]uint64, max+2)
	cumfreq[0] = 0
	for i := range freq {
		cumfreq[i+1] = cumfreq[i] + freq[i]
	}
	offset := cumfreq[min]
	totfreq := cumfreq[max+1] - offset

	return &Pmodel{freq, cumfreq, totfreq, offset}
}

func getTestData() (val, freq []uint64, min, max uint64, length int) {
	val = []uint64{0, 0, 1, 0, 0}
	length = len(val)

	min = math.MaxUint64
	max = 0
	for _, v := range val {
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	freq = make([]uint64, max+1)
	for _, v := range val {
		freq[v]++
	}

	return
}

func testEnc() {
	val, _, _, _, _ := getTestData()

	pm := newTestEncPmodel(val)
	rc := newRangeCoder()

	filename := "compressed"
	fp, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	w := bufio.NewWriter(fp)

	for _, v := range val {
		rc.encode(w, pm, v)
	}

	bits := rc.finishenc(w)

	fmt.Println(bits)
}

func testDec() {
	_, freq, min, max, length := getTestData()

	pm := newTestDecPmodel(freq, min, max)
	rc := newRangeCoder()

	filename := "compressed"
	fp, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	r := bufio.NewReader(fp)

	rc.startdec(r)

	val := make([]uint64, length)
	for i := 0; i < length; i++ {
		val[i] = rc.decode(r, pm)
	}

	fmt.Println(val)
}
